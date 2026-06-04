# Lab: Bad Service Name

## Scenario
A typo in the Service selector causes it to stop matching pod labels.
Traffic stops routing even though the pods are healthy.

## Symptoms
- `kubectl port-forward` to the Service times out
- `curl` requests hang or fail

## How to diagnose
```bash
kubectl get endpoints -n inventory
kubectl describe service inventory-api -n inventory | grep -A 3 "Selector\|Endpoints"
```
Look for `Endpoints: <none>` and a mismatched `Selector` value.

## Root cause
The Service `selector` must exactly match the pod `labels`. A typo means
zero pods are registered as endpoints. The Service exists but routes nowhere, 
pods are invisible. 

## Fix
```bash
kubectl patch service inventory-api -n inventory \
  -p '{"spec":{"selector":{"app":"inventory-api"}}}'
```

## What I learned
- Services find pods via label selectors. Typos drop all traffic
- `kubectl get endpoints` is the fastest way to confirm a Service has registered pods
- Pods can be healthy while being completely unreachable if the Service selector is wrong
