# From URL to Rendered Page

This guide gives you the full request path from typing a URL to seeing a page on screen.

A DevOps engineer does not usually own every layer, but you need to understand where failures happen and how the layers connect.

## The short version

When a user opens a URL, the path is roughly:

1. browser parses the URL
2. browser checks cache and existing connection state
3. DNS resolves the hostname to an IP address
4. transport is established with TCP or QUIC
5. TLS negotiates encryption and validates the certificate
6. HTTP request is sent
7. edge systems route and protect the request
8. app responds
9. browser parses HTML and fetches more assets
10. browser builds layout, paints, and composites the page

That is the end-to-end path you should be able to explain.

## Step 1: The browser parses the URL

Example:

```text
https://app.example.com/dashboard?tab=active
```

The browser separates:

- scheme: `https`
- host: `app.example.com`
- path: `/dashboard`
- query string: `tab=active`
- port: implied `443` for HTTPS unless explicitly set

Why it matters:

- the scheme influences transport and security
- the host drives DNS and TLS validation
- path and query influence routing and application behaviour

## Step 2: The browser checks what it already knows

Before it goes to the network, the browser may reuse:

- cached DNS
- cached HTTP responses
- existing TCP or QUIC connections
- cookies
- HSTS rules
- service worker logic

Why it matters:

Not every page load starts from zero. This is why “works in one browser tab but not another” can be real.

## Step 3: DNS resolution

The browser needs an IP address for the hostname.

The typical path is:

1. browser cache
2. OS resolver cache
3. configured recursive resolver
4. recursive resolver queries authoritative DNS servers
5. answer is returned and cached based on TTL

The important mental model:

- authoritative servers hold the source of truth for the zone
- recursive resolvers do the lookup work on behalf of clients
- TTL controls how long an answer may be reused

What a DevOps engineer touches here:

- Route53 or another DNS provider
- A, AAAA, CNAME, ALIAS, TXT records
- TTL choices
- cutovers and migrations

Common failures:

- wrong record value
- stale caches
- missing record
- split-horizon DNS confusion
- certificate valid for one name but DNS points traffic somewhere else

AWS equivalent:

- Route53 hosted zones and records

## Step 4: Transport setup

After DNS, the client connects to the destination.

### HTTP/1.1 and HTTP/2 usually use TCP

TCP gives:

- reliable ordered delivery
- retransmission
- flow control
- a three-way handshake

Very high-level TCP sequence:

1. client sends SYN
2. server replies SYN-ACK
3. client replies ACK

### HTTP/3 uses QUIC over UDP

Modern stacks may use HTTP/3 instead, which runs over QUIC and avoids some TCP-level head-of-line blocking issues.

You do not need to be a QUIC expert for most DevOps roles, but you should know it exists.

What a DevOps engineer touches here:

- security groups
- NACLs
- firewalls
- load balancer listeners
- idle timeouts

Common failures:

- port blocked
- target unreachable
- listener mismatch
- timeout or packet loss

AWS equivalent:

- security groups
- NACLs
- ALB or NLB listeners

## Step 5: TLS handshake

For HTTPS, the client and server negotiate encryption.

Key concepts:

- certificate presented by server
- hostname validation
- trusted CA chain
- SNI selects the correct certificate on multi-tenant endpoints
- ALPN negotiates HTTP/1.1, HTTP/2, or HTTP/3

What the client checks:

- certificate is valid now
- certificate chain is trusted
- certificate name matches the host

What a DevOps engineer touches here:

- ACM or another certificate manager
- ingress TLS config
- reverse proxy config
- certificate rotation

Common failures:

- expired certificate
- wrong certificate
- missing intermediate certificate
- TLS terminated at the wrong layer
- redirect loops after TLS offload

AWS equivalent:

- ACM certificates
- ALB HTTPS listeners
- CloudFront certificates

## Step 6: HTTP request

Once transport and TLS are ready, the browser sends the HTTP request.

Typical request data:

- method
- path
- headers
- cookies
- maybe a request body

Examples:

- `GET /dashboard`
- `Host: app.example.com`
- `Cookie: session=...`

Important browser behaviour:

- cookies may carry app session state
- headers may influence caching, content negotiation, auth, and routing

What a DevOps engineer touches here:

- reverse proxies
- ingress rules
- header forwarding
- auth gateways
- WAF rules
- cache control

Common failures:

- bad header forwarding
- wrong host routing
- auth cookie not forwarded
- body size limit exceeded
- CORS or proxy misconfiguration

## Step 7: Edge routing and platform path

In a modern setup, the request often passes through several components.

Typical path:

