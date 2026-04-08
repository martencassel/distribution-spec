
# Endpoint: end-1
# Method: GET
# API Endpoint: /v2/
# Success: 200
# Failure: 404 / 401

# Does the GET /v2/ endpoint have multiple functions ?

## 1. Liveness check
“Is the registry alive and reachable?”

## 2. Protocol Capability Check
“Does this server speak the OCI/Docker Distribution API?”

## 3. Authentication Mechanism Discovery
“What auth method should I use?”

## 4. Token Service Discovery
“Where do I get tokens?”

## Is the registry alive and ready to talk ?
A successful response(200 or 401) means:
1. The registry is reachable
2. The registry speaks the Distribution API

If the registry is misconfigured or down you will get:
1. 500 (server error)
2. connection failure
3. timeout

## Does the registry support the Distribution Spec ?
Yes, if it returns either:
- 200 OK
- 401 Unauthorized  (with a valid WWW-Authenticate header)
No, if it returns:
404 Not Found
501 Not Implemented
This is how clients detect whether they can use OCI push/pull flows.

## Does the registry support the spec but requires authentication?
Yes — indicated by:

```text
401 Unauthorized
WWW-Authenticate: Bearer realm="...",service="..."
```
Its a challenge telling the client how to authenticate.

## Can the client proceed with push or pull operations?
The client can proceed only if:

It receives 200 OK → proceed immediately
It receives 401 Unauthorized → perform auth flow, then proceed

The client must not proceed if:
It receives 404 or 501 → registry does not support the spec
It receives network errors → registry unavailable

## What authentication mechanism shall the client use?
The registry tells the client via the WWW-Authenticate header.

```text
WWW-Authenticate:
Bearer realm="https://auth.example.com/token",service="registry.example.com"
```
This tells us:

Use Bearer token authentication
Token endpoint is in realm
Token is scoped to the service

If the header instead says:

```text
Basic realm="Registry"
```
…then the registry uses Basic Auth.

## Where does the client get tokens?
From the realm parameter in the WWW-Authenticate header.

```text
realm="https://auth.example.com/token"
```
The client performs a token request there.

## What scopes are required ?

Scopes are also communicated via the WWW-Authenticate header, typically during subsequent requests.

```text
WWW-Authenticate: Bearer realm="...",service="...",scope="repository:myapp:pull"
```
Scopes tell the client:
- which repository
- which action (pull, push, delete)

## How does the client learn about the scopes required ?

1. Client probes registry

```text
GET /v2/
401 Unauthorized
WWW-Authenticate: Bearer realm="https://auth.example.com/token",service="registry.example.com"
```

2. Client attempts an operation, e.g. pull a manifest

```text
GET /v2/myapp/manifests/latest

Registry replies:
401 Unauthorized
WWW-Authenticate: Bearer realm="https://auth.example.com/token",
                  service="registry.example.com",
                  scope="repository:myapp:pull"
```

3. Client requests a token from the auth server

```text
GET https://auth.example.com/token?
    service=registry.example.com&
    scope=repository:myapp:pull
```
Auth server returns a token with that scope.

4. Client retries the manifest request with the token

GET /v2/myapp/manifests/latest
Authorization: Bearer <token>
Registry returns 200 OK.



# What headers are relevant ?

WWW-Authenticate: Tells the client what auth mechanism to use, where to get tokens, and later what scopes are required.

(Status Code): Not a header, but the only required signal: 200 or 401 → registry supports spec; 404 → does not.

Docker-Distribution-API-Version: Legacy Docker header. MAY appear. Clients SHOULD NOT depend on it.


