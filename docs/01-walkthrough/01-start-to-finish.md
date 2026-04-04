# Start-to-Finish Walkthrough

This is the straight path through the lab.

Use this when you want to move from zero to full platform confidence without having to decide what to do next.

## How to use this guide

- stay inside `nix develop`
- do one checkpoint at a time
- do not continue if the checkpoint fails
- use the linked stuck guides when needed

## Phase 0: Prepare the machine

Read first:

- [MacBook Setup](/Users/deanlofts/gits/devops-roundrobin/docs/00-macbook-setup.md)

Run:

```bash
cp .env.example .env
nix develop
make bootstrap
```

You are ready when:

- `make bootstrap` succeeds
- `make help` prints targets

## Phase 1: Understand the repo shape

Run:

```bash
find . -maxdepth 3 -type f | sort
```

Read:

- [`README.md`](/Users/deanlofts/gits/devops-roundrobin/README.md)
- [`Makefile`](/Users/deanlofts/gits/devops-roundrobin/Makefile)
- [`compose.yaml`](/Users/deanlofts/gits/devops-roundrobin/compose.yaml)

What you are learning:

- where app code lives
- where Docker and Compose live
- where Kubernetes manifests live
- where the Helm chart lives
- where the Terraform modules live

Why it matters:

A DevOps engineer spends a lot of time orienting quickly in unfamiliar repos.

## Phase 2: Run the app alone

Run:

```bash
make app-run
```

In another terminal:

```bash
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/metrics
```

Expected:

- `/` returns JSON describing the service
- `/health` returns `200`
- `/ready` returns `503` until Postgres and Redis exist
- `/metrics` returns Prometheus text output

What, why, when:

- `/health` answers “should this process be restarted?”
- `/ready` answers “should this process receive traffic yet?”
- `/metrics` answers “can we observe this service?”

Stop the app with `Ctrl+C` before moving on.

## Phase 3: Run the full local stack with Compose

Run:

```bash
make compose-up
make compose-ps
```

Check:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:18081/health
```

Expected:

- app is healthy
- app is ready
- worker is healthy
- Postgres and Redis are healthy in `docker compose ps`

What you are doing in role terms:

- bringing up a service stack
- validating dependencies
- checking readiness semantics

## Phase 4: Exercise the async workflow

Run:

```bash
curl -X POST http://localhost:8080/jobs \
  -H 'content-type: application/json' \
  -d '{"payload":"walkthrough job"}'

curl http://localhost:8080/jobs
```

Expected:

- the POST returns `202`
- the GET shows the job processed by `worker`

Why it matters:

This gives you a real producer-consumer workflow to debug later.

## Phase 5: Break a dependency and recover

Run:

```bash
docker compose stop postgres
curl http://localhost:8080/ready
docker compose logs --tail=50 app postgres
```

Then recover:

```bash
docker compose start postgres
curl http://localhost:8080/ready
```

What you should observe:

- health can remain up while readiness drops
- logs show dependency failures
- recovery does not require rebuilding images

This is exactly the kind of “service is up but not useful” distinction you need to be comfortable with in the job.

## Phase 6: Clean up Compose and move to Kubernetes

Run:

```bash
make compose-down
make kind-up
make ingress-install
make kind-load
```

Check:

```bash
kubectl get nodes -o wide
kubectl get pods -n ingress-nginx
```

Expected:

- one Ready node
- ingress controller running

What, why, when:

- `kind` gives you local Kubernetes without cloud cost
- ingress gives you realistic HTTP routing
- loading images avoids needing a remote registry for this lab

## Phase 7: Deploy raw Kubernetes resources

Run:

```bash
make k8s-apply-raw
kubectl get all -n round-robin-raw
kubectl get ingress -n round-robin-raw
```

Check:

```bash
curl http://app.127.0.0.1.nip.io/health
curl http://app.127.0.0.1.nip.io/ready
```

Expected:

- Postgres, Redis, app, and worker pods are running
- ingress serves the app

What you are learning:

- the primitive Kubernetes objects Helm and Terraform will later manage
- how Services, Deployments, Secrets, ConfigMaps, PVCs, and Ingress fit together

## Phase 8: Debug a raw Kubernetes incident

Run:

```bash
./scripts/drills/break-readiness-raw.sh
kubectl get pods -n round-robin-raw
kubectl describe deployment round-robin-app -n round-robin-raw
curl -i http://app.127.0.0.1.nip.io/ready
```

Recover:

```bash
./scripts/drills/restore-readiness-raw.sh
kubectl rollout status deployment/round-robin-app -n round-robin-raw
```

Why this matters:

It teaches you to use status, rollout history, readiness, and deployment config together.

## Phase 9: Replace raw resources with Helm

Run:

```bash
make k8s-delete-raw
make helm-lint
make helm-install
kubectl get all -n round-robin
```

Check:

```bash
curl http://app.127.0.0.1.nip.io/health
```

Expected:

- the stack now lives in namespace `round-robin`
- Helm controls the release lifecycle

What, why, when:

- raw manifests are good for understanding
- Helm is good when you need consistent packaging, values, upgrades, and rollbacks

## Phase 10: Practise Helm upgrade and rollback

Run:

```bash
helm upgrade round-robin infra/helm/round-robin -n round-robin --set app.replicaCount=2
kubectl get deploy -n round-robin
helm rollback round-robin -n round-robin
```

Then do a bad rollout:

```bash
./scripts/drills/break-image.sh round-robin round-robin-app app
kubectl rollout status deployment/round-robin-app -n round-robin
helm rollback round-robin -n round-robin
```

What you are learning:

- good change
- bad change
- controlled rollback

That is core production work.

## Phase 11: Understand Terraform before applying it

Read:

- [`infra/terraform/environments/local/main.tf`](/Users/deanlofts/gits/devops-roundrobin/infra/terraform/environments/local/main.tf)
- [`infra/terraform/modules/namespace/main.tf`](/Users/deanlofts/gits/devops-roundrobin/infra/terraform/modules/namespace/main.tf)
- [`infra/terraform/modules/platform/main.tf`](/Users/deanlofts/gits/devops-roundrobin/infra/terraform/modules/platform/main.tf)
- [`infra/terraform/modules/monitoring/main.tf`](/Users/deanlofts/gits/devops-roundrobin/infra/terraform/modules/monitoring/main.tf)

What, why, when:

- the environment root is where a specific environment is assembled
- child modules are where reusable units live
- you use this pattern when one repo needs local, dev, staging, and prod variants

## Phase 12: Run Terraform against the local cluster

Prepare:

```bash
cp infra/terraform/environments/local/terraform.tfvars.example \
  infra/terraform/environments/local/terraform.tfvars
