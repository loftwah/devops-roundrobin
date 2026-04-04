# AWS Platform Mapping

Use this guide to translate the local lab into current AWS-heavy platform work.

The point is not to pretend kind is AWS. The point is to make the local mechanics teach the same decisions, failure modes, and ownership boundaries you will see in a real estate.

## Identity and access

### Workforce access

Default AWS pattern today:

- workforce users sign in through AWS IAM Identity Center
- IAM Identity Center is federated with the corporate IdP
- permission sets map access into account-level IAM roles
- users get short-lived console or CLI credentials instead of long-lived IAM user keys

What to know cold:

- `aws configure sso` sets up a CLI profile for IAM Identity Center access
- `aws sts get-caller-identity` tells you which role and account your current credentials resolve to
- SCPs, permission sets, IAM policies, and resource policies can all block access even when login succeeds

Good local mapping in this repo:

- the SSO and federation mental model is covered in [`03-federated-identity-and-sso.md`](/Users/deanlofts/gits/devops-roundrobin/docs/02-reference/03-federated-identity-and-sso.md)

Missing from the repo before you join a new AWS-heavy team:

- one short hands-on drill that configures an AWS CLI SSO profile and checks identity with STS

### Workload identity

Current AWS pattern:

- ECS uses task roles for workloads
- EKS uses Pod Identity where possible, with IRSA still common in existing estates

What matters operationally:

- humans and workloads should not share credentials
- the trust policy matters as much as the permission policy
- short-lived credentials and clear role boundaries are the baseline

## Networking and request path

### VPC mental model

You should be able to explain:

- VPC
- public and private subnets
- route tables
- internet gateway
- NAT for private egress
- security groups
- network ACLs
- VPC endpoints

Local mapping here:

- Docker bridge networking and kind service networking teach dependency wiring
- ingress-nginx teaches HTTP routing and readiness dependency

Where the local lab is thinner than AWS:

- it does not currently force you to reason about subnet placement, egress design, or security-group boundaries
- it does not have a dedicated VPC troubleshooting drill

### DNS, ingress, ALB, and TLS

Good AWS translation:

- `app.127.0.0.1.nip.io` is standing in for a Route 53 record
- ingress-nginx is standing in for an ALB plus controller logic
- a Kubernetes Ingress rule is standing in for listener rules and target selection

Important differences to remember:

- in AWS, public DNS is usually Route 53
- public TLS is usually terminated with ACM certificates on an ALB or CloudFront
- target health is evaluated through target groups, not only Kubernetes readiness

If you join a company tomorrow, you should already be comfortable tracing:

1. Route 53 record
2. ALB listener and rule
3. target group health
4. service or ingress controller
5. pod or task readiness

## Compute choices

### ECS versus EKS

Good enough interview answer:

- choose ECS when you want a simpler AWS-native container control plane with less Kubernetes surface area to run
- choose EKS when you need Kubernetes APIs, portability, or an existing Kubernetes operating model

Operationally useful distinctions:

- ECS has less platform overhead for many teams
- EKS gives more flexibility and ecosystem depth, but also more cluster and add-on ownership
- both can run stateless app and worker services well

This repo prepares you better for:

- EKS-style thinking
- container packaging
- probes, rollouts, and ingress troubleshooting

This repo prepares you less directly for:

- ECS task definitions
- ECS services, capacity providers, and task roles

## Data services

### PostgreSQL and Redis

Accurate AWS mapping:

- local PostgreSQL maps to RDS or Aurora PostgreSQL
- local Redis maps to ElastiCache for Redis or Valkey

What the local lab teaches well:

- dependency health
- startup ordering
- secret injection
- failure impact on readiness

What it does not teach by itself:

- managed backups and snapshots
- parameter groups
- failover behaviour
- Multi-AZ or replication tradeoffs
- maintenance windows and service-linked events

## Observability

Local mapping:

- Prometheus metrics from the app and worker
- Grafana dashboards
- container and pod logs
- Kubernetes events

AWS translation:

- CloudWatch for logs, metrics, alarms, and events
- Amazon Managed Service for Prometheus for Prometheus-compatible metrics at scale
- Amazon Managed Grafana for dashboards

What is missing for stronger AWS readiness:

- one drill that starts from CloudWatch-style symptoms first, then pivots to workload logs and platform events
- one note on when a team keeps everything in CloudWatch versus when it adds managed Prometheus and Grafana

## Terraform patterns that transfer well

Good in this repo:

- environment root plus child modules
- separate namespace, platform, and monitoring concerns
- clear local provider configuration

Needs to be explained explicitly:

- Terraform state ownership is real ownership
- Terraform does not automatically take over resources created manually earlier in the walkthrough
- remote state, locking, encryption, and cross-account role assumption matter in AWS even though they stay out of scope locally

Common AWS estate patterns you should know before day one:

- S3 backend with encryption and locking
- CI or automation roles that assume into target accounts
- environment roots per account or stage
- shared modules for VPC, IAM, EKS or ECS, Route 53, RDS, and observability
- importing or deliberately recreating resources when changing the source of truth

## Incident handling

This repo teaches useful base instincts:

- separate health from readiness
- check metrics, logs, and events together
- prefer rollback when impact is active

To be closer to AWS-backed incidents, add these mental translations:

- Kubernetes ingress `503` maps to ALB listener or target-group failures
- DB auth failures map to Secrets Manager, RDS auth, or workload identity mistakes
- queue lag maps to SQS backlog, worker saturation, or downstream slowness
- “it deployed but is broken” often shows up first in CloudWatch alarms, target health, and rollout events

## What to add next

Highest-value additions for AWS-heavy job readiness:

1. a short AWS CLI SSO and STS drill
2. a VPC and request-path drill covering Route 53, ALB, ACM, subnets, and security groups
3. a workload identity note comparing ECS task roles, EKS Pod Identity, and IRSA
4. a Terraform ownership and import note for moving from manual changes to managed state
5. one ECS-specific comparison page so the repo does not feel EKS-only

## Primary references

- [AWS IAM Identity Center](https://docs.aws.amazon.com/singlesignon/latest/userguide/what-is.html)
- [AWS CLI IAM Identity Center configuration](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sso.html)
- [AWS STS get-caller-identity](https://docs.aws.amazon.com/cli/latest/reference/sts/get-caller-identity.html)
- [Amazon EKS workload access, including Pod Identity and IRSA](https://docs.aws.amazon.com/eks/latest/userguide/service-accounts.html)
- [Amazon Managed Service for Prometheus](https://docs.aws.amazon.com/prometheus/latest/userguide/what-is-Amazon-Managed-Service-Prometheus.html)
- [Amazon Managed Grafana](https://docs.aws.amazon.com/grafana/latest/userguide/what-is-Amazon-Managed-Service-Grafana.html)
- [Terraform import](https://developer.hashicorp.com/terraform/cli/commands/import)
