# Checkpoint Questions

Use this alongside the walkthrough when you want quick confirmation that you did the step properly.

This is not a theory quiz. The point is to tie each action to:

- what changed
- what you should expect to see
- how to verify it quickly
- what it probably means if the result is wrong

## How to use this

At each checkpoint:

1. run the command from the walkthrough
2. answer the questions without guessing
3. run the verification command
4. do not move on until the answer and the observed result line up

If you want to keep notes, copy the prompts into your own scratch file and write short answers under each one.

## Phase 0: Prepare the machine

### `cp .env.example .env`

Questions:

- Why do we copy `.env.example` instead of editing `.env.example` directly?
- What is the difference between a template env file and the env file the repo actually reads?
- Which values in `.env` look like local-machine choices rather than app logic?

Expected:

- a new `.env` file exists at the repo root
- it starts with the same values as `.env.example`
- the repo can now read local config without modifying the tracked template

Verify:

```bash
ls -l .env .env.example
diff -u .env.example .env
```

Green flags:

- both files exist
- `diff` shows no output yet

If not:

- `.env` is missing: copy it again
- `diff` shows changes: make sure they were intentional

### `nix develop`

Questions:

- Why is `nix develop` the entrypoint instead of relying on globally installed tools?
- Which tools does this repo expect to come from the Nix shell?
- What problem are we avoiding by pinning the toolchain?

Expected:

- you enter a shell where repo tooling is available on `PATH`
- the shell prints the repo-specific message from the Nix shell hook

Verify:

```bash
go version
kind version
kubectl version --client
helm version --short
terraform version
```

Green flags:

- the commands run without `command not found`
- versions print successfully inside the shell

If not:

- you are probably outside `nix develop`
- or Nix failed to provide the shell correctly

### `make bootstrap`

Questions:

- What does `bootstrap` actually check in this repo?
- Why is it useful to fail here instead of later?
- Which missing command would hurt you much later if this check did not exist?

Expected:

- `.env` is verified first
- required local commands are checked
- you see `Local prerequisites are available.`

Verify:

```bash
make bootstrap
```

Green flags:

- the command exits successfully
- the output is short and explicit

If not:

- read the failure literally
- fix the missing file or command before continuing

## Phase 1: Understand the repo shape

### `find . -maxdepth 3 -type f | sort`

Questions:

- Where does app code live?
- Where does worker code live?
- Where do Docker and Compose live?
- Where are raw Kubernetes manifests, Helm, and Terraform separated?

Expected:

- you can point to the main directories without hesitation
- you can explain why this is one repo with multiple operational layers

Verify:

```bash
find . -maxdepth 3 -type f | sort
```

Green flags:

- you can identify `app/`, `worker/`, `infra/k8s/`, `infra/helm/`, `infra/terraform/`, `scripts/`, `Makefile`, and `flake.nix`

If not:

- reread `README.md`, `Makefile`, and `compose.yaml`

## Phase 2: Run the app alone

### `make app-run`

Questions:

- What can the app do by itself before dependencies exist?
- Why should `/health` succeed while `/ready` can still fail?
- What does that tell you about health versus readiness?

Expected:

- the app starts on port `8080`
- `/health` returns `200`
- `/ready` returns `503` until Postgres and Redis exist
- `/metrics` responds

Verify:

```bash
curl http://localhost:8080/
curl -i http://localhost:8080/health
curl -i http://localhost:8080/ready
curl http://localhost:8080/metrics
```

Green flags:

- health is up
- readiness is not yet up
- metrics are exposed

If not:

- the app may not be running
- or the port is already taken

## Phase 3: Run the full local stack with Compose

### `make compose-up`

Questions:

- Which services should start here?
- What new dependencies now exist that did not exist in Phase 2?
- Why is Compose a useful step before Kubernetes?

Expected:

- app, worker, Postgres, and Redis containers are running
- app readiness becomes healthy
- worker health is available on `localhost:18081`

Verify:

```bash
make compose-ps
curl -i http://localhost:8080/health
curl -i http://localhost:8080/ready
curl -i http://localhost:18081/health
docker compose ps
```

Green flags:

- app is healthy and ready
- worker is healthy
- Postgres and Redis show healthy in Compose

If not:

- inspect `docker compose logs --tail=50 app postgres redis worker`

## Phase 4: Exercise the async workflow

### `POST /jobs` then `GET /jobs`

Questions:

- Which component accepts the job?
- Which component processes it?
- Which services are involved between submission and persistence?

Expected:

- `POST /jobs` returns `202`
- `GET /jobs` shows processed work

Verify:

```bash
curl -X POST http://localhost:8080/jobs \
  -H 'content-type: application/json' \
  -d '{"payload":"walkthrough job"}'

curl http://localhost:8080/jobs
```

Green flags:

- the app accepts the job
- the worker eventually processes it

If not:

- check app logs, worker logs, Redis availability, and Postgres availability

## Phase 5: Break a dependency and recover

### `docker compose stop postgres`

Questions:

- What should fail first: health or readiness?
- Why is that distinction operationally useful?
- What evidence should appear in logs?

Expected:

- app process keeps running
- readiness drops
- logs show dependency failure

Verify:

```bash
docker compose stop postgres
curl -i http://localhost:8080/ready
docker compose logs --tail=50 app postgres
```

Green flags:

- `/ready` degrades
- the process does not necessarily crash

If not:

- make sure you actually stopped Postgres

Recovery:

```bash
docker compose start postgres
curl -i http://localhost:8080/ready
```