```

Run:

```bash
make tf-prepare
make tf-init
make tf-plan
make tf-apply
```

Then inspect:

```bash
terraform -chdir=infra/terraform/environments/local output
```

What to pay attention to:

- what Terraform thinks it will create or change
- what is managed by Terraform versus manually
- that Terraform is now the owner of the `round-robin` release because the manual app release was removed first
- that `ingress-nginx` stays outside Terraform in this walkthrough because it was already installed as cluster bootstrap
- that this local root is intentionally managing namespace, Helm, and optional monitoring, not cluster creation

## Phase 13: Turn on monitoring

Run:

```bash
make monitor-install
kubectl get pods -n monitoring
```

Then open:

- `http://grafana.127.0.0.1.nip.io`
- `http://prometheus.127.0.0.1.nip.io`

Grafana login:

- user: `admin`
- password: `admin`

What you are learning:

- scraping
- dashboards
- release health from metrics, not just logs

## Phase 14: Run investigation drills

Run:

```bash
./scripts/drills/break-db-creds.sh round-robin round-robin-app
kubectl logs deployment/round-robin-app -n round-robin
kubectl get events -n round-robin --sort-by=.metadata.creationTimestamp
./scripts/drills/restore-db-creds.sh round-robin round-robin-app
```

Then:

```bash
kubectl scale deployment ingress-nginx-controller -n ingress-nginx --replicas=0
curl -i http://app.127.0.0.1.nip.io/health
kubectl scale deployment ingress-nginx-controller -n ingress-nginx --replicas=1
```

This is the closest part of the lab to real operational work.

## Phase 15: Rebuild from scratch

Run:

```bash
make tf-destroy
make kind-down
make compose-down
```

Then rebuild:

```bash
make kind-up
make ingress-install
make docker-build
make kind-load
make helm-install
```

Why it matters:

If you can rebuild cleanly, your understanding is real.

## End state

You are ready when you can do these without hesitation:

- explain the app, worker, Postgres, and Redis flow
- explain health versus readiness
- deploy with Compose
- deploy with raw Kubernetes
- deploy with Helm
- explain Terraform module structure
- run a rollback
- investigate a broken deployment with logs, events, readiness, and metrics

## If you get stuck

Read next:

- [Stuck Guide](/Users/deanlofts/gits/devops-roundrobin/docs/02-reference/01-when-stuck.md)
- [Concepts and Why](/Users/deanlofts/gits/devops-roundrobin/docs/02-reference/02-concepts-and-why.md)
- [Incident Playbooks](/Users/deanlofts/gits/devops-roundrobin/docs/03-incidents/01-common-devops-scenarios.md)
