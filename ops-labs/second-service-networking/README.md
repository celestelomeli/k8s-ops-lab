# Lab: Second Service Networking

## What this demonstrates
Service-to-service communication using Kubernetes DNS. The inventory API
calls a notifier service by name after a product is created (no IP addresses,
no configuration, just the service name).

## Architecture
```
POST /products
      ↓
inventory-api pod
      ↓  http://notifier:9090/notify
notifier pod  ←  CoreDNS resolves "notifier" to the Service ClusterIP
      ↓
logs: "notification sent: product created: <name>"
```

## How DNS works in Kubernetes
Every Service gets a DNS entry automatically:
```
short name:  notifier                              
full name:   notifier.inventory.svc.cluster.local  
```
CoreDNS runs in `kube-system` and resolves these names to Service ClusterIPs.
Pods never need to know IP addresses since they use service names.

## Services deployed
- `notifier` — stateless Go HTTP service on port 9090
- Called asynchronously (goroutine) so it never blocks the API response

## How to test
```bash
kubectl port-forward service/inventory-api 8080:80 -n inventory

curl -X POST localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"sprocket","quantity":25,"price":4.99}'

kubectl logs -n inventory -l app=notifier
```
Expect: `notification sent: product created: sprocket`

## What I learned
- Services are reachable by name within the cluster — Kubernetes DNS handles resolution
- The same service discovery that makes `DB_HOST=postgres` work powers all
  service-to-service communication
- Stateless services need no PVC — they hold no data between requests
- Docker Compose does this within one machine; Kubernetes does it across any
  number of nodes and cloud regions with the same simplicity
