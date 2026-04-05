# Exercises and Challenges

Use this when you want hands-on work that feels closer to a real DevOps job.

This file is not the reference walkthrough. It is the practice track.

The format is:

- situation
- your task
- constraints
- what to touch
- how to verify
- what "done" means

## How to use this

Work in order the first time.

For each exercise:

1. read the situation
2. try to solve it before looking everywhere
3. verify the result
4. write a short note: what broke, what you checked, what fixed it

Do not treat these like copy-paste labs. Treat them like tickets.

## Track 1: Operate the existing system

These exercises teach you how to run and inspect what already exists.

### Exercise 1: Orient yourself in the repo

Situation:

You joined a team and inherited this repository.

Your task:

- explain what each top-level directory is for
- explain which files are the main entry points

Constraints:

- do not run the app yet
- use the repo structure and docs first

What to touch:

- `README.md`
- `Makefile`
- `compose.yaml`
- `infra/`
- `scripts/`

Verify:

- you can explain the purpose of `app/`, `worker/`, `infra/k8s/`, `infra/helm/`, `infra/terraform/`, `scripts/`, `flake.nix`, and `Makefile`

Definition of done:

- you can describe the repo shape without looking at the tree output

### Exercise 2: Bring the app up without `make`

Situation:

The `Makefile` is convenient, but you need to know what it hides.

Your task:

- start the app directly
- verify its operational endpoints

Constraints:

- do not use `make app-run`

What to touch:

- `app/main.go`
- `internal/platform/config.go`

Verify:

```bash
go run ./app
curl -i http://localhost:8080/
curl -i http://localhost:8080/health
curl -i http://localhost:8080/ready
curl -i http://localhost:8080/metrics
```

Definition of done:

- you can explain why `/health` and `/ready` do not mean the same thing

### Exercise 3: Bring Compose up without `make`

Situation:

You need to understand the actual Compose command path.

Your task:

- start the local stack manually
- identify each running container and its role

Constraints:

- do not use `make compose-up`

What to touch:

- `compose.yaml`
- `.env`

Verify:

```bash
docker compose up --build -d
docker compose ps
curl -i http://localhost:8080/ready
curl -i http://localhost:18081/health
```

Definition of done:

- you can explain why the app becomes ready only after Postgres and Redis are available

### Exercise 4: Trace one job end to end

Situation:

A user submits work and expects it to be processed asynchronously.

Your task:

- submit a job
- prove where it enters
- prove where it is processed
- prove where the result ends up

Constraints:

- you must use both API output and logs

What to touch:

- app handler code
- worker processing code

Verify:

```bash
curl -X POST http://localhost:8080/jobs \
  -H 'content-type: application/json' \
  -d '{"payload":"exercise-job"}'

curl http://localhost:8080/jobs
docker compose logs --tail=50 app worker
```

Definition of done:

- you can explain the request path, queue path, and persistence path in plain language

## Track 2: Debug realistic failures

These exercises are closer to day-to-day operational work.

### Exercise 5: Dependency outage

Situation:

The app is up, but users report it is not usable.

Your task:

- break a dependency
- identify whether this is a health issue or a readiness issue
- recover it cleanly

Constraints:

- do not rebuild anything

Verify:

```bash
docker compose stop postgres
curl -i http://localhost:8080/health
curl -i http://localhost:8080/ready
docker compose logs --tail=50 app postgres
docker compose start postgres
curl -i http://localhost:8080/ready
```

Definition of done:

- you can explain why the process can stay alive while the service is not ready

### Exercise 6: Queue processing failure

Situation:

Jobs are accepted but are not being processed.

Your task:

- simulate or identify the failure
- determine whether the issue is with the app, worker, Redis, or Postgres

Constraints:

- do not guess from one signal
- use at least logs plus one runtime check

Verify:

- submit jobs
- inspect `docker compose logs worker`
- inspect `docker compose logs app`
- inspect readiness and job output

Definition of done:

- you can state which layer failed and why

### Exercise 7: Readiness failure in Kubernetes

Situation:

A deployment succeeded, but traffic through ingress is failing.

Your task:

- reproduce a readiness-related failure
- inspect the rollout
- restore service

Constraints:

- use Kubernetes evidence, not only curl output

What to touch:

- `scripts/drills/break-readiness-raw.sh`
- `scripts/drills/restore-readiness-raw.sh`

Verify:

```bash
./scripts/drills/break-readiness-raw.sh
kubectl get pods -n round-robin-raw
kubectl describe deployment round-robin-app -n round-robin-raw
kubectl get endpoints -n round-robin-raw
curl -i http://app.127.0.0.1.nip.io/ready
./scripts/drills/restore-readiness-raw.sh
```

