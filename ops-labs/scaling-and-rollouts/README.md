# Lab: Scaling and Rollouts

## Scenario
Scale the inventory API from 1 to 3 replicas and observe Kubernetes
distributing traffic across all pods. Then trigger a rolling update and
watch pods replace one at a time with zero downtime.

## Commands

### Scale up
```bash
kubectl scale deployment inventory-api --replicas=3 -n inventory
kubectl get pods -n inventory -w
```

### Watch logs across all pods
```bash
kubectl logs -n inventory -l app=inventory-api -f
```

### Trigger a rolling restart
```bash
kubectl rollout restart deployment/inventory-api -n inventory
kubectl get pods -n inventory -w
```

### Scale back down
```bash
kubectl scale deployment inventory-api --replicas=1 -n inventory
```

## What to observe
- All 3 pods share the same Service. Traffic load balances automatically
- Rolling restart replaces one pod at a time
- New pod becomes Ready before old pod is terminated, so there is always at least 2 healthy pods
- `Error` status during termination is normal because the container process being killed

## What I learned
- `kubectl scale` changes replica count instantly so no manifest edit needed
- The Service automatically routes to all healthy pods, no config change required
- Rolling updates guarantee minimum availability, zero downtime by default
- Docker Compose has no concept of replicas. This is one of the core reasons
  to use Kubernetes for production workloads
