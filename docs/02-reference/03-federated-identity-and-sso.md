# Federated Identity and SSO

This guide is the practical mental model you need for workforce access, application login, and common DevOps support work around identity.

## Why this matters

In most modern environments, people do not get long-lived local accounts for every system.

Instead:

- a central identity provider proves who the user is
- applications and cloud platforms trust that provider
- access is granted through short-lived sessions and roles

If you work in DevOps, you will almost certainly touch:

- SSO into AWS
- SSO into internal tools
- RBAC mapped from identity groups
- claims, roles, sessions, and temporary credentials
- provisioning and deprovisioning issues

## Core terms

### Identity provider

Usually called the IdP.

This is the system that authenticates the user.

Examples:

- Okta
- Microsoft Entra ID
- Google Workspace
- Active Directory through a bridge product

### Service provider or relying party

This is the application or platform that trusts the IdP and consumes the login result.

Examples:

- AWS IAM Identity Center
- Grafana
- an internal admin portal
- GitHub Enterprise

### Federation

Federation means one system accepts identity assertions from another trusted system instead of storing the user's primary password itself.

### SSO

Single Sign-On means one login session can be reused across multiple applications.

### Authentication versus authorization

- authentication answers who are you
- authorization answers what are you allowed to do

Do not blur them.

A user can authenticate successfully and still be denied access because authorization is wrong.

## The practical split you should know

There are three identity domains you will commonly deal with.

### Workforce identity

This is employee login.

Examples:

- logging into AWS
- logging into Grafana
- logging into GitHub

Best current AWS pattern:

- use AWS IAM Identity Center for workforce access
- federate it with your corporate IdP
- grant access through permission sets and IAM roles
- avoid creating long-lived IAM users for humans unless there is a very specific exception

### Customer or application identity

This is end-user login to an application.

Examples:

- users logging into a SaaS app
- login with Google or Microsoft

Common protocols:

- OpenID Connect for authentication
- OAuth 2.0 for delegated authorization

### Workload identity

This is machine-to-machine trust.

Examples:

- a pod accessing AWS resources
- GitHub Actions assuming an AWS role
- a CI system calling cloud APIs

This is not workforce SSO, even though the word identity still applies.

## The protocols you should understand

### SAML

SAML is older, XML-heavy, and still common for workforce SSO into enterprise tools.

You will still encounter it in:

- legacy enterprise apps
- some admin consoles
- older SaaS integrations

Good enough to know:

- browser is redirected to the IdP
- IdP authenticates user
- IdP sends a signed SAML assertion back
- the application trusts that assertion and starts a session

### OAuth 2.0

OAuth 2.0 is an authorization framework.

It is about delegated access.

Example:

- app A wants access to API B on behalf of a user

Important point:

OAuth alone is not “login.”

People often speak loosely here. Be precise.

### OpenID Connect

OIDC sits on top of OAuth 2.0 and adds identity.

This is the modern default for application authentication and increasingly common for workforce federation too.

Good enough to know:

- user is redirected to the IdP
- user authenticates there
- app receives an authorization code
- app exchanges the code for tokens
- app validates the ID token and starts a session

### SCIM

SCIM is about provisioning.

It is not login.

It is used to sync:

- users
- groups
- lifecycle changes

This matters because many “SSO is broken” incidents are actually provisioning or group mapping problems.

## The modern default mental model

If someone says “set up SSO,” the likely best-practice interpretation today is:

1. use a central IdP
2. use OIDC where supported
3. use SAML where required
4. use SCIM for provisioning if the platform supports it
5. map groups to roles
6. use short-lived sessions and roles, not shared credentials

## AWS-specific mental model

For workforce access to AWS, think in this sequence:

1. user signs into the corporate IdP
2. IdP federates into AWS IAM Identity Center
3. IAM Identity Center identifies the user and their groups
4. a permission set maps to an IAM role in a target AWS account
5. AWS issues short-lived credentials for that role
6. the user gets console or CLI access through those temporary credentials

That is the normal modern path.

## What IAM Identity Center actually does

Think of IAM Identity Center as the workforce access broker for AWS.

It helps with:

- federated user access
- account assignment
- permission sets
- centralized login

It does not replace IAM itself.

It sits above IAM roles and uses them.

## Permission sets versus IAM roles

This distinction matters in AWS.

- permission set: IAM Identity Center concept describing what access should be granted
- IAM role: the actual role created or used in the account

In practice:

- you assign a permission set to a user or group for an account
- AWS uses that to provide access through a role

