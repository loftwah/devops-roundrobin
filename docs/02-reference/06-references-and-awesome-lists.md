# References and Awesome Lists

Use this page when you want a controlled set of high-signal references.

Rule:

- official docs first for correctness
- awesome lists second for breadth and discovery

Do not let awesome lists turn into a distraction trap.

## Best official references for this repo

### Core lab stack

- [kind documentation](https://kind.sigs.k8s.io/)
- [Kubernetes probes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [ingress-nginx deploy guide for kind](https://kubernetes.github.io/ingress-nginx/deploy/#kind)
- [Helm docs](https://helm.sh/docs/)
- [Helm chart best practices](https://helm.sh/docs/chart_best_practices/)
- [Terraform docs](https://developer.hashicorp.com/terraform/docs)
- [Terraform module development guidance](https://developer.hashicorp.com/terraform/language/modules/develop)
- [Terraform Helm provider tutorial](https://developer.hashicorp.com/terraform/tutorials/kubernetes/helm-provider)
- [Prometheus Go instrumentation guide](https://prometheus.io/docs/guides/go-application/)
- [kube-prometheus-stack chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)

### AWS references

- [AWS IAM Identity Center](https://docs.aws.amazon.com/singlesignon/latest/userguide/what-is.html)
- [AWS IAM identity guide](https://docs.aws.amazon.com/IAM/latest/UserGuide/gs-identities.html)
- [AWS security reference architecture for workforce identity](https://docs.aws.amazon.com/prescriptive-guidance/latest/security-reference-architecture-identity-management/workforce-iam-identity-center.html)
- [Amazon EKS docs](https://docs.aws.amazon.com/eks/)
- [Amazon ECS docs](https://docs.aws.amazon.com/ecs/)
- [Amazon RDS docs](https://docs.aws.amazon.com/rds/)
- [Amazon ElastiCache docs](https://docs.aws.amazon.com/elasticache/)
- [Amazon Route 53 docs](https://docs.aws.amazon.com/route53/)
- [AWS Load Balancer docs](https://docs.aws.amazon.com/elasticloadbalancing/)
- [AWS CloudWatch docs](https://docs.aws.amazon.com/cloudwatch/)

### Protocols and web path references

- [OpenID Connect Core](https://openid.net/specs/openid-connect-core-1_0-18.html)
- [RFC 1034: DNS concepts](https://www.rfc-editor.org/rfc/rfc1034)
- [RFC 1035: DNS specification](https://www.rfc-editor.org/rfc/rfc1035)
- [RFC 9293: TCP](https://www.rfc-editor.org/rfc/rfc9293)
- [RFC 8446: TLS 1.3](https://www.rfc-editor.org/rfc/rfc8446)
- [RFC 9110: HTTP semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [RFC 9114: HTTP/3](https://www.rfc-editor.org/rfc/rfc9114)

## Awesome lists worth using

These are not source-of-truth documents. They are discovery tools.

### [awesome-sre](https://github.com/dastergon/awesome-sre)

Why it is worth using:

- broad and well-known SRE resource collection
- useful for reliability, observability, on-call, postmortems, and production engineering

Use it when:

- you want deeper reading after finishing the lab
- you want books, talks, and operational practices

### [awesome-kubernetes](https://github.com/ramitsurana/awesome-kubernetes)

Why it is worth using:

- broad entry point for Kubernetes tools and ecosystem references

Use it when:

- you want to discover adjacent tooling
- you want to see categories of Kubernetes resources, not just one product

Warning:

- it is broad and noisy
- use it for discovery, not for deciding best practice by itself

### [awesome-tf](https://github.com/shuaibiyy/awesome-tf)

Why it is worth using:

- Terraform and OpenTofu ecosystem references in one place

Use it when:

- you want to compare tooling around Terraform workflows
- you want broader ecosystem awareness after you understand the basics

### [awesome-sysadmin](https://github.com/awesome-foss/awesome-sysadmin)

Why it is worth using:

- good breadth for operations, CI/CD, monitoring, PKI, and infrastructure-adjacent tooling

Use it when:

- you need ideas for the wider systems and platform space
- you want exposure beyond this repo’s stack

### [awesome-networking](https://github.com/nyquist/awesome-networking)

Why it is worth using:

- useful for DNS, TCP, routing, monitoring, automation, and network operations discovery

Use it when:

- you are strengthening your network fundamentals
- you want references around the “URL to page render” path

### [awesome-platform-engineering](https://github.com/toptechevangelist/awesome-platform-engineering)

Why it is worth using:

- useful for seeing how modern platform engineering is framed

Use it when:

- you want to connect this lab to broader platform engineering concepts

Warning:

- more useful for landscape awareness than for day-one operational correctness

## How to use these without wasting time

Use this order:

1. official docs for the exact tool or protocol
2. repo docs for the local implementation
3. awesome lists only when you want breadth or alternatives

If you are one week from starting a job, your goal is not to consume every link.

Your goal is:

- know the core flows cold
- know where to look when stuck
- know a few strong external references

## Best references for immediate interview and job readiness

If you only use a few, use these:

1. [Start-to-Finish Walkthrough](../01-walkthrough/01-start-to-finish.md)
2. [Cheatsheets](05-cheatsheets.md)
3. [When You Get Stuck](01-when-stuck.md)
4. [Federated Identity and SSO](03-federated-identity-and-sso.md)
5. [From URL to Rendered Page](04-from-url-to-rendered-page.md)
6. [awesome-sre](https://github.com/dastergon/awesome-sre)
7. [awesome-networking](https://github.com/nyquist/awesome-networking)
