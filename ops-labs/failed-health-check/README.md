# Lab: Failed Health Check

## Scenario
Both liveness and readiness probes are pointed at the wrong port. The app
is healthy but Kubernetes thinks it's broken and keeps restarting it.

## Symptoms
- New pod shows RESTARTS climbing (1, 2, 3...)
- Pod stays `0/1 Ready` and never enters Service rotation
- Old pod keeps running and serving traffic (rolling update protection)

## How to diagnose
```bash
kubectl describe pod -n inventory -l app=inventory-api | grep -A 15 "Events:"
```
Look for:
```
Liveness probe failed: connection refused on port 9999
Readiness probe failed: connection refused on port 9999
Container api failed liveness probe, will be restarted
```

## Root cause
Probes were pointing at port 9999 instead of 8080. The app was healthy but 
nothing was listening on 9999 so every probe got `connection refused`.
Kubernetes cannot distinguish between a broken app and misconfigured probes.
It acts on what the probes report.

## Fix
```bash
kubectl apply -f inventory-api/k8s/api-deployment.yaml
```
Restores the correct probe configuration from the source of truth YAML.

## Broken config
```yaml
livenessProbe:
  httpGet:
    port: 9999   # wrong — app listens on 8080
readinessProbe:
  httpGet:
    port: 9999   # wrong — app listens on 8080
```

## What I learned
- Probes are Kubernetes asking the pod if it is "healthy". Misconfigured probes
  look identical to a broken app from Kubernetes' perspective
- The rolling update kept the old pod alive the entire time to ensure zero downtime
- Image present + container started does not mean the app is healthy
