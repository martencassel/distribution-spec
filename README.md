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

