# DevOps Round Robin Workbook

Local-first, realistic, incremental.

This repository is a progressive DevOps lab that takes you from an empty platform repo to a small production-style stack with an app, worker, PostgreSQL, Redis, Kubernetes, Helm, Terraform modules, observability, rollback drills, and full rebuild practice. Everything runs locally with Docker and kind. Every exercise also maps to the AWS equivalent so you rehearse the local mechanics and the cloud mental model at the same time.

The kind bootstrap defaults to an official `kindest/node` image pinned in `.env.example`. That pin is intentional. On some current Docker Desktop and Apple Silicon combinations, the newest default node image can fail during kubelet startup because of a containerd CRI compatibility issue. You can still override `KIND_NODE_IMAGE` when you want to test a newer Kubernetes version.

## What “best way possible” means here

This workbook intentionally uses current mainstream patterns rather than outdated shortcuts:

- `nix develop` for a pinned toolchain instead of relying on global installs
- Docker Compose v2 with `compose.yaml`
- multi-stage Dockerfiles with non-root distroless runtime images
- kind for local Kubernetes because it mirrors upstream Kubernetes behaviour closely
- Kubernetes `startupProbe`, `readinessProbe`, and `livenessProbe`
- `ingressClassName` instead of deprecated ingress annotations
- Helm `upgrade --install` workflow
- Terraform with a root module plus child modules, official `helm` and `kubernetes` providers, and environment-specific roots
- Prometheus Operator via `kube-prometheus-stack`

The local stack includes PostgreSQL and Redis inside containers or the cluster because you need something real to operate. In AWS, you would normally replace those with RDS and ElastiCache rather than self-managing them on EKS.

## Companion Docs

If you want a lower-distraction path with setup help, walkthroughs, recovery steps, and explanations, use these alongside the workbook:

- [Docs Index](docs/README.md)
- [MacBook Setup](docs/00-macbook-setup.md)
- [Start-to-Finish Walkthrough](docs/01-walkthrough/01-start-to-finish.md)
- [Checkpoint Questions](docs/01-walkthrough/02-checkpoint-questions.md)
- [What You Just Proved](docs/01-walkthrough/03-what-you-just-proved.md)
- [Exercises and Challenges](docs/01-walkthrough/04-exercises-and-challenges.md)
- [When You Get Stuck](docs/02-reference/01-when-stuck.md)
- [Concepts, What, Why, When](docs/02-reference/02-concepts-and-why.md)
- [Federated Identity and SSO](docs/02-reference/03-federated-identity-and-sso.md)
- [From URL to Rendered Page](docs/02-reference/04-from-url-to-rendered-page.md)
- [Cheatsheets](docs/02-reference/05-cheatsheets.md)
- [References and Awesome Lists](docs/02-reference/06-references-and-awesome-lists.md)
- [AWS Platform Mapping](docs/02-reference/07-aws-platform-mapping.md)
- [Common DevOps Scenarios](docs/03-incidents/01-common-devops-scenarios.md)

## Quick Start

```bash
cp .env.example .env
nix develop
make bootstrap
```

First local checkpoints:

```bash
make test
make docker-build
make compose-up
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl -X POST http://localhost:8080/jobs -H 'content-type: application/json' -d '{"payload":"hello"}'
curl http://localhost:8080/jobs
```

When you are done with the Docker phase:

```bash
make compose-down
make kind-up
make ingress-install
make kind-load
make k8s-apply-raw
```

When you are ready for Helm and Terraform:

```bash
make helm-lint
make helm-install
make tf-init
make tf-plan
make tf-apply
```

## Repo Shape

```text
.
├── app/
├── worker/
├── internal/
├── dashboards/
├── infra/
│   ├── k8s/
│   │   ├── base/
│   │   └── raw/
│   ├── helm/
│   │   └── round-robin/
│   └── terraform/
│       ├── environments/local/
│       └── modules/
├── scripts/
│   ├── dev/
│   ├── drills/
│   └── kind/
├── compose.yaml
├── flake.nix
├── Makefile
└── README.md
```

## Working Rules

