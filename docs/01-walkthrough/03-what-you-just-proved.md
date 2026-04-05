# What You Just Proved

Use this with the walkthrough when you want each step to feel like real learning instead of command repetition.

This file answers the question:

"I ran the command. What did that actually prove, what should I read next, and what do I need to understand before moving on?"

## How to use this

For each walkthrough phase:

1. run the step
2. read the matching section here
3. inspect the listed files
4. do not move on until you can explain the step in your own words

If you can only say "the command worked", you are not done with that phase yet.

## Phase 0: Prepare the machine

### `cp .env.example .env`

What you just did:

- created the local config file the repo actually reads
- kept the tracked template file unchanged

What this proves:

- the repo now has a concrete local configuration source
- future commands can read expected environment values

What to read next:

- [`.env.example`](../../.env.example)
- [`Makefile`](../../Makefile)

What to understand before moving on:

- why `.env.example` is versioned but `.env` is local
- which values are machine-specific
- which values affect app behavior versus infrastructure wiring

If this were broken, look here first:

- whether `.env` exists
- whether the values make sense for your machine

### `nix develop`

What you just did:

- entered the pinned toolchain for this repo

What this proves:

- your shell now has the repo’s expected CLI tools on `PATH`
- you are using the repo’s chosen versions, not random global ones

What to read next:

- [`flake.nix`](../../flake.nix)

What to understand before moving on:

- why this repo uses Nix instead of assuming global installs
- which tools are coming from the Nix shell
- what version drift would look like without this

If this were broken, look here first:

- Nix itself
- whether you are actually inside the dev shell

### `make bootstrap`

What you just did:

- ran the first guardrail before touching the stack

What this proves:

- `.env` exists
- required commands are available in the current shell

What to read next:

- [`Makefile`](../../Makefile)
- [`scripts/dev/check-prereqs.sh`](../../scripts/dev/check-prereqs.sh)

What to understand before moving on:

- what `bootstrap` actually checks
- why failing early is better than finding missing tools halfway through

If this were broken, look here first:

- missing `.env`
- missing CLI tools

## Phase 1: Understand the repo shape

What you just did:

- mapped the repo structure before running more infrastructure

What this proves:

- you can orient yourself in the codebase
- you know where the major layers live

What to read next:

- [`README.md`](../../README.md)
- [`Makefile`](../../Makefile)
- [`compose.yaml`](../../compose.yaml)

What to understand before moving on:

- `app/` is the web service
- `worker/` is the background processor
- `infra/k8s/` is raw Kubernetes
- `infra/helm/` is packaged Kubernetes
- `infra/terraform/` is the higher-level control layer
- `scripts/` contains operational helpers

If this were broken, look here first:

- the repo structure itself is not the issue
- the issue is usually that you moved on without knowing where to inspect next

## Phase 2: Run the app alone

What you just did:

- started only the app process without its dependencies

What this proves:

- the app binary can start
- the app exposes its operational endpoints
- the app distinguishes process health from dependency readiness

What to read next:

- [`app/`](../../app)
- [`internal/platform/config.go`](../../internal/platform/config.go)

What to understand before moving on:

- `/health` answers whether the process should be restarted
- `/ready` answers whether the service should receive traffic
- `/metrics` exposes observability data
- the app can be alive without being useful yet

If this were broken, look here first:

- whether the app process started
- whether port `8080` is free
- whether you are in the right shell

## Phase 3: Run the full local stack with Compose

What you just did:

- brought up the app, worker, Postgres, and Redis together

What this proves:

- service-to-service connectivity works in the local container environment
- the app becomes ready when its dependencies exist
- the worker has the dependencies it needs

What to read next:

- [`compose.yaml`](../../compose.yaml)
- [`worker/`](../../worker)

What to understand before moving on:

- which containers exist and why
- which ports are exposed to your Mac
- why the worker gets its own health endpoint
- how Compose is wiring service names like `postgres` and `redis`

If this were broken, look here first:

- `docker compose ps`
- `docker compose logs`
- host port conflicts

## Phase 4: Exercise the async workflow

What you just did:

- pushed work into the system and confirmed the worker processed it

What this proves:

- the app accepts jobs
- the worker consumes jobs
- Redis is working as the queue layer
- Postgres is working as the persistence layer

What to read next:

- app job handlers
- worker job processing code

What to understand before moving on:

- which component accepts the request
- which component performs background work
- where the queued state lives
- where the processed result ends up

If this were broken, look here first:

- app logs
- worker logs
- Redis health
- Postgres health

## Phase 5: Break a dependency and recover

What you just did:

- forced a partial outage and recovered from it

What this proves:

- the app can remain alive while becoming unready
- readiness reflects downstream dependency health
- recovery can happen without rebuilding the service

What to read next:

- readiness-related code in the app
- Compose logs during the failure window

What to understand before moving on:

- why health and readiness are different signals
- why "process up" does not mean "service usable"
- why logs are necessary to explain the failing dependency

If this were broken, look here first:

- whether the dependency really stopped
- whether readiness actually checks that dependency

## Phase 6: Move to Kubernetes

What you just did:

- left local Compose orchestration and created a disposable Kubernetes environment

What this proves:

- you can run a local cluster
- ingress prerequisites are in place
- the repo can prepare local image delivery for the cluster

What to read next:

- [`scripts/kind/create-cluster.sh`](../../scripts/kind/create-cluster.sh)
- [`scripts/kind/load-images.sh`](../../scripts/kind/load-images.sh)

What to understand before moving on:

- what `kind` is doing here
- why the repo uses a local registry
- why host ports `80` and `443` are mapped
- why image loading is needed before the cluster can run your app

