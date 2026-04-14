### OP‑UP‑01 — Start Upload Session

| Concern              | Specification                                                                 |
|----------------------|-------------------------------------------------------------------------------|
| **ID**               | OP‑UP‑01                                                                      |
| **Logical Operation**| Start upload session                                                           |
| **HTTP Method**      | POST                                                                          |
| **Endpoint**         | `/v2/<name>/blobs/uploads/`                                                   |
| **Request Body**     | MUST be empty                                                                 |
| **Request Headers**  | MUST NOT include `Content-Length: > 0`                                        |
| **Response Headers** | MUST include `Location: /v2/<name>/blobs/uploads/<uuid>`<br>MAY include `Range` |
| **Response Codes**   | `202 Accepted`, `401 Unauthorized`, `403 Forbidden`                           |

### OP‑UP‑06 — Single‑Request Monolithic Upload (POST with digest)

| Concern              | Specification                                                                 |
|----------------------|-------------------------------------------------------------------------------|
| **ID**               | OP‑UP‑06                                                                      |
| **Logical Operation**| Upload the entire blob in one POST request                                    |
| **HTTP Method**      | POST                                                                          |
| **Endpoint**         | `/v2/<name>/blobs/uploads/?digest=<digest>`                                   |
| **Request Body**     | MUST contain full blob byte stream                                            |
| **Request Headers**  | MUST include `Content-Length: <length>`<br>MUST include `Content-Type: application/octet-stream`<br>`Content-Length` MUST match actual body size |
| **Digest Invariant** | `<digest>` MUST match the blob’s actual digest                                |
| **Response Headers** | On success: MUST include `Location: <blob-location>`                          |
| **Response Codes**   | `201 Created` on success<br>`202 Accepted` if registry does NOT support single‑POST monolithic uploads (client MUST continue with PUT)<br>`400` for digest mismatch |
| **Notes**            | If registry returns `202`, client MUST follow POST‑then‑PUT upload method     |
