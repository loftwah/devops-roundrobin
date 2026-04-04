# Common DevOps Scenarios

These are realistic scenarios you can practise in this repo.

Use them like tickets or mini-incidents.

## Scenario 1: Service deploy completed but traffic is failing

### Situation

A new app deployment finished, but users are seeing failures.

### Common real-world symptoms

- ALB target group unhealthy
- ingress returns `503`
- deploy marked complete but application is unusable

### In this lab

Run:

```bash
./scripts/drills/break-readiness-raw.sh
curl -i http://app.127.0.0.1.nip.io/ready
kubectl describe deployment round-robin-app -n round-robin-raw
kubectl get endpoints -n round-robin-raw
```

### What you should do

1. confirm whether it is a health issue or a readiness issue
2. inspect deployment config
3. inspect endpoints and pod readiness
4. restore the change

Recover:

```bash
./scripts/drills/restore-readiness-raw.sh
```

### AWS equivalent

- EKS pod readiness
- ALB target health
- ECS task health if running in ECS instead

## Scenario 2: App cannot connect to the database after a change

### Situation

A config change went out and the app can no longer talk to Postgres.

### In this lab

Run:

```bash
./scripts/drills/break-db-creds.sh round-robin round-robin-app
kubectl logs deployment/round-robin-app -n round-robin
curl -i http://app.127.0.0.1.nip.io/ready
```

### What you should do

1. confirm the app is running but not ready
2. inspect logs for DB auth or connection errors
3. inspect deployment env vars or Secret references
4. restore the correct value

Recover:

```bash
./scripts/drills/restore-db-creds.sh round-robin round-robin-app
```

### AWS equivalent

- Secrets Manager rotation issue
- bad ECS task secret injection
- broken EKS Secret or environment reference
- RDS credential mismatch

## Scenario 3: Bad image shipped

### Situation

A deploy references an image that does not exist or cannot start.

### In this lab

Run:

```bash
./scripts/drills/break-image.sh round-robin round-robin-app app
kubectl get pods -n round-robin
kubectl describe deployment round-robin-app -n round-robin
kubectl get events -n round-robin --sort-by=.metadata.creationTimestamp
```

### What you should do

1. identify `ImagePullBackOff` or rollout failure
2. confirm the image reference is wrong
3. use rollback rather than improvising

Recover:

```bash
helm rollback round-robin -n round-robin
```

Or for raw Deployments:

```bash
./scripts/drills/rollback-deployment.sh round-robin round-robin-app
```

### AWS equivalent

- bad ECR image tag
- bad ECS task definition revision
- broken EKS Deployment image reference

## Scenario 4: Queue backlog or async processing issue

### Situation

Users can submit work but background processing is not keeping up.

### In this lab

Generate work:

```bash
for i in $(seq 1 20); do
  curl -s -X POST http://localhost:8080/jobs \
    -H 'content-type: application/json' \
    -d "{\"payload\":\"job-$i\"}" >/dev/null
done
```

Then inspect:

```bash
curl http://localhost:8080/jobs
docker compose logs worker
kubectl logs deployment/round-robin-worker -n round-robin
```

### What you should think about

- is the worker healthy?
- is Redis reachable?
- is Postgres reachable for persistence?
- do you need more worker replicas?

### AWS equivalent

- SQS backlog
- ECS worker scaling
- EKS worker scaling

## Scenario 5: Ingress or edge path failure

### Situation

The app is running, but traffic from the edge cannot reach it.

### In this lab

Run:

```bash
kubectl scale deployment ingress-nginx-controller -n ingress-nginx --replicas=0
curl -i http://app.127.0.0.1.nip.io/health
kubectl scale deployment ingress-nginx-controller -n ingress-nginx --replicas=1
```

### What you should do

1. separate edge failure from app failure
2. inspect ingress controller health
3. inspect ingress object and backend Service
4. restore the controller

### AWS equivalent

- ALB failure
- bad target group
- Route53 or ingress controller issue

## Scenario 6: Monitoring needs to explain the incident

### Situation

An alert fired, and you need to understand whether the issue is increasing errors, latency, or worker throughput.

### In this lab

Use:

- Grafana
- Prometheus
- app metrics
- worker metrics
- `kubectl logs`
- `kubectl get events`

### What you should do

1. identify the symptom from metrics
2. pivot into logs
3. confirm impact in platform events
4. recover and confirm metrics return to normal

### AWS equivalent

- CloudWatch dashboards and alarms
- AMP and AMG
- CloudWatch Logs

## Scenario 7: A full rebuild is required

### Situation

The environment is messy or you need to prove the platform is reproducible.

### In this lab

Run:

```bash
make tf-destroy
make kind-down
make compose-down
make kind-up
make ingress-install
make docker-build
make kind-load
make helm-install
```

### What this tests

- toolchain repeatability
- image build path
- cluster bootstrap
- deployment process

### AWS equivalent

- rebuild an environment from Terraform and CI

## How to use these scenarios for interview prep

For each one, practise saying:

1. what the symptom was
2. what layer you checked first
3. what data you used
4. how you recovered
5. what you would improve afterward

That maps extremely well to real-world DevOps conversations.
