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




# Chunked Uploads Detailed

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
