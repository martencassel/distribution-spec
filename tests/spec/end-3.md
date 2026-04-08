Endpoint: end-3
Method: GET / HEAD
API Endpoint: /v2/<name>/manifests/<reference>
Success: 200
Failure: 404

# Does the endpoint have multiple functions ?

# Function 1 - Retrieve a Manifest (GET)
Return the full manifest document.
Returns full JSON body.

# Function 2 - Check Manifest Existence / Metadata (HEAD)
This is a different function even though it’s the same endpoint.

# Function 3 — Trigger Auth Scope Discovery
This is a different function even though it’s the same endpoint.

# Function 4 — Content Negotiation (Media Type Selection)
Allow clients to request specific manifest types.

# Function 5 — Reference Resolution
Tag to digest resolution.
The reference identifier could either be a tag or a digest.
The Registry must resolve the reference tag to a digest.
The Registry must return the manifest for the digest.

# Function 6 — Conditional Requests (Optional)
“Only give me the manifest if it has changed.”
“Give me the manifest only if the digest does NOT match this value.”

GET /v2/myapp/manifests/latest
If-None-Match: "sha256:abc123"

GET /v2/myapp/manifests/latest
If-Match: "sha256:abc123"


