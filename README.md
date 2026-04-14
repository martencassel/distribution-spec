### **Content Negotation**

```http
send GET with Accept: supported_types
ct = response.Content-Type (ignore parameters)
if ct not in supported_types: reject
if manifest.mediaType exists and != ct: reject
accept manifest
```

### **Push Manifest**

```http
# Prepare manifest JSON
manifest = {
    "mediaType": manifest_type,   # MUST be present
    ...                           # config, layers, etc.
}

# Compute digest of exact bytes
digest = sha256(manifest.bytes)

# Determine reference (tag or digest)
reference = <tag_or_digest>

# Prepare headers
content_type = manifest.mediaType          # MUST match mediaType
# MUST NOT include parameters
headers = {
    "Content-Type": content_type
}

# Send PUT request
resp = PUT /v2/<name>/manifests/<reference>
        Headers: headers
        Body: manifest.bytes

# Registry behavior:
# - MUST store exact bytes
# - MUST return 201 Created on success
# - MUST include:
#       Location: <pullable-manifest-url>
#       Docker-Content-Digest: digest

if resp.status == 201:
    server_digest = resp.headers["Docker-Content-Digest"]
    if server_digest != digest:
        reject("digest mismatch")
    return success

# Error cases
if resp.status == 404:
    reject("repository does not exist")

if resp.status == 413:
    reject("manifest too large")

# Other errors
reject("unexpected response")
```

```text
manifest.bytes = encode(manifest)
digest = sha256(manifest.bytes)

PUT /manifests/<ref>
    Content-Type: manifest.mediaType
    Body: manifest.bytes

if resp.status != 201: fail
if resp.digest != digest: fail
```

### **Pushing manifest with subject**

```http
PUT /v2/<name>/manifests/<reference>
Content-Type: <manifest mediaType>
Body: <manifest bytes>
--- response ---
201 Created
OCI-Subject: sha256:<digest-of-subject>
Location: <pullable-manifest-url>
Docker-Content-Digest: sha256:<digest-of-pushed-manifest>
```

```http
push manifest
if no OCI-Subject header:
    index = pull referrers list or empty
    if index.type unexpected: fail
    if digest not in index: add descriptor
    push updated index (optionally with If-Match)
```

```json
{
  "mediaType": "application/vnd.oci.image.manifest.v1+json",

  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:aaa...",
    "size": 123
  },

  "layers": [],

  "subject": {
    "mediaType": "application/vnd.oci.image.manifest.v1+json",
    "digest": "sha256:<<< SUBJECT DIGEST >>>",
    "size": 456
  }
}
```

```http
# Assume: manifest includes a "subject" field
manifest.bytes = encode(manifest)
digest = sha256(manifest.bytes)

# Push manifest normally
resp = PUT /v2/<name>/manifests/<reference>
        Content-Type: manifest.mediaType
        Body: manifest.bytes

# If registry supports subject processing, it MUST set OCI-Subject
oci_subject = resp.headers.get("OCI-Subject")

if oci_subject is not None:
    # Registry handled subject automatically
    return success
```

### **Pushing blobs**

A blob can be uploaded in two ways:

1. **Chunked upload**
   Upload the blob in multiple parts using an upload session.

```text
session = POST /uploads
for chunk in chunks:
    PATCH session with chunk
PUT session?digest=sha256
```

2. **Monolithic upload**
   Upload the entire blob in a single request.

   There are two API patterns for monolithic uploads:
   - **POST → PUT**: Start the upload with `POST`, then complete it with a `PUT` containing the full blob.
   - **Single POST**: Upload the entire blob directly in one `POST` request that includes the digest.

```text
# POST → PUT
session = POST /uploads
PUT session?digest=sha256 with full_blob
```

```text
# or single POST
POST /uploads?digest=sha256 with full_blob
```

### **Chunked Uploads Detailed**

```text
# Start upload session
resp = POST /v2/<name>/blobs/uploads/
upload_url = resp.location

# Send chunks
for chunk in blob.split_into_chunks():
    PATCH upload_url
        Body: chunk
    upload_url = response.location   # server may update the URL

# Finalize upload
PUT upload_url?digest=<sha256-of-blob>
    Body: empty
```


### **API Reference: Blob Upload Operations**

## 1. Start upload session

```http
POST /v2/<name>/blobs/uploads/
→ 202 Accepted
Location: <upload-location>
```
- <upload-location> MUST contain a unique session ID.
- MAY be absolute or relative.
- Client MUST treat it as opaque.

## 2. Chunked upload
Send chunk

```http
PATCH <upload-location>
Content-Type: application/octet-stream
Content-Length: <chunk-size>
<chunk-bytes>
```
- MAY return updated Location: header.

## Finalize

```http
PUT <upload-location>?digest=<sha256>
Content-Length: 0

201 Created
Location: <blob-location>
```

## 3. Monolithic upload (POST → PUT)

```http
POST /v2/<name>/blobs/uploads/
→ 202 Accepted
Location: <upload-location>

PUT <upload-location>?digest=<sha256>
Content-Type: application/octet-stream
Content-Length: <size>

<full-blob>
```

## 4. Monolithic upload (single POST)

```http
POST /v2/<name>/blobs/uploads/?digest=<sha256>
Content-Type: application/octet-stream
Content-Length: <size>

<full-blob>

→ 201 Created
→ Location: <blob-location>
```


### From issues

Most clients use POST → monolithic PATCH → PUT.