If this were broken, look here first:

- Docker runtime health
- `kind` cluster status
- ingress controller pods

## Phase 7: Deploy raw Kubernetes resources

What you just did:

- deployed the app stack as primitive Kubernetes objects

What this proves:

- the raw manifests are valid enough to run the full system
- the cluster can schedule the app, worker, Postgres, and Redis
- ingress can route traffic to the app

What to read next:

- [`infra/k8s/raw/`](../../infra/k8s/raw)
- [`infra/k8s/base/`](../../infra/k8s/base)

What to understand before moving on:

- how Deployments, Services, ConfigMaps, Secrets, PVCs, and Ingress work together
- why the raw layer exists before Helm

If this were broken, look here first:

- pod status
- events
- service endpoints
- ingress resources

## Phase 8: Debug a raw Kubernetes incident

What you just did:

- practised debugging Kubernetes state instead of just applying manifests

What this proves:

- you can use rollout status, deployment description, and external symptoms together
- you can restore the deployment after breaking it

What to read next:

- the drill scripts in [`scripts/drills/`](../../scripts/drills)

What to understand before moving on:

- how readiness failures surface in Kubernetes
- what `kubectl describe` adds beyond `kubectl get`
- how to separate symptom from cause

If this were broken, look here first:

- deployment status
- pod events
- readiness endpoint behavior

## Phase 9: Replace raw resources with Helm

What you just did:

- switched from primitive manifests to release packaging

What this proves:

- the same stack can be managed as a Helm release
- the app lifecycle can now be handled through chart values and release history

What to read next:

- [`infra/helm/round-robin/Chart.yaml`](../../infra/helm/round-robin/Chart.yaml)
- [`infra/helm/round-robin/values.yaml`](../../infra/helm/round-robin/values.yaml)
- [`infra/helm/round-robin/templates/`](../../infra/helm/round-robin/templates)

What to understand before moving on:

- what Helm adds over raw manifests
- how values become rendered Kubernetes objects
- why namespace and naming now differ from the raw phase

If this were broken, look here first:

- `helm lint`
- rendered templates
- release history

## Phase 10: Practise Helm rollback

What you just did:

- exercised both a good change and a bad change

What this proves:

- you can use Helm as an operational control point
- rollback is part of the normal workflow, not just an emergency idea

What to read next:

- Helm history and rollback behavior in the running cluster

What to understand before moving on:

- what a successful upgrade looks like
- what a broken rollout looks like
- how rollback restores the previously good state

If this were broken, look here first:

- deployment rollout status
- Helm history
- current rendered values

## Phase 11: Understand Terraform before applying it

What you just did:

- paused before applying Terraform to understand its ownership boundaries

What this proves:

- you are not treating Terraform like a magic black box

What to read next:

- [`infra/terraform/environments/local/main.tf`](../../infra/terraform/environments/local/main.tf)
- [`infra/terraform/modules/namespace/main.tf`](../../infra/terraform/modules/namespace/main.tf)
- [`infra/terraform/modules/platform/main.tf`](../../infra/terraform/modules/platform/main.tf)
- [`infra/terraform/modules/monitoring/main.tf`](../../infra/terraform/modules/monitoring/main.tf)

What to understand before moving on:

- what the root module assembles
- what the child modules encapsulate
- why local cluster creation is not owned here

If this were broken, look here first:

- you likely skipped understanding ownership and will be confused about what Terraform should control

## Phase 12: Run Terraform

What you just did:

- handed cluster application ownership to Terraform

What this proves:

- Terraform can initialize, plan, and apply against the running local cluster
- the Helm release can be managed through Terraform instead of only manually

What to read next:

- Terraform state and outputs in the local environment directory

What to understand before moving on:

- why `plan` matters before `apply`
- what Terraform owns versus what Helm owns
- why the manual Helm release is removed before Terraform takes over

If this were broken, look here first:

- provider init
- kubeconfig and cluster access
- ownership conflicts with an already-installed release

## Phase 13: Turn on monitoring

What you just did:

- added observability tooling on top of the running platform

What this proves:

- the stack can expose dashboards and metrics endpoints through ingress
- you now have more than logs for diagnosis

What to read next:

- monitoring-related Terraform and ingress configuration
- the Grafana dashboard JSON files

What to understand before moving on:

- why metrics complement logs
- what Prometheus and Grafana each do
- how you would use them during an incident

If this were broken, look here first:

- monitoring pods
- monitoring ingress
- ingress controller health

## Phase 14: Investigation drills

What you just did:

- practised failure analysis with deliberate breakage

What this proves:

- you can form a debugging loop instead of guessing
- you can recover after identifying the right layer

What to read next:

- incident drills and the stuck guide

What to understand before moving on:

- which signals tell you app failure versus config failure versus routing failure
- why logs, events, readiness, and ingress all matter

If this were broken, look here first:

- the first signal that changed after you injected failure

## Phase 15: Rebuild from scratch

What you just did:

- tore the environment down and rebuilt it from nothing

What this proves:

- the workflow is repeatable
- the infrastructure is not dependent on hidden manual steps
- your understanding is strong enough to reconstruct the system

What to read next:

- teardown and bootstrap scripts

What to understand before moving on:

- why clean teardown matters
- why rebuild is the strongest confidence check in the repo

If this were broken, look here first:

- leftover Docker state
- leftover cluster state
- steps executed out of order

## Final standard

You understand a phase when you can answer all four of these without guessing:

1. What did this command create, change, or verify?
2. Which file defines that behavior?
3. What would I expect to see if it were broken?
4. Where would I look first?

If you cannot answer those, go back one phase and inspect the files before continuing.