## Phase 6: Move to Kubernetes

### `make kind-up`

Questions:

- What does `kind-up` create besides the cluster itself?
- Why does this repo wire a local registry?
- Why do ports `80` and `443` matter here?

Expected:

- the `kind` cluster exists
- the node becomes Ready
- the local `kind-registry` exists

Verify:

```bash
kubectl get nodes -o wide
docker ps
```

Green flags:

- one Ready control-plane node
- `kind-registry` exists

If not:

- inspect Docker runtime health and `kind get clusters`

### `make ingress-install`

Questions:

- What problem does ingress solve in this lab?
- Why install it before the app manifests?

Expected:

- ingress controller pods are running in `ingress-nginx`

Verify:

```bash
kubectl get pods -n ingress-nginx
```

## Phase 7: Deploy raw Kubernetes resources

### `make k8s-apply-raw`

Questions:

- Which primitives are you exercising here that Helm will later package?
- Which namespace should these resources land in?
- How do Service, Deployment, Secret, ConfigMap, PVC, and Ingress fit together?

Expected:

- raw resources appear in `round-robin-raw`
- ingress serves the app

Verify:

```bash
kubectl get all -n round-robin-raw
kubectl get ingress -n round-robin-raw
curl -i http://app.127.0.0.1.nip.io/health
curl -i http://app.127.0.0.1.nip.io/ready
```

Green flags:

- pods become Running
- ingress exists
- app responds through ingress

## Phase 8: Debug a raw Kubernetes incident

### Break and restore readiness

Questions:

- Which signals tell you the rollout is unhealthy?
- Where would you look first: pod status, deployment description, logs, events, or the app endpoint?
- What does a readiness failure look like from outside the cluster?

Expected:

- rollout health degrades
- readiness endpoint fails
- restore returns the deployment to healthy state

Verify:

```bash
kubectl get pods -n round-robin-raw
kubectl describe deployment round-robin-app -n round-robin-raw
curl -i http://app.127.0.0.1.nip.io/ready
```

## Phase 9: Replace raw resources with Helm

### `make helm-install`

Questions:

- What changes when Helm becomes the control layer?
- Which namespace should now contain the app stack?
- Why is Helm better than raw manifests for repeated installs and rollbacks?

Expected:

- resources now live in `round-robin`
- Helm controls the release lifecycle

Verify:

```bash
helm list -A
kubectl get all -n round-robin
curl -i http://app.127.0.0.1.nip.io/health
```

## Phase 10: Practise Helm rollback

Questions:

- What is the difference between a successful change and a bad rollout?
- How do you know rollback actually restored service rather than just changing Helm metadata?

Expected:

- good upgrade works
- bad rollout fails visibly
- rollback restores the working release

Verify:

```bash
helm history round-robin -n round-robin
kubectl rollout status deployment/round-robin-app -n round-robin
curl -i http://app.127.0.0.1.nip.io/health
```

## Phase 11: Understand Terraform before applying it

Questions:

- Which file is the environment root?
- Which files are reusable child modules?
- Why keep that separation?

Expected:

- you can explain root module versus child module without hand-waving

Verify:

- read the Terraform files listed in the walkthrough
- explain in one sentence what each one is responsible for

## Phase 12: Run Terraform

### `make tf-init`, `make tf-plan`, `make tf-apply`

Questions:

- What will Terraform manage here?
- What is deliberately still outside Terraform in this walkthrough?
- Why do we remove the manual Helm release before handing ownership to Terraform?

Expected:

- `plan` shows intended changes before apply
- `apply` creates the managed resources successfully

Verify:

```bash
make tf-plan
make tf-apply
terraform -chdir=infra/terraform/environments/local output
```

Green flags:

- the plan is readable
- apply succeeds
- outputs render after apply

## Phase 13: Turn on monitoring

Questions:

- What changes once metrics and dashboards exist?
- Why are logs alone not enough?

Expected:

- monitoring pods run
- Grafana and Prometheus become reachable through ingress

Verify:

```bash
kubectl get pods -n monitoring
curl -I http://grafana.127.0.0.1.nip.io
curl -I http://prometheus.127.0.0.1.nip.io
```

## Phase 14: Investigation drills

Questions:

- Which tools tell you configuration failure versus networking failure versus app failure?
- Why inspect logs, events, readiness, and ingress together instead of only one of them?

Expected:

- you can break a component, gather evidence, and restore it deliberately

Verify:

- run the drill
- write down what signal told you the truth first

## Phase 15: Rebuild from scratch

Questions:

- Can you tear down the stack cleanly?
- Can you rebuild it without improvising new steps?
- If rebuild succeeds, what does that say about your understanding?

Expected:

- Compose is down
- `kind` cluster is gone
- `kind-registry` is removed
- rebuild works again from nothing

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
```

Green flags:

- teardown leaves no lab containers behind except unrelated local services
- rebuild succeeds in the documented order

## Final confidence check

You should be able to answer these out loud without looking:

- Why do we have both `.env.example` and `.env`?
- Why do we use `nix develop` before anything else?
- What does `make bootstrap` protect us from?
- Why is `/health` different from `/ready`?
- What is Compose teaching that Kubernetes later builds on?
- Why does `kind` need a local registry in this repo?
- Why learn raw manifests before Helm?
- Why hand Helm ownership to Terraform only after removing the manual release?
- What would you inspect first if the app returned `503` through ingress?

If you can answer those clearly and verify each phase without guessing, you are not just following commands anymore. You understand the system.
