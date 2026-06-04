# Lab: ImagePullBackOff

## Scenario
A deployment references an image tag that doesn't exist in any
registry. The kubelet can't pull it and the pod never starts.

## Symptoms
- Pod stuck in `ErrImagePull` then `ImagePullBackOff`
- `READY: 0/1` — pod never starts
- Old pod stays running (rolling update protection)

## How to diagnose
```bash
kubectl describe pod -n inventory -l app=inventory-api | grep -A 10 "Events:"
```
Look for:
```
Failed to pull image "inventory-api:doesnotexist": pull access denied,
repository does not exist or may require authorization
```

## Root cause
The kubelet on the node tries to pull the image from a registry when it's not present 
in the node's local image store. If the image doesn't exist in any registry, the pull
fails and Kubernetes backs off with increasing wait times between retries.

## Fix
```bash
kubectl set image deployment/inventory-api api=inventory-api:v3 -n inventory
```

## What I learned
- `imagePullPolicy: IfNotPresent` only skips the pull if the image is already
  on the node. It does not prevent pulling if the image is absent
- For local development with kind, images must be explicitly loaded:
  `kind load docker-image <image>:<tag> --name <cluster>`
- `ErrImagePull` = first failure attempting to pull
- `ImagePullBackOff` = backing off after repeated failures, wait time increases
- The `:latest` tag defaults to `imagePullPolicy: Always` — always pulls from
  registry even if the image exists on the node
