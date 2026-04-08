# 1. Identity Service

## Purpose
Establish who the user is before any token can be issued.

## Core Requirements
- Must accept Basic Auth (username/password or PAT).
- Must support whatever credential types your registry ecosystem uses.

## Credential Validation
- Must validate credentials against an identity backend

## Identity Resolution
Must return a canonical user identity:


## Failure Modes
- Reject missing credentials
- Reject invalid credentials
- Reject expired or revoked PATs
- Reject disabled users

# Output:
A verified identity:

```bash
{
  "user_id": "12345",
  "username": "samalba",
  "groups": ["dev", "ops"],
  "is_service_account": false
}
```



# 2. Authorization Service

## Purpose
Determine what the authenticated user is allowed to do.

## Core Requirements
Must parse the requested scope string, e.g.:

repository:samalba/my-app:pull,push

## Permission Evaluation
- Must check whether the authenticated user has rights to:
pull
push
delete
list
any other registry‑specific action

## This requires an authorization backend, such as:
- RBAC database
- ACLs stored in the registry
- Organization/team membership mapping
- Policy engine (OPA, Cedar, Rego, etc.)

## Decision Logic
- If user lacks permission → deny token issuance
- If user has permission → return allowed actions

## Output

A list of allowed actions to embed in the token:
{
  "repository": "samalba/my-app",
  "actions": ["pull"]
}


# 3. Token Service (Issuer)

## Purpose
Issue a short‑lived bearer token encoding the allowed actions.

## Requirements
Generate opaque token
Embed access section (allowed actions)
Set expiry metadata
Sign token (if JWT) or store it (if opaque)

# 4. Registry Enforcement Service

## Purpose
Validate the token and enforce access control on actual operations.

## Requirements
- Validate token signature or lookup
- Check expiry
- Compare requested operation with token’s allowed actions
- Allow or deny the operation


####

1. Identity Service (AuthN)
Validates credentials

Resolves user identity

Rejects invalid users

2. Authorization Service (AuthZ)
Evaluates whether the authenticated user may perform the requested actions

Produces the allowed action set for token embedding


# 1. Anonymous is just “no subject”
In the OCI auth model, the /token endpoint is an authorization token service, not an authentication service.

GET /token?scope=repository:foo/bar:pull

# 2. The registry still needs to return a token
OCI registries must return a token for any allowed request, even anonymous ones.

So even anonymous users get a token — but a token that encodes:
no subject
only public permissions
only allowed actions

# 3. What goes into an anonymous JWT?
sub is omitted or set to "anonymous"
access contains only public permissions
exp is short
iss, aud are normal

```json
{
  "iss": "https://registry.example.com",
  "aud": "registry.example.com",
  "sub": "anonymous",
  "access": [
    {
      "type": "repository",
      "name": "library/nginx",
      "actions": ["pull"]
    }
  ],
  "exp": 1711220400
}
```

# Thinking

Authenticated user → token with scoped permissions
Anonymous user → token with scoped permissions

A subject with no identity and minimal permissions.



# Permission(Subject, Resource, Action, Constraints)

# Subject Hiearchy

Subject {
    id
    type: user | group | org | service-account
    parent: Subject?   // for hierarchy
    attributes: {...}
}

Permissions can be inherited:
	org grants → inherited by teams → inherited by users
	group grants → inherited by members

org:acme
  └── team:backend
        └── user:marten

# Resource Hiearchy

registry
  └── namespace (org)
        └── repository
              └── tag
              └── digest



Resource {
    id
    type: registry | namespace | repository | artifact
    parent: Resource?
    attributes: {...}
}


# Example

