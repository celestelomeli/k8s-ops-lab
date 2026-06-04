# Lab: Broken DB Password

## Scenario
A database password rotation introduced a wrong credential in the Kubernetes
Secret. The API pod crashes on startup and enters CrashLoopBackOff.

## Symptoms
- Pod status cycles: `Running → Error → CrashLoopBackOff`
- Readiness probe fails: `connection refused` on `/health`
- Liveness probe fails: `connection refused` on `/health`

## How to diagnose
```bash
kubectl logs -n inventory -l app=inventory-api
```
Look for:
```
could not connect to database: password authentication failed for user "inventory"
```

```bash
kubectl describe pod -n inventory -l app=inventory-api
```
Look for probe failure events in the `Events:` section.

## Root cause
The Secret `inventory-secret` contains the wrong `DB_PASSWORD`. The app retries
10 times then calls `log.Fatal()`, crashing before port 8080 is ever bound.
Probes get `connection refused` — not because the app is slow, but because it
already exited.

## Fix
```bash
kubectl delete secret inventory-secret -n inventory
kubectl create secret generic inventory-secret \
  --namespace inventory \
  --from-literal=DB_PASSWORD=yourcorrectpassword

kubectl rollout restart deployment/inventory-api -n inventory
```

## What I learned
- Secrets are injected at container startup and changing a Secret does not
  update a running pod. You must restart to pick up new values.
- Readiness probe failure removes the pod from Service rotation before
  the liveness probe triggers a restart. Traffic stops before the termination.
- `log.Fatal()` exits the process immediately, preventing the HTTP server
  from ever starting. Probes will show `connection refused`, not an HTTP error.
  Connection refused is a network-level error,trying to open a TCP connection to an 
  address and port, and nothing answering.