1. DNS points to CDN or load balancer
2. CDN may serve cached assets or forward request
3. WAF may inspect or block
4. load balancer selects a target group
5. ingress or reverse proxy routes to the service
6. service forwards to the app pod or container

In this repo’s local Kubernetes path, the equivalent is:

1. browser resolves `app.127.0.0.1.nip.io`
2. request reaches host port 80
3. ingress-nginx receives the request
4. ingress rule maps host to the app Service
5. Service routes to a ready app pod

In AWS, a common equivalent is:

1. Route53
2. CloudFront optionally
3. WAF optionally
4. ALB
5. ingress controller or direct target group
6. ECS tasks or EKS pods

## Step 8: The application handles the request

At the app layer, several things can happen:

- session is validated
- authorization is checked
- database or cache is queried
- business logic runs
- HTML, JSON, or redirect is returned

In this lab:

- the app exposes `/`, `/health`, `/ready`, `/metrics`, and `/jobs`
- readiness checks Postgres and Redis connectivity

What a DevOps engineer touches here:

- env vars
- Secrets and ConfigMaps
- database connectivity
- rollout behaviour
- logs and metrics

Common failures:

- bad config
- DB connection failure
- Redis connection failure
- pod not ready
- app crash

## Step 9: The browser parses the response

If the response is HTML, the browser:

1. parses HTML
2. discovers CSS, JS, fonts, images, and API calls
3. issues more requests for those assets

This means “the page loaded” can still hide problems:

- HTML may arrive
- CSS may fail
- JS bundle may fail
- API calls may fail

You need to think in waterfalls, not in a single request.

## Step 10: Rendering

Very roughly:

1. browser builds the DOM from HTML
2. browser parses CSS
3. browser builds style information
4. browser creates layout
5. browser paints pixels
6. browser composites layers to the screen

If JS frameworks are involved:

- JS may render client-side content
- hydration may attach behaviour after HTML arrives

This matters because the network request can be fine while the page is still broken.

## The full practical path in one line

User enters URL -> browser parses -> caches checked -> DNS lookup -> TCP or QUIC connection -> TLS handshake -> HTTP request -> CDN or LB -> ingress or proxy -> service -> app -> app dependencies -> response -> browser fetches assets -> render and hydrate.

## Where DevOps usually gets involved

This is the part you should be able to reason about in production.

### DNS

- wrong record
- bad TTL
- cutover issues

### Network

- blocked ports
- bad listener
- unreachable target

### TLS

- expired cert
- wrong SAN
- missing chain

### Edge routing

- ingress rule wrong
- ALB target unhealthy
- service has no endpoints

### App platform

- pod not ready
- wrong env vars
- bad deploy

### App dependencies

- DB unreachable
- cache unreachable
- secrets wrong

### Browser and frontend

- broken JS bundle
- wrong base URL
- API returning 401 or 500

## A good incident debugging sequence

When “the site is down,” go in this order:

1. can the hostname resolve?
2. can the endpoint accept a TCP or TLS connection?
3. does the edge return anything?
4. is the load balancer or ingress healthy?
5. do ready app instances exist?
6. do app logs show dependency failures?
7. are downstream dependencies healthy?
8. does the browser fail on secondary assets or API calls?

That order prevents random guessing.

## How this maps to this repo

Compose path:

- browser to `localhost:8080`
- Docker port mapping
- app process
- Postgres and Redis dependencies

kind path:

- browser to `app.127.0.0.1.nip.io`
- host port 80
- ingress-nginx
- Kubernetes Service
- app pod
- Postgres and Redis in cluster

## A good interview-level answer

If asked how a page renders from a URL, a strong answer is:

“The browser parses the URL, checks caches, resolves DNS, establishes transport with TCP or QUIC, negotiates TLS, sends the HTTP request, then the request passes through edge components like CDN, WAF, load balancer, and ingress before reaching the application. The application may call dependencies like a database or cache, returns HTML or JSON, and then the browser fetches any additional assets, builds the DOM and CSSOM, lays out the page, paints it, and composites it on screen. Operationally, outages can happen at any layer, so I debug in that order instead of jumping straight to app logs.”

## References

- [RFC 1034: Domain names - concepts and facilities](https://www.rfc-editor.org/rfc/rfc1034)
- [RFC 1035: Domain names - implementation and specification](https://www.rfc-editor.org/rfc/rfc1035)
- [RFC 9293: Transmission Control Protocol](https://www.rfc-editor.org/rfc/rfc9293)
- [RFC 8446: TLS 1.3](https://www.rfc-editor.org/rfc/rfc8446)
- [RFC 9110: HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [RFC 9114: HTTP/3](https://www.rfc-editor.org/rfc/rfc9114)
