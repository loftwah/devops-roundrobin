# MacBook Setup

This guide is the lowest-distraction path through the lab on a Mac.

## Goal

End up with:

- one container runtime
- Nix installed
- no need to install `kubectl`, `helm`, `terraform`, `kind`, `jq`, `yq`, `trivy`, `k6`, `go`, `postgres`, or `redis` globally

## Recommended setup

Use:

- Nix for all CLI tooling in this repo
- Docker Desktop or OrbStack for container runtime

That is the cleanest path because:

- the toolchain is pinned in [`flake.nix`](/Users/deanlofts/gits/devops-roundrobin/flake.nix)
- you avoid version drift
- you avoid global installs fighting each other
- the repo becomes portable between machines

## What must exist outside Nix

You still need these installed on the Mac itself:

1. Nix
2. a Docker-compatible runtime

Everything else for this lab is supplied by `nix develop`.

## Recommended order

1. Ensure Nix works:

```bash
nix --version
```

2. Ensure your container runtime works:

```bash
docker --version
docker info
docker compose version
```

3. Clone or open the repo and create the env file:

```bash
cp .env.example .env
```

4. Enter the pinned toolchain shell:

```bash
nix develop
```

5. Verify the repo prerequisites:

```bash
make bootstrap
```

## Expected result

Inside `nix develop`, these should work:

```bash
go version
kind version
kubectl version --client
helm version --short
terraform version
jq --version
yq --version
trivy --version
k6 version
```

## Why this setup is the best fit here

### Why Nix

Use Nix when:

- you want a reproducible local toolchain
- you are switching between repos with different tool versions
- you want to avoid Homebrew sprawl

Do not use Nix here because it is fashionable. Use it because this workbook includes many CLI tools and you do not want to debug installation drift a week before a new job.

### Why Docker Desktop or OrbStack

Use a Docker-compatible runtime when:

- you need Docker builds
- you need Docker Compose
- you need kind, which runs Kubernetes nodes as containers

## Known local choices in this repo

### Worker host port

The worker is mapped to `localhost:18081`, not `8081`.

Why:

- `8081` is commonly already taken on Macs by other local tooling
- this avoids an unnecessary early failure

### Postgres and Redis host ports

The repo defaults to:

- `localhost:5432` for Postgres
- `localhost:6379` for Redis

If your Mac already has either service running, override these in [`.env.example`](/Users/deanlofts/gits/devops-roundrobin/.env.example) after copying it to `.env`:

```bash
POSTGRES_HOST_PORT=15432
REDIS_HOST_PORT=16379
```

### kind node image pin

The repo pins `KIND_NODE_IMAGE` in [`.env.example`](/Users/deanlofts/gits/devops-roundrobin/.env.example).

Why:

- on this class of Mac and Docker setup, the newest default kind node image can fail during kubelet startup because of a containerd CRI compatibility problem
- a pinned official kind node image gives you a predictable path through the lab

Override it only if you deliberately want to test a newer Kubernetes version.

## Minimal daily workflow

When you come back to the repo later:

```bash
cd /Users/deanlofts/gits/devops-roundrobin
nix develop
make help
```

## If something fails immediately

### `make bootstrap` fails

Run:

```bash
make deps
```

That will tell you which command is missing.

### `docker info` fails

Your container runtime is not running.

Fix:

- start Docker Desktop or OrbStack
- rerun `docker info`

### `nix develop` fails

Check:

```bash
nix flake metadata
```

If that fails, the problem is Nix or network access, not the lab itself.

## What you do not need to do

- do not install `kubectl` globally
- do not install `helm` globally
- do not install `terraform` globally
- do not install `kind` globally
- do not mix Homebrew versions of these tools into the same workflow unless you are debugging the toolchain itself

## Reference commands

Repo root:

```bash
pwd
```

Should be:

```text
/Users/deanlofts/gits/devops-roundrobin
```

Primary entry points:

- [`Makefile`](/Users/deanlofts/gits/devops-roundrobin/Makefile)
- [`README.md`](/Users/deanlofts/gits/devops-roundrobin/README.md)
- [`flake.nix`](/Users/deanlofts/gits/devops-roundrobin/flake.nix)