## OIDC tokens you should know

### ID token

Used to convey identity information to the application.

Typical contents:

- issuer
- subject
- audience
- expiry
- user claims

### Access token

Used to access an API.

It is for authorization, not for proving browser session state to your frontend.

### Refresh token

Used to obtain new access or ID tokens without asking the user to log in again, depending on the application and provider.

## Common claims and mappings

When SSO “works but access is wrong,” these are usually involved:

- email
- subject
- groups
- roles
- audience
- issuer

Typical real problem:

- login succeeds
- application accepts the token
- user lands in the wrong role because the expected group claim is missing or renamed

## The flow you should be able to explain

### Workforce SSO to AWS

1. user opens AWS access portal
2. browser is redirected to the IdP
3. user authenticates with password, MFA, or passkey
4. IdP returns a trusted assertion to IAM Identity Center
5. IAM Identity Center resolves user, group membership, and account assignments
6. user chooses an AWS account and role or a default assignment is applied
7. AWS issues a short-lived session
8. user reaches the console or uses CLI credentials derived from that session

### App login with OIDC

1. user visits the app
2. app redirects to the IdP authorization endpoint
3. user authenticates
4. IdP returns an authorization code to the app callback
5. app exchanges code for tokens
6. app validates issuer, audience, expiry, signature, and nonce as appropriate
7. app starts its own session, often with a secure cookie
8. app uses identity and group claims for authorization

## Common DevOps tasks around SSO

These are realistic things you might handle.

### Onboard an internal app into SSO

Typical work:

- choose OIDC or SAML
- configure redirect URIs or ACS URLs
- configure client ID and secret or metadata
- map groups to roles
- test login, logout, and error paths

### Debug “login works, access denied”

Typical checks:

- user is in the right group
- group claim is being sent
- app expects the same claim name and format
- role mapping matches the intended group

### Debug “login loop” or callback failure

Typical checks:

- redirect URI mismatch
- wrong issuer URL
- wrong client ID or secret
- bad clock skew
- cookies blocked or session not persisted

### Debug expired or invalid sessions

Typical checks:

- token expiry
- refresh token usage
- cookie configuration
- TLS offload or proxy header issues

### Debug AWS access problems after successful SSO

Typical checks:

- IAM Identity Center assignment exists
- correct AWS account selected
- correct permission set applied
- role trust and resulting IAM permissions are correct
- policy or SCP is not blocking the action

## Common failure modes

### Authentication failure

User cannot complete login.

Likely areas:

- IdP config
- MFA
- callback URL
- cert or signing issue

### Authorization failure

User logs in but cannot do the thing they need.

Likely areas:

- missing group
- wrong group
- wrong role mapping
- permission set issue
- IAM policy issue

### Provisioning failure

User should have access but does not appear in the target platform correctly.

Likely areas:

- SCIM sync
- group push
- stale directory data

### Session failure

User logs in, then gets logged out or redirected back to login repeatedly.

Likely areas:

- cookies
- session storage
- reverse proxy config
- token expiry handling

## Questions you should ask during an incident

1. Did authentication fail, or did authorization fail after authentication?
2. Is the issue with one user, one group, one app, or all users?
3. Was there a recent IdP, app, certificate, redirect URI, or group mapping change?
4. Is provisioning involved, or only runtime login?
5. Is the target system trusting the correct issuer?

## A good interview-level answer

If asked how SSO typically works, a strong answer is:

“Users authenticate against a central IdP. The target application or platform trusts that IdP through OIDC or SAML. Authentication and authorization stay separate: the IdP proves identity, then claims or group mappings drive access. In AWS, the modern workforce pattern is IAM Identity Center federated with a corporate IdP, with permission sets mapping to IAM roles and short-lived sessions rather than long-lived IAM users.”

## What to study next after this

After you understand the flow, focus on these practical areas:

- OIDC authorization code flow
- SAML browser flow
- group and claim mapping
- IAM Identity Center permission sets
- STS and temporary credentials
- SCIM provisioning

## References

- [AWS IAM Identity Center: What is IAM Identity Center?](https://docs.aws.amazon.com/singlesignon/latest/userguide/what-is.html)
- [AWS IAM: Plan access to your AWS accounts](https://docs.aws.amazon.com/IAM/latest/UserGuide/gs-identities.html)
- [AWS Prescriptive Guidance: IAM Identity Center](https://docs.aws.amazon.com/prescriptive-guidance/latest/security-reference-architecture-identity-management/workforce-iam-identity-center.html)
- [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0-18.html)