- Do one task at a time.
- Do not skip verification.
- Do not move forward until the current task works.
- Break things on purpose where noted.
- Keep notes directly in this file or in your own journal.

For every task, capture:

1. What situation were you responding to?
2. What was your concrete task?
3. What actions did you take?
4. What was the result?
5. What is the AWS equivalent?
6. Where would you inspect it in AWS?
7. How would you know it is broken?
8. How would you recover it?

## Local Access Points

- App via Compose: [http://localhost:8080](http://localhost:8080)
- Worker via Compose: [http://localhost:18081](http://localhost:18081)
- App via Kubernetes ingress: [http://app.127.0.0.1.nip.io](http://app.127.0.0.1.nip.io)
- Grafana via monitoring task: [http://grafana.127.0.0.1.nip.io](http://grafana.127.0.0.1.nip.io)
- Prometheus via monitoring task: [http://prometheus.127.0.0.1.nip.io](http://prometheus.127.0.0.1.nip.io)

## Target Scenario

You are operating a small internal platform with:

- a web app
- a background worker
- PostgreSQL
- Redis
- Kubernetes
- Helm
- Terraform
- observability
- CI-style workflows
- rollout and rollback
- failure recovery

The web app exposes `/`, `/health`, `/ready`, `/metrics`, and `/jobs`. The worker consumes jobs from Redis and stores processed results in PostgreSQL. That gives you something concrete to build, deploy, break, debug, and recover.

---

# Task 1 - Create the Repo Skeleton

## STAR

Situation: You are starting from an almost-empty repository and need a shape that feels like a real platform repo.

Task: Build the basic structure so future work lands in clear places.

Action:

```bash
tree -a -L 3
```

Inspect:

- `app/`
- `worker/`
- `infra/k8s/`
- `infra/helm/`
- `infra/terraform/`
- `scripts/`
- `flake.nix`
- `Makefile`

Result: The repo now looks like something a platform team could actually extend.

AWS equivalent: Team repository and platform code structure.

Where in AWS: GitHub or CodeCommit, not the AWS Console.

Definition of done:

- repo structure exists
- `README.md` explains the scenario
- `flake.nix` and `Makefile` establish a repeatable entry point

---

# Task 2 - Build a Tiny Web App

## STAR

Situation: You need a service that behaves like something an orchestrator can operate.

Task: Run the app locally and verify the operational endpoints.

Action:

```bash
make app-run
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:8080/metrics
```

Expected behaviour:

- `/` returns service metadata
- `/health` stays green unless you force a health failure
- `/ready` reports dependency state
- `/metrics` exposes Prometheus metrics
- logs are structured JSON via `slog`

Result: You have an app with platform-friendly contracts before adding any orchestration.

AWS equivalent: A service running on ECS or EKS behind a load balancer.

Where in AWS:

- ECR for the image
- ECS or EKS for the runtime
- CloudWatch Logs for app logs

How you know it is broken:

- `/health` fails
- `/ready` shows dependency outages
- request logs show non-200 responses or startup errors

Recovery: Fix config, restore dependencies, or redeploy.

---

# Task 3 - Containerise the App

## STAR

Situation: The app works on your host but is not yet portable or deployable.

Task: Build and run the app as a non-root container.

Action:

```bash
make docker-build
docker run --rm -p 8080:8080 round-robin-app:local
```

Inspect:

- `app/Dockerfile`
- multi-stage build
- distroless runtime
- non-root user

Result: You have a container image that is closer to real delivery practices than “run it on my laptop”.

AWS equivalent: Build an image, push it to ECR, and run it in ECS or EKS.

Where in AWS:

- ECR image repository
- ECS task definition or EKS Deployment

---

# Task 4 - Build the Docker Compose Stack

## STAR

Situation: The app is containerised, but the platform dependencies do not exist yet.

Task: Run app, worker, PostgreSQL, and Redis together.

Action:

```bash
make compose-up
make compose-ps
curl http://localhost:8080/ready
curl -X POST http://localhost:8080/jobs -H 'content-type: application/json' -d '{"payload":"compose drill"}'
curl http://localhost:8080/jobs
```

Break it:

```bash
docker compose stop postgres
curl http://localhost:8080/ready
docker compose start postgres
```

Result: You can now practise service dependencies, startup order, and recovery on your own machine.

AWS equivalent: RDS + ElastiCache + app service in a VPC.

Where in AWS:

- RDS console and CloudWatch metrics
- ElastiCache console
- ECS service or EKS workload status

Recovery: Restore the dependency, confirm health checks, then re-run the workflow.

---

# Task 5 - Move Config to Environment Variables

## STAR

Situation: Hard-coded configuration becomes unmanageable once you run multiple environments.

Task: Understand and practise the environment-driven config model.

Action:

- Read `.env.example`
- Copy it to `.env`
- Change values deliberately and see how the app reacts

Useful drill:

```bash
docker compose down
cp .env.example .env
sed -n '1,200p' .env
make compose-up
```

Result: You are using the same shape of config injection you would use in CI, ECS, or Kubernetes.

AWS equivalent: SSM Parameter Store and Secrets Manager, injected into ECS tasks or Kubernetes workloads.

Where in AWS:

- Systems Manager Parameter Store
- Secrets Manager
- ECS task definition env and secret references
- EKS Secret and ConfigMap manifests

---

# Task 6 - Create a Local Kubernetes Cluster

## STAR

Situation: Compose proves local service wiring, but you still need cluster operations practice.

Task: Create a clean local Kubernetes cluster with kind.

Action:

```bash
make kind-up
kubectl config current-context
kubectl get nodes -o wide
```

What this does:

- creates a kind cluster
- wires a local registry
- opens host ports 80 and 443 for ingress

Result: You now have a disposable Kubernetes environment you can destroy and rebuild safely.

AWS equivalent: EKS cluster creation.

Where in AWS:

- EKS clusters
- EC2 or Fargate capacity behind the cluster

---

# Task 7 - Deploy with Raw Kubernetes Manifests

## STAR

Situation: You need to know what Helm and Terraform eventually abstract away.

Task: Apply the base manifests directly and inspect each object.

Action:

```bash
make kind-load
make k8s-apply-raw
kubectl get all -n round-robin-raw
kubectl describe deployment round-robin-app -n round-robin-raw
kubectl get configmap,secret,svc,ingress,pvc -n round-robin-raw
```

Break it:

```bash
kubectl -n round-robin-raw set env deployment/round-robin-app POSTGRES_HOST=wrong-host
kubectl -n round-robin-raw rollout status deployment/round-robin-app
kubectl -n round-robin-raw rollout undo deployment/round-robin-app
```

Result: You understand the workload model before packaging it.

AWS equivalent: Raw workloads on EKS.

Where in AWS:

- EKS workloads
- Kubernetes events and pod descriptions

---

# Task 8 - Add and Test Health Probes

## STAR

Situation: A scheduler needs more than “the process started”.

Task: Inspect how startup, readiness, and liveness interact.

Action:

```bash
kubectl describe deployment round-robin-app -n round-robin-raw
curl http://app.127.0.0.1.nip.io/health
curl http://app.127.0.0.1.nip.io/ready
```

Break it:

```bash
./scripts/drills/break-readiness-raw.sh
kubectl get pods -n round-robin-raw -w
./scripts/drills/restore-readiness-raw.sh
```

Result: You can see exactly how Kubernetes decides whether a pod is serving traffic or needs a restart.

AWS equivalent:

- ALB target health
- ECS container health checks
- EKS pod probes

Where in AWS:

- Target group health
- ECS service events
- EKS pod status and events

---

# Task 9 - Add Ingress

## STAR

Situation: Services exist inside the cluster, but nothing routes external traffic yet.

Task: Install ingress-nginx and route the app through it.

Action:

```bash
make ingress-install
kubectl get ingress -A
curl http://app.127.0.0.1.nip.io/
```

Result: You now have host-based routing into the cluster.

AWS equivalent: ALB plus Route53.

Where in AWS:

- EC2 Load Balancers
- Route53 hosted zones and records

---

# Task 10 - Package the Platform with Helm

## STAR

Situation: Raw manifests work, but managing them directly does not scale.

Task: Install the same platform through a chart.

Action:

```bash
kubectl delete -k infra/k8s/raw --ignore-not-found
make helm-lint
make helm-install
helm list -A
kubectl get all -n round-robin
curl http://app.127.0.0.1.nip.io/
```

Practice upgrade and rollback:

```bash
helm upgrade round-robin infra/helm/round-robin -n round-robin --set app.replicaCount=2
helm rollback round-robin -n round-robin
```

Result: You have moved from hand-managed objects to a repeatable application package.

AWS equivalent: Helm-managed app deployment on EKS.

Where in AWS:

- EKS workloads
- your GitOps or CI deployment history

---

# Task 11 - Add the Terraform Control Layer

## STAR

Situation: Helm manages an application release, but you still need infrastructure orchestration and lifecycle control.

Task: Hand ownership from manual Helm to Terraform cleanly, then use Terraform to manage the namespace and Helm release against the local cluster.

Action:

```bash
cp infra/terraform/environments/local/terraform.tfvars.example infra/terraform/environments/local/terraform.tfvars
make tf-prepare
make tf-init
make tf-plan
make tf-apply
terraform -chdir=infra/terraform/environments/local output
```

Important note:

- For this local lab, kind cluster creation stays outside Terraform on purpose.
- Terraform’s official Kubernetes and Helm providers are strongest when the cluster already exists.
- Terraform does not automatically adopt resources you created manually in earlier steps, including the app release.
- The walkthrough keeps `ingress-nginx` outside Terraform by default because it was already installed as cluster bootstrap in the Kubernetes phase.
- In a real AWS estate, you would usually choose one owner up front or import existing resources into Terraform state before switching ownership.
- In AWS, Terraform would also manage EKS, IAM, networking, and remote state.

Result: You now control the release lifecycle through Terraform instead of hand-running Helm.

AWS equivalent: Terraform provisioning for EKS add-ons and application releases.

Where in AWS:

- Terraform Cloud or CI logs
- EKS
- IAM
- VPC

---

# Task 12 - Structure Terraform with Real Modules

## STAR

Situation: A single large Terraform directory quickly turns into a maintenance problem.

Task: Understand the root module and child module split in this repo.

Action:

Inspect:

- `infra/terraform/environments/local/`
- `infra/terraform/modules/namespace/`
- `infra/terraform/modules/platform/`
- `infra/terraform/modules/monitoring/`

Focus on:

- input variables
- outputs
- provider boundaries
- local environment root calling reusable child modules

Result: You are working with a structure that scales much better than one huge `main.tf`.

AWS equivalent: Reusable Terraform modules for VPC, EKS, RDS, IAM, Route53, and platform services.

Where in AWS: Not a console location. This is about how infrastructure code is organised.

---

# Task 13 - Add the Background Worker

## STAR

Situation: Real platforms almost always have asynchronous work outside the request path.

Task: Verify the worker as a separate deployable unit and confirm end-to-end job flow.

Action:

Compose path:

```bash
make compose-up
curl -X POST http://localhost:8080/jobs -H 'content-type: application/json' -d '{"payload":"worker drill"}'
curl http://localhost:8080/jobs
```

Kubernetes path:

```bash
kubectl get deployments -n round-robin-raw
kubectl logs deployment/round-robin-worker -n round-robin-raw
curl -X POST http://app.127.0.0.1.nip.io/jobs -H 'content-type: application/json' -d '{"payload":"worker drill"}'
curl http://app.127.0.0.1.nip.io/jobs
```

Helm and Terraform path:

- worker is already part of the chart and Terraform module

Result: You can reason about multiple deployable units instead of a single monolith.

AWS equivalent:

- ECS service
- EKS Deployment
- worker backed by SQS or Redis queue

Where in AWS:

- ECS services
- EKS workloads
- SQS queues if you replace Redis with AWS-native queueing

---

# Task 14 - Practise Persistence

## STAR

Situation: Some state should survive restarts and some should not.

Task: Identify what persists in this lab and what the AWS equivalents would be.

Action:

Inspect:

- Compose named volumes for PostgreSQL and Redis
- PostgreSQL PVC in Kubernetes
- job history in `processed_jobs`

Drill:

```bash
docker compose down
docker volume ls | grep devops-roundrobin
make compose-up
curl http://localhost:8080/jobs
```

Result: You can distinguish ephemeral workloads from persistent data.

AWS equivalent:

- RDS for PostgreSQL
- ElastiCache for Redis
- EBS or EFS for stateful workloads
- S3 for object storage

Where in AWS:

- RDS
- ElastiCache
- EBS or EFS
- S3

---

# Task 15 - Add Metrics, Prometheus, and Grafana

## STAR

Situation: Logs tell you what happened once. Metrics tell you how the system is trending.

Task: Turn on the monitoring stack and inspect app and worker metrics.

Action:

Enable monitoring:

```bash
make monitor-install
kubectl get pods -n monitoring
```

Then visit:

- [http://grafana.127.0.0.1.nip.io](http://grafana.127.0.0.1.nip.io)
- [http://prometheus.127.0.0.1.nip.io](http://prometheus.127.0.0.1.nip.io)

Grafana credentials:

- username: `admin`
- password: `admin`

Result: You now have dashboards, scrape targets, and a real observability feedback loop.

AWS equivalent:

- CloudWatch metrics
- Amazon Managed Service for Prometheus
- Amazon Managed Grafana

Where in AWS:

- CloudWatch
- AMP
- AMG

---

# Task 16 - Practise Log-Driven Debugging

## STAR

Situation: Metrics show that something is wrong. You still need logs and events to find the cause.

Task: Use container logs, pod descriptions, and events to diagnose failures.

Action:

```bash
kubectl logs deployment/round-robin-app -n round-robin
kubectl logs deployment/round-robin-worker -n round-robin
kubectl describe pod -n round-robin
kubectl get events -n round-robin --sort-by=.metadata.creationTimestamp
```

Compose equivalent:

```bash
make compose-logs
```

Result: You can move from symptom to probable cause instead of guessing.

AWS equivalent: CloudWatch Logs plus ECS or EKS events.

Where in AWS:

- CloudWatch Logs
- ECS service events
- EKS pod events

---

# Task 17 - Use a CI-Style Workflow Locally

## STAR

Situation: Manual commands are fine for learning but not good enough for repeated delivery.

Task: Use the `Makefile` as your local CI workflow.

Action:

```bash
make test
make docker-build
make helm-lint
make tf-plan
```

Optional personal extension:

- add GitHub Actions later that run the same Make targets

Result: Your local workflow already resembles a small CI pipeline.

AWS equivalent:

- GitHub Actions
- CodeBuild
- CodePipeline

Where in AWS:

- CodePipeline
- CodeBuild

---

# Task 18 - Practise Scaling

## STAR

Situation: Stable deployments are not enough. You need predictable scaling behaviour.

Task: Turn on autoscaling and generate load.

Action:

Helm path:

```bash
helm upgrade round-robin infra/helm/round-robin -n round-robin --set autoscaling.enabled=true
make load-test
kubectl get hpa -n round-robin
kubectl top pods -n round-robin
```

Result: You practise requests, limits, HPA configuration, and load-driven behaviour.

AWS equivalent:

- EKS HPA
- ECS target-tracking autoscaling

Where in AWS:

- EKS workloads and metrics
- ECS service autoscaling

---

# Task 19 - Practise Rollback

## STAR

Situation: A bad deployment is inevitable. The important part is controlled recovery.

Task: Intentionally break a release and roll it back.

Action:

Helm rollback path:

```bash
./scripts/drills/break-image.sh round-robin round-robin-app app
kubectl rollout status deployment/round-robin-app -n round-robin
helm rollback round-robin -n round-robin
```

Raw Kubernetes rollback path:

```bash
kubectl -n round-robin-raw set image deployment/round-robin-app app=ghcr.io/example/does-not-exist:broken
kubectl -n round-robin-raw rollout undo deployment/round-robin-app
```

Result: You build muscle memory for controlled recovery instead of panic edits.

AWS equivalent:

- ECS deployment rollback
- Helm rollback on EKS

Where in AWS:

- ECS deployment history
- EKS release history through Helm or GitOps tooling

---

# Task 20 - Practise Security Basics

## STAR

Situation: Security hygiene must be built into the delivery path, not bolted on later.

Task: Review the repo for baseline container and secret practices.

Action:

Inspect:

- distroless non-root Dockerfiles
- secrets outside source code
- environment-driven config

Run:

```bash
make scan-image
```

Result: You have a baseline secure-by-default setup for a training platform.

AWS equivalent:

- IAM for runtime identity
- ECR vulnerability scanning
- Secrets Manager
- Security Hub and Inspector in larger environments

Where in AWS:

- ECR image scan results
- Secrets Manager
- IAM roles

---

# Task 21 - Run Failure Drills

## STAR

Situation: Knowing the happy path is not enough. Operations skill comes from rehearsed failure.

Task: Break the platform on purpose and recover it using metrics, logs, and events.

Action:

Database credential drill:

```bash
./scripts/drills/break-db-creds.sh round-robin round-robin-app
kubectl logs deployment/round-robin-app -n round-robin
./scripts/drills/restore-db-creds.sh round-robin round-robin-app
```

Ingress drill:

```bash
kubectl scale deployment ingress-nginx-controller -n ingress-nginx --replicas=0
curl http://app.127.0.0.1.nip.io/
kubectl scale deployment ingress-nginx-controller -n ingress-nginx --replicas=1
```

Probe drill:

```bash
./scripts/drills/break-readiness-raw.sh
kubectl get pods -n round-robin-raw
./scripts/drills/restore-readiness-raw.sh
```

Result: You start thinking like an operator, not just an implementer.

AWS equivalent: Real incident handling across ALB, ECS or EKS, RDS, and CloudWatch.

Where in AWS:

- CloudWatch dashboards and alarms
- ALB target health
- ECS or EKS events
- RDS metrics

---

# Task 22 - Destroy and Rebuild Everything

## STAR

Situation: The real test of infrastructure discipline is whether you can rebuild from scratch cleanly.

Task: Tear the lab down and bring it back up from nothing.

Action:

```bash
make tf-destroy
make kind-down
make compose-down
docker volume prune -f
make kind-up
make ingress-install
make docker-build
make kind-load
make tf-apply
```

Result: You prove that the platform is reproducible rather than hand-crafted.

AWS equivalent: Full environment recreation via Terraform and deployment automation.

Where in AWS:

- Terraform state and plan history
- EKS and dependent services

---

## Exit Criteria

You can:

- deploy the stack end to end locally
- explain every moving part in STAR format
- map every local task to the AWS equivalent
- use Helm and Terraform intentionally rather than mechanically
- debug failures with logs, metrics, events, and rollout history
- destroy and rebuild the environment without fear

## Best-Practice Notes for This Repo

- Use `nix develop` whenever possible so the toolchain stays reproducible.
- Treat Docker Compose as the service-dependency lab and kind as the orchestration lab.
- Learn the raw Kubernetes objects before relying on Helm.
- Use Terraform modules with small, clear ownership boundaries.
- Prefer `terraform plan` before `apply`.
- Prefer `helm upgrade --install` over bespoke scripting.
- Practise recovery drills while everything is fresh.

## Official References

- [kind documentation](https://kind.sigs.k8s.io/)
- [Kubernetes probes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [ingress-nginx deployment guide for kind](https://kubernetes.github.io/ingress-nginx/deploy/#kind)
- [Helm chart best practices](https://helm.sh/docs/chart_best_practices/)
- [Terraform module development guidance](https://developer.hashicorp.com/terraform/language/modules/develop)
- [Terraform Helm provider tutorial](https://developer.hashicorp.com/terraform/tutorials/kubernetes/helm-provider)
- [Prometheus client instrumentation for Go](https://prometheus.io/docs/guides/go-application/)