Granting
namespace:acme/* → pull
gives
acme/app1
acme/app2
acme/app3


# Action Set

OCI defines a small but important action vocabulary:

pull
push
delete
list
mount
catalog

Action {
    id
    implies: [other actions]?  // optional
}

# Permission {
    subject: Subject
    resource: Resource
    actions: [Action]
    constraints: {...}   // optional: time, IP, token type, etc.
    effect: allow | deny
}

Permissio(
    SubjectHierarchy,
    ResourceHierarchy,
    ActionSet,
    Constraints
)
n

# Challenge

The logical access‑control model is hierarchical and expressive…
but humans hate configuring raw S–R–A (Subject–Resource–Action) rules.


Permission(SubjectHierarchy, ResourceHierarchy, ActionSet)

So vendors introduce opinionated abstractions that collapse the hierarchy into manageable units.


JFrog collapses the complexity into four user‑facing concepts:
Users
Groups
Repositories (local, remote, virtual)
Permission Targets

# 1. Users → Subject leaf nodes
JFrog hides the subject hierarchy by:
Treating users as atomic subjects
Letting groups represent hierarchy instead of nested org/team structures

#2. Groups → Subject hierarchy (flattened)

#

Subject hierarchy (org → team → user)
- Users + flat Groups

Resource hierarchy (namespace → repo → artifact)
- Repositories only

Action vocabulary
- Predefined roles (read/write/delete)

Permission rules
- Permission Targets (bundled S–R–A)

Complex inheritance
- No inheritance, no overrides

Fine‑grained ACLs
- Coarse, repo‑level permissions

# Subjects
Subject ::= User | Group
User ::= { id: UserId }
Group ::= { id: GroupId }

# Resources
Resource ::= Repository
Repository ::= { id: RepoId, type: Local | Remote | Virtual }


# Actions
Action ::= READ | ANNOTATE | DEPLOY | DELETE | MANAGE

DEPLOY ⇒ READ
DELETE ⇒ DEPLOY
MANAGE ⇒ READ ∧ ANNOTATE ∧ DEPLOY ∧ DELETE

# Permission Targets

PermissionTarget ::= {
    name: String,
    resources: Set<Resource>,
    subjects: Set<Subject>,
    actions: Set<Action>
}

# Effective Permissions

EffectiveActions(subject, resource) =
    ⋃ { PT.actions |
        PT ∈ PermissionTargets ∧
        subject ∈ PT.subjects ∧
        resource ∈ PT.resources }

subject ∈ PT.subjects
    if subject is a User and subject ∈ PT.subjects
    OR subject is a User and ∃ g ∈ Groups(subject) such that g ∈ PT.subjects


# Evaluation Rule

Request ::= { subject: Subject, resource: Resource, action: Action }

Authorization Decision:

Allow(Request) ⇔
    action ∈ EffectiveActions(subject, resource)

# Token Integration

TokenCapability(token) =
    { (resource, actions) |
        actions = EffectiveActions(subject, resource) }



“Can subject S perform action A on resource R?”

# Subject

Logical fields:

```text
subject_id (user, service account, robot account)
type (human, service account, group)
groups (optional)
attributes (optional: org, team, metadata)
```

The registry must know who is requesting access.

# Resource
Represents what is being accessed.

repository:<namespace>/<name>

```text
resource_type = repository
resource_id = namespace/name
artifact (optional: tag, digest)
attributes (optional: visibility, owner)
```

# Action

OCI defines a small, finite set:

pull
push
delete
list
mount (optional)
catalog (optional)

# Permission / Policy Rule

This is the relationship between Subject, Resource, and Action.
{
  "subject": "user-or-group",
  "resource": "repository:namespace/name",
  "actions": ["pull", "push"]
}


# Model

Permission(subject, resource, action)


Subject S is allowed to perform Action A on Resource R.

# Extended Logical Model (Optional but Common)

```
group_id → list of subjects
role_id → list of actions
subject/group → role → resource
```
# 1. Capability‑Based Security (1970s → now)

A capability is an unforgeable token that grants the holder the right to perform specific actions.

Reference tokens behave exactly like capabilities:

They are opaque
They point to server‑side rights
Possession = authority
Revocation is instant (delete the capability)


# 2. Token‑by‑Reference (OAuth2 Introspection Model, 2015 → now)

Token‑by‑value:
- JWT, Docker registry tokens, OAuth2 access tokens
- All permissions encoded inside the token
- Hard to revoke

Token‑by‑reference:
- Token is just an ID
- Server stores the permissions
- Server introspects the token on every use
- Revocation is instant

# This model was formalized in:
OAuth2 Token Introspection (RFC 7662)
OAuth2 Reference Tokens (used by Okta, Auth0, Ping, etc.)

"OAuth2 reference tokens without the external introspection endpoint."

# A capability‑based, token‑by‑reference authorization model with OAuth2‑style introspection semantics.

# Use Cases

## Who is the user ?
A CI pipeline uses a PAT to authenticate.
Identity service returns:

```json
user_id = ci-bot-123
groups = ["build", "deploy"]
```

# What is the user allowed to do?
User requests:
scope = repository:my-app:push

Authorization service checks:
- Does user have push rights on my-app?
- If yes → allow
- If no → deny token issuance

# Store the capability that the token refers to.

Stores token metadata:
token ID
subject
allowed actions
expiry
revocation status
Supports instant revocation
Supports permission updates

# -

Reference tokens are opaque.
The token itself contains no permissions.
The registry must look them up here.

# Token ID
ref:8f3a2c9d...

Stored capability:

actions = ["pull"]
expires_at = 2026-03-23T20:00:00Z

# Create a short‑lived, opaque token that points to stored permissions
Generates a random opaque token ID
Stores the capability in the token store
Returns the token to the client
Sets TTL and metadata

This replaces JWT generation.
The token is now a reference, not a self‑contained object.

# Client Receives
Authorization: Bearer ref:8f3a2c9d...


# Introspection Service: What does this token allow?

Registry sends the opaque token to the introspection service
Service looks up the token in the reference token store
Validates
existence
expiry
revocation
allowed actions
Returns a decision:

# Example

Registry asks:

introspect(ref:8f3a2c9d...)

Response:
actions = ["pull"]
valid = true

# Registry Enforcement Layer

Allow or deny the actual OCI operation.

Client tries to push.
Token only allows pull.
Registry denies.


Here’s the crisp, operator‑centric answer you’re looking for.

---

# 🎯 Requirements for the **Reference Token String Format Scheme**

A reference token is *just an identifier*, but the identifier itself must satisfy several strict requirements so that the entire capability‑based, token‑by‑reference model remains:

- unforgeable
- unguessable
- revocable
- efficient to store and lookup
- safe to expose to untrusted clients

Below is the complete set of requirements for designing the **reference token string format**.

---

# ✅ **1. Unforgeability**
The token string must be impossible for a client to guess or construct.

**Requirements:**
- Must be generated from a cryptographically secure RNG.
- Minimum entropy: **128 bits**, ideally **192–256 bits**.
- No deterministic or user‑derived components.

**Why:**
Possession = authority. If a token can be guessed, the system is broken.

---

# ✅ **2. Opaqueness**
The token must reveal **zero information** about:

- the subject
- permissions
- expiry
- issuer
- environment
- internal structure

**Requirements:**
- No embedded metadata.
- No encoding of claims.
- No structured format like JWT, PASETO, or base64‑encoded JSON.

**Why:**
Opaque tokens behave like capabilities: the server is the only source of truth.

---

# ✅ **3. Prefix‑Based Routing (Optional but Recommended)**
A short prefix can help route tokens to the correct backend store.

Example:
```
ref:8f3a2c9d...
```

**Requirements:**
- Prefix must not reduce entropy.
- Prefix must not encode sensitive metadata.
- Prefix must be static and non‑semantic.

**Why:**
Operational convenience: multiple token stores, multiple issuers, versioning.

---

# ✅ **4. URL‑Safe and Header‑Safe**
The token must be safe to transmit in:

- HTTP headers
- URLs
- logs (if redacted)
- CLI tools

**Requirements:**
- Use a URL‑safe alphabet (Base64URL or hex).
- Avoid characters requiring escaping (`+`, `/`, `=`).
- Avoid whitespace.

**Why:**
Reduces friction for clients and avoids encoding bugs.

---

# ✅ **5. Fixed or Bounded Length**
The token should have a predictable length for:

- storage
- indexing
- rate‑limiting
- log scrubbing
- regex detection

**Requirements:**
- Fixed length strongly preferred (e.g., 43 chars for 256‑bit Base64URL).
- If variable, must be bounded and non‑semantic.

**Why:**
Operational ergonomics and consistent parsing.

---

# ✅ **6. Stateless Validation of Format**
The token format must allow the server to:

- quickly reject malformed tokens
- avoid unnecessary DB lookups

**Requirements:**
- Validate prefix (if used).
- Validate length.
- Validate allowed characters.

**Why:**
Protects the introspection service from garbage input.

---

# ✅ **7. No Embedded Signature**
Unlike JWTs, reference tokens must **not** contain:

- signatures
- HMACs
- public keys
- timestamps

**Requirements:**
- Token is *not* self‑verifying.
- All validation happens server‑side.

**Why:**
This is the essence of token‑by‑reference.

---

# 🎯 **Summary Table**

| Requirement | Why it matters |
|------------|----------------|
| High entropy (≥128 bits) | Prevent guessing/forgery |
| Opaque (no metadata) | Capability semantics |
| URL/header safe | Works everywhere |
| Fixed length | Operational simplicity |
| Prefix optional | Routing/versioning |
| Stateless format validation | Protects introspection service |
| No embedded claims/signatures | Ensures token‑by‑reference |

---

# 🧩 **Example Valid Format**
```
ref:8f3a2c9d7b1e4f0a9c3d2e1f6b8a7c4d
```

- `ref:` = static prefix
- `8f3a2c9d…` = 256‑bit random hex

---

# 🧩 **Example Minimal Format (no prefix)**
```
k9F2xYp7QmC8r1Zt4Vw6n0Hs3Bd5Lj8PqR0Tu2WxYzA
```

- 256‑bit Base64URL
- No metadata
- Opaque

---

If you want, I can also design:

- a complete reference token spec
- a storage schema
- introspection API contract
- revocation model
- lifecycle and TTL strategy

Just tell me what direction you want to take next.

# Registry Token

## 1. iss - Issuer
"iss": "https://registry.example.com"
Used by the registry to verify the token’s origin.

## 2. sub — Subject
Identifies the authenticated principal.
"sub": "user:marten"
This is not used for access control directly — it’s for auditing and traceability.

## 3. aud — Audience
Must match the registry hostname.
"aud": "registry.example.com"
Prevents token replay against other registries.

## 4. exp — Expiration
Short‑lived, typically 5–30 minutes.
"exp": 1711220400

## 5. nbf / iat — Not‑before / Issued‑at

## 6. access — The OCI‑specific authorization claim
"access": [
  {
    "type": "repository",
    "name": "acme/app",
    "actions": ["pull", "push"]
  }
]

# Optional
jti — Token ID
scope - "scope": "repository:acme/app:pull,push"
realm - "realm": "https://registry.example.com/token"

# OCI registries return both token and access_token because:
access_token is the OAuth2‑style field name
token is the legacy Docker field name
Different clients expect different names
The OCI spec allows both
Vendors include both for maximum compatibility
They are the same JWT, just exposed under two names.


# Token format

## JSON Web Token
A internet standard for creating data with optional signature and/or optional encryption whose payload
holds JSON that asserts some number of claims.
Tokens are signed using a private secret or a public/private key.

A claim such as "logged in as administrator". The token contains this claim.
The token can be signed. Any party could verify if the token is legitimate, using public key.
Tokens are compact, URL safe, in SSO context.

{
  "alg": "RS256",
  "typ": "JWT",
  "x5c": [
    "MIIEFjCCAv6gAwIBAgIUOrLWyQZq2anewzZxX7RKlt7l5TMwDQYJKoZIhvcNAQELBQAwgYYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRIwEAYDVQQHEwlQYWxvIEFsdG8xFTATBgNVBAoTDERvY2tlciwgSW5jLjEUMBIGA1UECxMLRW5naW5lZXJpbmcxITAfBgNVBAMTGERvY2tlciwgSW5jLiBFbmcgUm9vdCBDQTAeFw0yNTA5MjUwMDMwMDBaFw0yNjA5MjUwMDMwMDBaMIGFMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTESMBAGA1UEBxMJUGFsbyBBbHRvMRUwEwYDVQQKEwxEb2NrZXIsIEluYy4xFDASBgNVBAsTC0VuZ2luZWVyaW5nMSAwHgYDVQQDExdEb2NrZXIsIEluYy4gRW5nIEpXVCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANtsyvxWl6kN2bPFwZ6uUhKWsA+Y1nE0VopwYG1mh1cxgXKJY1gloFWxyvOrBY71CWG8pN6HeHktaDxqQ7BzaJOmqbb6G8qAgdR07uvF9bMZDWdelJ4brbbXBLZ34NF0rfsAimPnAdI0Vp/QKbgaKuigELRq02FhnjnlTi+HMfMUaNhfb4cRHXlwgyjirnp8L7syg4qTMcCYfhZvYbvK62orykg9Ph7hBOOnk2Em7jAh/BhUc1MAHYtpxVHzssFsBoNiHaFH487p+vR9EQVQ69x/yEiKrRWskD0Ze8QL3y1ZSJns4t/2hDTEM1eIjriBoRc9dswRMJo0xKz0tnH/nQ8CAwEAAaN7MHkwDgYDVR0PAQH/BAQDAgGmMBMGA1UdJQQMMAoGCCsGAQUFBwMBMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFGqK3OPv9W13jUvyM/62x7F6eXbZMB8GA1UdIwQYMBaAFC60UPNeBkogY22tXPcBMHFvG3CsMA0GCSqGSIb3DQEBCwUAA4IBAQCE9Ll2+jOkHpwNNSDbzXlBkniqja83FcfsPRUoWhanCxN+quX+x+2sxoH2ZqtLrVlhlhfQ0Nds5Mwh6YvaL3louPnKXXlRG+yGb1eUTikDUHtm1RnbBeIOqiLkeLddR62C3AAP9PBGcmTgur99ECm4ZG5BK4/dQtkliw5PbHKo1PmD2ev5tOn7OQfG1h7sqM4YVDqlvXwosuas8Nxp3c+b0oCuw5/j6AQwLXCcGYtxyaslwPZ3lJGXgqiJYDbt+raRe/Iyyvf8Q7HPssgoUNO6lDy/NOrTcNQc9JljytfTPCsE+q0oQ6rTjXT3/GMHWobmhheajzv5qw2nKCoD4kGk"
  ]
}

{
  "access": [
    {
      "actions": [
        "pull"
      ],
      "name": "library/alpine",
      "parameters": {
        "pull_limit": "100",
        "pull_limit_interval": "21600"
      },
      "type": "repository"
    }
  ],
  "aud": "registry.docker.io",
  "exp": 1774332398,
  "iat": 1774332098,
  "iss": "auth.docker.io",
  "jti": "dckr_jti_H_AdWbcX5GiX3ChT4XBB4p8rB3Y=",
  "nbf": 1774331798,
  "sub": "" # Anonymous auth
}

# Registry Auth Flow

1. Client attempts to begin a push/pull operation with the registry
2. If the registry requires authorization it will return a 401 Unauthorized HTTP response with information
   on how to authenticate.
3. The registry client makes a request to the authorization service for a Bearer Token
4. The authorization service returns an opaque Bearer token representing the clients authorized access.
5. The client retries the original request with the Bearer token embedded in the requests Authorization header
6. The Registry authorizes the client by validating the Bearer token and the claim set embedded
   within it and begins the push/pull session as usual

# Requirements
1. Registry clients can understand and respond to token auth challenges returned by the resource server.
2. An authorization server capable of managing access controls to their resources hosted by any given service.
3. A Docker Registry capable of trusting the authorization server to sign 
