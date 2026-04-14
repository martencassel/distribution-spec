### Monolithic Blob Upload Methods

https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pushing-a-blob-monolithically

| Method ID   | Description                                      | Steps Required | Notes                                                                 |
|-------------|--------------------------------------------------|----------------|-----------------------------------------------------------------------|
| OP‑UP‑06    | Upload the entire blob in one POST request       | 1 step         | Client sends the whole blob directly in a POST request               |
| OP‑UP‑07    | Upload the blob using POST (get URL) then PUT    | 2 steps        | Client first obtains an upload URL, then uploads the blob with PUT   |


### OP‑UP‑07 — Upload the Blob Using POST (to get a session) and Then PUT (to send the blob)

Reference: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#post-then-put

| Concern              | Specification                                                                 |
|----------------------|-------------------------------------------------------------------------------|
| **ID**               | OP‑UP‑07                                                                      |
| **Logical Operation**| Upload the blob in two steps: first get an upload URL, then send the blob     |
| **Step 1 Method**    | POST                                                                          |
| **Step 1 Endpoint**  | `/v2/<name>/blobs/uploads/`                                                   |
| **Step 1 Request Body** | MUST be empty                                                              |
| **Step 1 Response**  | MUST return `202 Accepted`<br>MUST include `Location: <location>` containing a unique upload session ID |
| **Location Rules**   | `<location>` MAY be absolute or relative<br>MAY point to another server<br>MUST be used exactly as returned (except for absolute/relative conversion) |
| **Step 2 Method**    | PUT                                                                           |
| **Step 2 Endpoint**  | `<location>?digest=<digest>`                                                  |
| **Step 2 Request Body** | MUST contain the full blob byte stream                                     |
| **Step 2 Request Headers** | MUST include `Content-Length: <length>` (exact size of the blob)<br>MUST include `Content-Type: application/octet-stream` |
| **Digest Invariant** | `<digest>` MUST match the actual digest of the blob                           |
| **Step 2 Success Response** | MUST return `201 Created`<br>MUST include `Location: <blob-location>` (a pullable blob URL) |
| **Failure Cases**    | `400 Bad Request` if digest or content length does not match                  |


### OP‑UP‑06 — Upload the Entire Blob in One POST Request

Reference: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#single-post

| Concern              | Specification                                                                 |
|----------------------|-------------------------------------------------------------------------------|
| **ID**               | OP‑UP‑06                                                                      |
| **Logical Operation**| Upload the whole blob in a single POST request                                |
| **HTTP Method**      | POST                                                                          |
| **Endpoint**         | `/v2/<name>/blobs/uploads/?digest=<digest>`                                   |
| **Request Body**     | MUST contain the full blob byte stream                                        |
| **Request Headers**  | MUST include `Content-Length: <length>` (exact size of the blob)<br>MUST include `Content-Type: application/octet-stream` |
| **Digest Invariant** | `<digest>` in the URL MUST match the actual digest of the blob                |
| **Registry Behavior**| If the registry supports single‑POST uploads: it stores the blob immediately   |
| **Success Response** | `201 Created`<br>Response MUST include `Location: <blob-location>` (a pullable URL) |
| **Fallback Behavior**| If the registry does NOT support this upload style: it MUST return `202 Accepted` and a `Location` header for an upload session; the client MUST continue the upload using a final PUT request |
| **Failure Cases**    | `400 Bad Request` if digest or content length does not match                  |


### OP‑UP‑08‑A — Start a Chunked Upload (POST)

| Concern              | Specification                                                                 |
|----------------------|-------------------------------------------------------------------------------|
| **ID**               | OP‑UP‑08‑A                                                                    |
| **Role in Workflow** | Step 1 of chunked upload                                                      |
| **Purpose**          | Ask the registry for an upload URL (session ID)                               |
| **HTTP Method**      | POST                                                                          |
| **Endpoint**         | `/v2/<name>/blobs/uploads/`                                                   |
| **Request Body**     | MUST be empty                                                                 |
| **Request Headers**  | MUST include `Content-Length: 0`                                              |
| **Response**         | MUST return `202 Accepted`                                                    |
| **Response Headers** | MUST include `Location: <location>` (unique upload URL)<br>MAY include `OCI-Chunk-Min-Length: <size>` |
| **Notes**            | `<location>` MUST be used exactly as returned (except absolute/relative conversion) |


### OP‑UP‑08‑B — Upload a Chunk (PATCH)

| Concern              | Specification                                                                 |
|----------------------|-------------------------------------------------------------------------------|
| **ID**               | OP‑UP‑08‑B                                                                    |
| **Role in Workflow** | Step 2 of chunked upload                                                      |
| **Purpose**          | Send a piece of the blob to the upload URL                                    |
| **HTTP Method**      | PATCH                                                                         |
| **Endpoint**         | `<location>` (exactly as returned by OP‑UP‑08‑A)                              |
| **Request Body**     | MUST contain a chunk of the blob                                              |
| **Request Headers**  | MUST include `Content-Type: application/octet-stream`<br>MUST include `Content-Length: <chunk-size>` |
| **Response**         | MUST return `202 Accepted`                                                    |
| **Response Headers** | SHOULD include `Range` showing how many bytes the registry has received       |
| **Notes**            | `<location>` may contain important query parameters; clients MUST NOT rebuild it manually |

### OP‑UP‑08‑C — Finish a Chunked Upload (PUT)

| Concern              | Specification                                                                 |
|----------------------|-------------------------------------------------------------------------------|
| **ID**               | OP‑UP‑08‑C                                                                    |
| **Role in Workflow** | Step 3 of chunked upload                                                      |
| **Purpose**          | Tell the registry that all chunks are uploaded and provide the final digest   |
| **HTTP Method**      | PUT                                                                           |
| **Endpoint**         | `<location>?digest=<digest>`                                                  |
| **Request Body**     | MUST be empty (`Content-Length: 0`)                                           |
| **Request Headers**  | MUST include the correct `digest=<digest>`                                    |
| **Response**         | MUST return `201 Created`                                                     |
| **Response Headers** | MUST include `Location: <blob-location>` (a pullable blob URL)                |
| **Failure Cases**    | `400 Bad Request` if the digest does not match the uploaded content           |


```text
OP‑UP‑08‑A (POST)
    ↓ returns Location
OP‑UP‑08‑B (PATCH)
    ↓ same Location
OP‑UP‑08‑C (PUT with digest)
```
