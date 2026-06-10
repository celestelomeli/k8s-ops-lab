# Phase 4 Observations

## Managed control plane, not managed workloads

I was surprised that moving to EKS did not mean writing less YAML. The control plane (etcd, API server, scheduler, controller manager) is now AWS's problem, which is real value, but everything above it is still mine. I still had to write 14+ manifests to get three services running: namespace, configmap, secret store, external secret, two deployments plus postgres, three services, a PVC, a storage class, an HPA, and a metrics-server install. "Managed Kubernetes" manages the cluster, not the app. The mental model that stuck: EKS replaces kubeadm, not kubectl.

## Secrets: from secret.yaml to AWS Secrets Manager via ESO

Replaced the hardcoded `secret.yaml` with the External Secrets Operator (ESO) reading from AWS Secrets Manager.

How it fits together:
- The secret lives in AWS Secrets Manager, never in git.
- ESO runs in the cluster and syncs that value into a normal Kubernetes Secret on a refresh interval.
- A `ClusterSecretStore` says "where to read from" (Secrets Manager, us-east-1). An `ExternalSecret` says "read key X and write it into Kubernetes Secret Y."
- ESO authenticates to AWS through IRSA, not static keys.

**IRSA (IAM Roles for Service Accounts):** the cluster has an OIDC provider. A Kubernetes service account is annotated with an IAM role ARN. When the ESO pod presents its service account token, AWS trusts the OIDC issuer and hands back temporary credentials scoped to that role. So the pod gets `secretsmanager:GetSecretValue` on `k8s-ops-lab/*` and nothing more, with no long-lived access keys anywhere. This is the part that felt genuinely production-grade.

A gotcha worth noting is that ESO installed into the `default` namespace, not `external-secrets`. The IAM trust policy condition has to match the real service account location (`system:serviceaccount:default:external-secrets`), or the role assumption fails.

## The bug that ate the afternoon: JSON secrets

Every service crashed with `password authentication failed for user "inventory"`, and I spent a long time blaming scram-sha-256 auth. The real cause was much simpler: the secret in Secrets Manager was stored as JSON:

```json
{"DB_PASSWORD": "..."}
```

My `ExternalSecret` grabbed the whole secret and dumped it into the `DB_PASSWORD` key, so the password the app actually sent was the entire string `{"DB_PASSWORD":"..."}`. Postgres, initialized from the same secret, then also baked in that blob as the password, which is why it was hard to spot.

The fix was to tell ESO which JSON field to extract:

```yaml
remoteRef:
  key: k8s-ops-lab/inventory/db-password
  property: DB_PASSWORD
```

Lesson: when a Secrets Manager value is JSON, you almost always want `property:`. 

## Persistent storage and the lost+found trap

Postgres uses a PVC backed by EBS through the EBS CSI driver (another IRSA role). Fresh EBS volumes are formatted ext4, which creates a `lost+found` directory at the root. Postgres refuses to initialize into a non-empty data directory, so it crash-looped immediately.

Fix is the standard one: point `PGDATA` at a subdirectory (`/var/lib/postgresql/data/pgdata`) so initdb sees an empty dir while the volume root keeps its `lost+found`. 

## Image architecture mismatch

Built the images on an Apple Silicon Mac, so they came out arm64. EKS nodes are amd64, and the pods died with `exec format error` because a binary built for one CPU type cannot run on the other. Fix: `docker build --platform linux/amd64`.

The real lesson is that the machine you build on decides the image's CPU type, and a laptop is an inconsistent build machine. On a real team you avoid this two ways. First, an automated build server (CI) always builds on the same kind of machine, so the image always matches the cluster and nobody has to remember the platform flag. Second, a multi-arch image bundles both arm64 and amd64 into one tag, and each machine automatically pulls the version that matches it, so the same image runs anywhere.

## Large CRDs need server-side apply

Installing ESO failed with the 262KB annotation limit, because `kubectl apply` stashes a full copy of the resource in an annotation for diffing. ESO's CRDs are huge. `kubectl apply --server-side` moves that bookkeeping to the API server and sidesteps the limit. Good to know for any operator with big CRDs.

## HPA needs metrics-server, which EKS does not ship

The HPA sat at `cpu: <unknown>/50%` and never scaled. EKS does not install metrics-server by default (unlike Docker Desktop). Once installed, the HPA read real CPU and scaling worked: I drove load from an in-cluster pod, watched CPU cross 50%, and saw replicas climb from 2 toward 5, then settle back after the cooldown window.


## Observability via CloudWatch Container Insights

Enabled through the `amazon-cloudwatch-observability` EKS add-on. It runs a CloudWatch agent DaemonSet (metrics) and Fluent Bit DaemonSet (logs) on every node. Having the `CloudWatchAgentServerPolicy` on the nodes only grants permission, the add-on is what actually collects. Container Insights gave per-pod CPU/memory/network dashboards, restart counts (today's CrashLoopBackOff churn shows up as a spike), and searchable container logs under `/aws/containerinsights/k8s-ops-lab/`. It bills as custom metrics plus log ingestion, so it is not free to leave running.

## Making it less abstract: k9s

In Kubernetes nothing is visible. The "thing running" is spread across machines you cannot point at, and `kubectl get` only shows a frozen snapshot. k9s (`brew install k9s`) is a live full-screen terminal dashboard for the cluster. Typing `:pods`, `:svc`, `:hpa` jumps between resource types, and single keys act on the highlighted row: `l` for live logs, `d` to describe, `s` to shell into the container. Watching pods sit there and update in real time, instead of re-running a command, made the cluster feel like a real running system for the first time.

The exercise that taught me the most was self-healing. Watching `kubectl get pods -w` in one terminal while deleting a pod in another, the pod dies and a new one is recreated within seconds with no action from me. That is the whole idea of Kubernetes made visible: I declared the desired state, and the system constantly works to make reality match it.

## What I would do differently

- Store DB credentials in Secrets Manager as plain strings, or standardize on JSON everywhere with `property:` from the start. The mismatch — storing as JSON but reading it as a plain string — caused the whole auth rabbit hole.
- Build images in CI for the target architecture instead of locally.
- For auth failures, inspect the actual credential being sent (print it, decode it, compare it byte-for-byte) before theorizing about encryption modes or hashing. 