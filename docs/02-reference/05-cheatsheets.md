# Cheatsheets

This is the fast command reference for the lab and for common DevOps tasks around it.

Use this while working. Use the deeper guides when you need reasoning and context.

## Rule of use

Use cheatsheets for recall, not for blind execution.

If you do not know what a command is changing, stop and read:

- [Concepts, What, Why, When](/Users/deanlofts/gits/devops-roundrobin/docs/02-reference/02-concepts-and-why.md)
- [When You Get Stuck](/Users/deanlofts/gits/devops-roundrobin/docs/02-reference/01-when-stuck.md)

## Repo and toolchain

Enter the pinned toolchain:

```bash
cd /Users/deanlofts/gits/devops-roundrobin
cp .env.example .env
nix develop
make bootstrap
make help
```

Quick validation:

```bash
go version
kind version
kubectl version --client
helm version --short
terraform version
docker compose version
```

## Compose

Bring the stack up:

```bash
make compose-up
make compose-ps
```

Stop and clean it:

```bash
make compose-down
```

Useful checks:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:18081/health
docker compose logs --tail=50 app postgres redis worker
```

Queue flow:

```bash
curl -X POST http://localhost:8080/jobs \
  -H 'content-type: application/json' \
  -d '{"payload":"hello"}'

curl http://localhost:8080/jobs
```

## Docker

Build images:

```bash
make docker-build
```

Inspect local images:

```bash
docker images | grep round-robin
```

Inspect running containers:

```bash
docker ps
docker logs <container-name>
```

## kind

Create cluster:

```bash
make kind-up
```

Delete cluster and local registry:

```bash
make kind-down
```

Fix kubeconfig context:

```bash
kind export kubeconfig --name round-robin
kubectl config current-context
kubectl cluster-info
```

Load images:

```bash
make kind-load
```

Install ingress:

```bash
make ingress-install
```

## Kubernetes

Apply raw manifests:

```bash
make k8s-apply-raw
```

Delete raw manifests:

```bash
make k8s-delete-raw
```

Most-used inspection commands:

```bash
kubectl get nodes -o wide
kubectl get pods -A
kubectl get all -n round-robin
kubectl get all -n round-robin-raw
kubectl get ingress -A
kubectl get events -n round-robin --sort-by=.metadata.creationTimestamp
kubectl describe deployment round-robin-app -n round-robin
kubectl logs deployment/round-robin-app -n round-robin
kubectl logs deployment/round-robin-worker -n round-robin
```

Watch rollout:

```bash
kubectl rollout status deployment/round-robin-app -n round-robin
kubectl rollout history deployment/round-robin-app -n round-robin
kubectl rollout undo deployment/round-robin-app -n round-robin
```

Set env var quickly:

```bash
kubectl -n round-robin set env deployment/round-robin-app KEY=value
kubectl -n round-robin set env deployment/round-robin-app KEY-
```

Set image quickly:

```bash
kubectl -n round-robin set image deployment/round-robin-app app=round-robin-app:local
```

## Ingress and app access

App via ingress:

```bash
curl http://app.127.0.0.1.nip.io/health
curl http://app.127.0.0.1.nip.io/ready
```

Check endpoints behind ingress:

```bash
kubectl get ingress -n round-robin
kubectl get svc -n round-robin
kubectl get endpoints -n round-robin
kubectl get pods -n ingress-nginx
```

## Helm

Lint:

```bash
make helm-lint
```

Install or upgrade:

```bash
make helm-install
helm upgrade round-robin infra/helm/round-robin -n round-robin
```

List releases:

```bash
helm list -A
helm history round-robin -n round-robin
```

Rollback:

```bash
make helm-rollback
helm rollback round-robin -n round-robin
```

Render templates locally:

```bash
helm template round-robin infra/helm/round-robin
```

## Terraform

Prepare handoff from manual Helm:

```bash
make tf-prepare
```

Initialize:

```bash
make tf-init
```

Plan:

```bash
make tf-plan
```

Apply:

```bash
make tf-apply
```

Destroy:

```bash
make tf-destroy
```

Useful direct commands:

```bash
terraform -chdir=infra/terraform/environments/local validate
terraform -chdir=infra/terraform/environments/local output
terraform -chdir=infra/terraform/environments/local state list
```

## Monitoring

Install monitoring:

```bash
make monitor-install
kubectl get pods -n monitoring
```

Open:

- `http://grafana.127.0.0.1.nip.io`
- `http://prometheus.127.0.0.1.nip.io`

Grafana login:

- user: `admin`
- password: `admin`

## Security and scanning

Scan images:

```bash
make scan-image
```

Review secrets and config use:

```bash
kubectl get secret -n round-robin
kubectl get configmap -n round-robin
```

## Load and scaling

Run load test:

```bash
make load-test
```

Check scaling:

```bash
kubectl get hpa -n round-robin
kubectl top pods -n round-robin
kubectl get deploy -n round-robin
```

## Incident drills

Break DB credentials:

```bash
./scripts/drills/break-db-creds.sh round-robin round-robin-app
./scripts/drills/restore-db-creds.sh round-robin round-robin-app
```

Break image:

```bash
./scripts/drills/break-image.sh round-robin round-robin-app app
./scripts/drills/rollback-deployment.sh round-robin round-robin-app
```

Break readiness:

```bash
./scripts/drills/break-readiness-raw.sh
./scripts/drills/restore-readiness-raw.sh
```

Break ingress:

```bash
kubectl scale deployment ingress-nginx-controller -n ingress-nginx --replicas=0
kubectl scale deployment ingress-nginx-controller -n ingress-nginx --replicas=1
```

## DNS, TLS, HTTP quick checks

DNS:

```bash
dig app.127.0.0.1.nip.io
nslookup app.127.0.0.1.nip.io
```

HTTP:

```bash
curl -i http://app.127.0.0.1.nip.io/health
curl -I http://app.127.0.0.1.nip.io/health
```

TLS and certs for a real external endpoint:

```bash
openssl s_client -connect example.com:443 -servername example.com
curl -Iv https://example.com
```

## SSO and identity debugging quick prompts

Use these questions before you touch config:

1. Is this authentication failure or authorization failure?
2. Is login broken for one user, one group, or everyone?
3. Did the token arrive, but the role mapping fail?
4. Is this a provisioning issue instead of a login issue?
5. In AWS, is the user missing the right IAM Identity Center assignment?

## Best companions to this cheatsheet

- [Start-to-Finish Walkthrough](/Users/deanlofts/gits/devops-roundrobin/docs/01-walkthrough/01-start-to-finish.md)
- [When You Get Stuck](/Users/deanlofts/gits/devops-roundrobin/docs/02-reference/01-when-stuck.md)
- [Federated Identity and SSO](/Users/deanlofts/gits/devops-roundrobin/docs/02-reference/03-federated-identity-and-sso.md)
- [From URL to Rendered Page](/Users/deanlofts/gits/devops-roundrobin/docs/02-reference/04-from-url-to-rendered-page.md)
