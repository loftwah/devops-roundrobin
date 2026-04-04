# Concepts, What, Why, When

This is the explanation layer for the lab.

## App and worker

### What

The app accepts HTTP requests and enqueues jobs.

The worker consumes jobs and records processed results.

### Why

This mirrors a common production split:

- synchronous request path
- asynchronous background processing path

### When you would use this pattern

Use it when:

- user-facing requests should stay fast
- work can be retried
- tasks take longer than a normal request budget

AWS equivalent:

- ECS or EKS service for the app
- ECS or EKS worker
- SQS instead of Redis in many AWS-native designs

## Health versus readiness

### What

- health says whether the process should stay alive
- readiness says whether it should receive traffic

### Why

A service can be healthy but not ready.

Examples:

- process started, but DB is down
- process is fine, but warm-up is incomplete

### When

Use:

- liveness or health checks to restart stuck processes
- readiness checks to protect traffic from bad pods

## Docker Compose

### What

Compose is the local service orchestration layer in this lab.

### Why

It is the fastest way to validate:

- images
- dependency wiring
- ports
- env vars
- local persistence

### When

Use Compose before Kubernetes when you want to debug the application and dependency graph without cluster complexity.

## kind

### What

kind runs Kubernetes nodes as containers.

### Why

It gives you:

- a disposable cluster
- realistic Kubernetes APIs
- a local workflow with no cloud bill

### When

Use kind when you need to practise Kubernetes operations, manifests, ingress, and rollout behaviour.

## Raw Kubernetes manifests

### What

These are the primitive objects:

- Namespace
- ConfigMap
- Secret
- Service
- Deployment
- StatefulSet
- Ingress

### Why

You need to understand these before Helm and Terraform abstract them away.

### When

Use raw manifests when:

- learning
- debugging
- testing a minimal deploy path

## Helm

### What

Helm packages Kubernetes resources into a reusable chart with values.

### Why

It reduces duplication and makes install, upgrade, and rollback more repeatable.

### When

Use Helm when:

- the same workload exists in more than one environment
- you need configurable deployment values
- rollback and release history matter

Do not use Helm as an excuse not to understand the objects underneath it.

## Terraform modules

### What

Terraform modules are reusable infrastructure building blocks.

In this lab:

- root module: `infra/terraform/environments/local`
- child modules: `namespace`, `platform`, `monitoring`

### Why

Modules prevent:

- giant unstructured Terraform directories
- copy-paste drift
- unclear ownership

### When

Use modules when:

- the same pattern appears repeatedly
- you need environment-specific composition
- you want stable inputs and outputs

## Observability

### What

Observability here means:

- metrics
- logs
- events

### Why

You need all three:

- metrics show trend and scale
- logs show detail
- events show platform decisions

### When

Use metrics first to find symptoms.

Use logs and events to find causes.

## Rollback

### What

Rollback means returning to a known-good release.

### Why

Bad deployments happen. A safe rollback path is a core operational skill.

### When

Rollback when:

- a deploy is clearly bad
- impact is active
- diagnosis would take longer than safe recovery

Do not keep pushing speculative fixes into a broken rollout if rollback is available and faster.

## Why this lab maps well to AWS

Local here:

- app and worker containers
- Postgres and Redis containers
- kind cluster
- ingress-nginx
- Prometheus and Grafana

AWS there:

- ECS or EKS workloads
- RDS
- ElastiCache
- ALB and Route53
- CloudWatch, AMP, AMG
- Terraform-managed infra and add-ons

The implementation changes, but the operational thinking is the same.
