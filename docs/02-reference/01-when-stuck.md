# When You Get Stuck

This guide is the fast recovery path.

Use it when a command fails and you want the next sensible move without spiralling.

## First rule

Do not change three things at once.

Use this sequence:

1. identify which layer failed
2. inspect that layer directly
3. make one change
4. re-test

## Layer map

There are four main layers in this lab:

1. local toolchain
2. container runtime and Compose
3. Kubernetes and ingress
4. Helm and Terraform control layer

## Toolchain failures

### `nix develop` fails

Check:

```bash
nix flake metadata
```

If this fails:

- the issue is Nix or connectivity
- do not debug app code yet

### `make bootstrap` fails

Check:

```bash
make deps
```

This tells you whether Docker, kind, kubectl, Helm, Terraform, jq, yq, or Trivy are missing from the active shell.

### `go: command not found`

You are outside `nix develop`.

Fix:

```bash
nix develop
```

## Compose failures

### `docker compose up` fails immediately

Check:

```bash
docker info
docker compose config
```

Common causes:

- Docker runtime is not running
- `.env` is missing
- a host port is already taken

### App starts but `/ready` fails

Check:

```bash
docker compose ps
docker compose logs --tail=50 app postgres redis worker
curl -i http://localhost:8080/ready
```

Interpretation:

- if Postgres is down, readiness should fail
- if Redis is down, readiness should fail
- if `/health` fails too, the process itself is in trouble

### Jobs are not being processed

Check:

```bash
curl -X POST http://localhost:8080/jobs -H 'content-type: application/json' -d '{"payload":"debug"}'
curl http://localhost:8080/jobs
docker compose logs worker
docker compose logs app
```

Likely causes:

- worker is not running
- Redis is unavailable
- Postgres is unavailable for persistence

## Kubernetes failures

### `make kind-up` fails

Check:

```bash
docker ps
kind get clusters
```

This repo pins a known-good node image in [`.env.example`](/Users/deanlofts/gits/devops-roundrobin/.env.example). Do not change it unless you are debugging kind itself.

### `kubectl` cannot connect

Run:

```bash
kind export kubeconfig --name round-robin
kubectl config current-context
kubectl cluster-info
```

Expected context:

```text
kind-round-robin
```

### Pods are Pending or CrashLooping

Check:

```bash
kubectl get pods -A
kubectl describe pod <pod-name> -n <namespace>
kubectl logs <pod-name> -n <namespace>
kubectl get events -n <namespace> --sort-by=.metadata.creationTimestamp
```

Interpret by symptom:

- `ImagePullBackOff`: wrong image or image not loaded into kind
- `CrashLoopBackOff`: app starts and dies repeatedly
- `Pending`: scheduling, PVC, or resource issue

### Ingress returns `503`

Check:

```bash
kubectl get pods -n ingress-nginx
kubectl get ingress -A
kubectl get svc -n round-robin
kubectl get endpoints -n round-robin
```

Usually this means:

- ingress controller is not running
- backing Service has no healthy endpoints yet
- app readiness has not passed yet

Wait a little before assuming the deployment is broken, especially right after install.

## Helm failures

### `helm install` or `helm upgrade` fails

Check:

```bash
make helm-lint
helm list -A
kubectl get all -n round-robin
kubectl get events -n round-robin --sort-by=.metadata.creationTimestamp
```

Then inspect the rendered values path:

```bash
helm template round-robin infra/helm/round-robin
```

### Release is deployed but app is unavailable

Check:

```bash
kubectl get pods -n round-robin
kubectl logs deployment/round-robin-app -n round-robin
curl -i http://app.127.0.0.1.nip.io/ready
```

## Terraform failures

### `terraform init` fails

Check:

```bash
terraform -chdir=infra/terraform/environments/local version
terraform -chdir=infra/terraform/environments/local init -backend=false
```

If provider downloads fail, it is usually network or registry access.

### `terraform plan` fails on Kubernetes connectivity

Check:

```bash
kubectl config current-context
terraform -chdir=infra/terraform/environments/local validate
```

Then confirm:

- kind cluster exists
- kubeconfig points at `kind-round-robin`

### You are unsure whether Terraform or Helm owns something

Read:

- [`infra/terraform/modules/platform/main.tf`](/Users/deanlofts/gits/devops-roundrobin/infra/terraform/modules/platform/main.tf)
- [`infra/helm/round-robin/values.yaml`](/Users/deanlofts/gits/devops-roundrobin/infra/helm/round-robin/values.yaml)

Remember:

- Terraform in this lab owns the Helm release and some namespaces
- Helm owns the app resources inside the release

## Monitoring failures

### Grafana or Prometheus not loading

Check:

```bash
kubectl get pods -n monitoring
kubectl get ingress -n monitoring
kubectl get events -n monitoring --sort-by=.metadata.creationTimestamp
```

Then verify ingress controller health again.

## Decision table

If command fails before anything starts:

- suspect toolchain or Docker

If containers run but app is not ready:

- suspect Postgres or Redis

If Kubernetes pods run but ingress gives `503`:

- suspect readiness or endpoints

If Helm deploys but service is still broken:

- inspect the rendered chart and pod events

If Terraform errors on provider or cluster access:

- suspect kubeconfig or provider init, not app code
