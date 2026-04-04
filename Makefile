SHELL := /bin/bash
PROJECT_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
ENV_FILE ?= $(PROJECT_ROOT)/.env

ifneq ("$(wildcard $(ENV_FILE))","")
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))
endif

KIND_CLUSTER_NAME ?= round-robin
KUBE_NAMESPACE ?= round-robin
RAW_KUBE_NAMESPACE ?= round-robin-raw
APP_IMAGE ?= round-robin-app:local
WORKER_IMAGE ?= round-robin-worker:local
TF_DIR := infra/terraform/environments/local
HELM_CHART := infra/helm/round-robin

.PHONY: help bootstrap env-check deps app-run worker-run test lint \
	build-app build-worker docker-build compose-up compose-down compose-logs \
	compose-ps kind-up kind-down kind-load ingress-install k8s-apply-raw \
	k8s-delete-raw helm-lint helm-install helm-upgrade helm-rollback \
	helm-uninstall tf-prepare tf-init tf-plan tf-apply tf-destroy \
	monitor-install scan-image load-test

help:
	@printf "\nAvailable targets:\n"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9._-]+:.*##/ {printf "  %-22s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

bootstrap: env-check deps ## Prepare the repo for local work

env-check: ## Verify that .env exists
	@test -f $(ENV_FILE) || (echo "Missing $(ENV_FILE). Copy .env.example to .env first." && exit 1)

deps: ## Check core local dependencies
	@./scripts/dev/check-prereqs.sh

app-run: env-check ## Run the app directly on the host
	go run ./app

worker-run: env-check ## Run the worker directly on the host
	go run ./worker

test: ## Run unit tests
	go test ./...

lint: ## Run basic validation for local manifests and IaC
	helm lint $(HELM_CHART)
	kubectl kustomize infra/k8s/raw >/dev/null
	terraform -chdir=$(TF_DIR) fmt -check -recursive
	terraform -chdir=$(TF_DIR) init -backend=false >/dev/null
	terraform -chdir=$(TF_DIR) validate

build-app: ## Build the app binary
	CGO_ENABLED=0 go build -o bin/app ./app

build-worker: ## Build the worker binary
	CGO_ENABLED=0 go build -o bin/worker ./worker

docker-build: ## Build app and worker images
	docker build -t $(APP_IMAGE) -f app/Dockerfile .
	docker build -t $(WORKER_IMAGE) -f worker/Dockerfile .

compose-up: env-check ## Start the full local Docker Compose stack
	docker compose up --build -d

compose-down: ## Stop the local Docker Compose stack
	docker compose down --volumes --remove-orphans

compose-logs: ## Tail Docker Compose logs
	docker compose logs -f

compose-ps: ## Show Docker Compose service status
	docker compose ps

kind-up: env-check ## Create the kind cluster and local registry wiring
	./scripts/kind/create-cluster.sh

kind-down: env-check ## Delete the kind cluster
	./scripts/kind/delete-cluster.sh

kind-load: env-check docker-build ## Load local images into the kind cluster
	./scripts/kind/load-images.sh

ingress-install: env-check ## Install ingress-nginx into kind
	./scripts/kind/install-ingress.sh

k8s-apply-raw: env-check kind-load ## Apply the raw Kubernetes manifests
	kubectl apply -k infra/k8s/raw

k8s-delete-raw: env-check ## Remove the raw Kubernetes manifests
	kubectl delete -k infra/k8s/raw --ignore-not-found

helm-lint: ## Lint the round-robin Helm chart
	helm lint $(HELM_CHART)

helm-install: env-check kind-load ## Install the Helm release directly
	helm upgrade --install round-robin $(HELM_CHART) \
		--namespace $(KUBE_NAMESPACE) \
		--create-namespace \
		-f $(HELM_CHART)/values.yaml

helm-upgrade: helm-install ## Upgrade the Helm release

helm-rollback: ## Roll back the last Helm release revision
	helm rollback round-robin --namespace $(KUBE_NAMESPACE)

helm-uninstall: ## Remove the manual Helm release from the local cluster
	helm uninstall round-robin --namespace $(KUBE_NAMESPACE) --ignore-not-found
	@if kubectl get namespace $(KUBE_NAMESPACE) >/dev/null 2>&1; then \
		kubectl delete namespace $(KUBE_NAMESPACE) --ignore-not-found; \
		kubectl wait --for=delete namespace/$(KUBE_NAMESPACE) --timeout=180s; \
	fi

tf-prepare: env-check ## Clear the manual app release before Terraform takes ownership
	@$(MAKE) helm-uninstall

tf-init: env-check ## Initialise Terraform for the local environment
	terraform -chdir=$(TF_DIR) init

tf-plan: env-check ## Create a Terraform plan for the local environment
	terraform -chdir=$(TF_DIR) plan -out=tfplan

tf-apply: env-check ## Apply Terraform for the local environment
	terraform -chdir=$(TF_DIR) apply

tf-destroy: env-check ## Destroy Terraform-managed local resources
	terraform -chdir=$(TF_DIR) destroy

monitor-install: env-check ## Install kube-prometheus-stack via Terraform
	terraform -chdir=$(TF_DIR) apply -var='enable_monitoring=true'

scan-image: ## Scan locally built images with Trivy
	trivy image $(APP_IMAGE)
	trivy image $(WORKER_IMAGE)

load-test: ## Run a short k6 load test against the local app
	k6 run scripts/load-test.js