Definition of done:

- you can explain why ingress can fail even when pods still exist

### Exercise 8: Bad image rollout

Situation:

A deployment points at an image that does not exist or will not start.

Your task:

- detect the failure
- identify the bad image reference
- recover with rollback

Constraints:

- do not improvise a random fix
- use the platform rollback path

Verify:

```bash
./scripts/drills/break-image.sh round-robin round-robin-app app
kubectl get pods -n round-robin
kubectl get events -n round-robin --sort-by=.metadata.creationTimestamp
helm history round-robin -n round-robin
helm rollback round-robin -n round-robin
```

Definition of done:

- you can explain what `ImagePullBackOff` means operationally

## Track 3: Change the platform safely

These exercises are about controlled changes, which is a large part of the job.

### Exercise 9: Inspect a container build properly

Situation:

You are asked whether the app image is production-friendly.

Your task:

- inspect the app Dockerfile
- explain the build stages
- explain why the runtime image choice matters

What to touch:

- `app/Dockerfile`

Verify:

- build the image
- run the image
- explain why it runs as it does

Definition of done:

- you can explain multi-stage builds and non-root runtime choices without hand-waving

### Exercise 10: Install the Helm release manually

Situation:

You need to understand the real Helm command, not just `make helm-install`.

Your task:

- lint the chart
- install or upgrade the release manually
- inspect the resulting resources

Constraints:

- do not use `make helm-install`

What to touch:

- `infra/helm/round-robin/Chart.yaml`
- `infra/helm/round-robin/values.yaml`
- `infra/helm/round-robin/templates/`

Verify:

```bash
helm lint infra/helm/round-robin
helm upgrade --install round-robin infra/helm/round-robin \
  --namespace round-robin \
  --create-namespace \
  -f infra/helm/round-robin/values.yaml
kubectl get all -n round-robin
```

Definition of done:

- you can explain how values become rendered manifests

### Exercise 11: Read a Terraform plan like an operator

Situation:

Terraform is about to change the environment and you need to understand what it will do before apply.

Your task:

- initialize the local Terraform root
- create a plan
- explain what Terraform owns here

Constraints:

- do not apply before you can explain the plan

What to touch:

- `infra/terraform/environments/local/main.tf`
- `infra/terraform/modules/`

Verify:

```bash
terraform -chdir=infra/terraform/environments/local init
terraform -chdir=infra/terraform/environments/local plan -out=tfplan
terraform -chdir=infra/terraform/environments/local show tfplan
```

Definition of done:

- you can explain what Terraform will create, update, or own without guessing

### Exercise 12: Rebuild from scratch

Situation:

The environment is messy and you need to prove the platform is reproducible.

Your task:

- tear everything down cleanly
- rebuild it in the documented order
- verify the stack works again

Constraints:

- do not rely on hidden state
- do not skip verification

Verify:

```bash
make tf-destroy
make kind-down
make compose-down
docker ps
make kind-up
make ingress-install
make docker-build
make kind-load
make helm-install
curl -i http://app.127.0.0.1.nip.io/health
```

Definition of done:

- the rebuild succeeds from a clean state

## Track 4: Interview and job-readiness challenges

These are not command-only exercises. They test whether you can think and explain like an engineer.

### Challenge 1: Explain the system from memory

Without looking at the repo, explain:

- what the app does
- what the worker does
- where the queue lives
- where processed results are stored
- why health and readiness are separate
- why the repo has raw Kubernetes, Helm, and Terraform

### Challenge 2: Explain a failure investigation

Pick one failure and explain:

- symptom
- likely layers involved
- first checks
- confirming evidence
- recovery path

Use one of:

- Postgres outage
- Redis outage
- bad image
- ingress failure
- readiness failure

### Challenge 3: Explain the AWS mapping

Explain the local-to-AWS mapping for:

- Compose
- `kind`
- ingress
- Postgres
- Redis
- Helm
- Terraform
- monitoring

### Challenge 4: Explain ownership boundaries

Explain who owns what here:

- the app code
- the container build
- Compose
- raw Kubernetes manifests
- Helm chart
- Terraform root

If you cannot explain ownership, you will struggle to change the system safely.

## Suggested order

First pass:

1. Exercises 1-4
2. Exercises 5-8
3. Exercises 9-12
4. Challenges 1-4

Second pass:

- redo the same exercises with less reliance on docs
- explain each result out loud as if you were onboarding a teammate

## What success looks like

This doc is working if you become able to do three things:

1. run the platform
2. debug the platform
3. explain the platform

That is much closer to the real job than simply repeating commands.
