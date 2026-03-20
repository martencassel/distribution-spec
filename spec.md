# Pull

## Overview
1. The process of pulling an object centers around retrieving two components: the manifest and one or more blobs.
2. Typically, the first step in pulling an object is to retrieve the manifest.
3. However, you MAY retrieve content from the registry in any order.

## Pulling Manifests
Reference: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests

### 1. Request URL
To pull a manifest, perform a `GET` request to:

```text
/v2/<name>/manifests/<reference>
```

### 2. Name And Reference
1. `<name>` refers to the namespace of the repository.
2. `<reference>` MUST be either:
   - the digest of the manifest, or
   - a tag.
3. `<reference>` MUST NOT be in any other format.

### 3. `<name>` Format
Throughout this document, `<name>` MUST match:

```text
[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*(\/[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*)*
```

Equivalent form:

```text
<name> = match path-segment (\/ path-segment )*
```

### 4. Name Length Compatibility
1. Many clients impose a 255-character limit on:
   - registry hostname (and optional port)
   - `/`
   - `<name>`
2. Examples:

```text
registry-1.docker.io:443/v2/library/alpine
registry-1.docker.io:443/v2/one/two/three/four/.../manifests/latest
```

3. If the registry name is `registry.example.org:5000`, those clients would be limited to a `<name>` of 229 characters (255 minus 25 for host+port and minus 1 for `/`).
4. For compatibility, registries should avoid `<name>` values that exceed this limit.

### 5. Tag Format
Throughout this document, `<reference>` as a tag:
1. MUST be at most 128 characters in length.
2. MUST match:

```text
[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}
```

Examples:

```text
latest
oijsdaijoasdoijsdajiodsajiosadoisdajoisdajoioisjasajdoi
aA1bB90aA_aA9._-912921391239
```

### 6. Content Negotiation
1. The client SHOULD include an `Accept` header indicating supported manifest content types.
2. In a successful response, the `Content-Type` header indicates the type of the returned manifest.
3. The registry SHOULD NOT include parameters on the `Content-Type` header.
4. The client SHOULD ignore parameters on the `Content-Type` header.
5. The `Content-Type` header SHOULD match what the client pushed as the manifest's `Content-Type`.
6. If the manifest has a `mediaType` field, clients SHOULD reject unless `mediaType` matches `Content-Type`.
7. For more information, see Content Negotiation and RFC7231.

Example request:

```http
GET /v2/library/alpine
Accept:
```

### 7. Manifest Response Rules
1. A `GET` request to an existing manifest URL MUST return `200 OK` with the expected manifest.
2. A successful response MUST contain `Docker-Content-Digest`.
3. The `Docker-Content-Digest` value MAY differ from the provided digest.
4. If it differs, hashing algorithms may not match.
5. See Content Digests (Appendix 3) for detecting the hashing algorithm.
6. Most clients MAY ignore this value, but if used, clients MUST verify it matches the returned manifest.
7. If `<reference>` is a digest, clients SHOULD verify the returned manifest matches it.
8. If the manifest is not found, the response code MUST be `404 Not Found`.

### 8. Exercises
Manifest pull examples:

```bash
curl http://registry.local:8080/v2/<name>/manifests/<reference>
curl https://registry-1.docker.io/v2/library/alpine/manifests/latest
```

Valid `<name>` examples:

```text
library/alpine
one/two/three/four/repo
alpine
```

Reference rule:

```text
Reference = digest of the manifest OR a tag
```

## Definitions

### Digest
Digest: a unique identifier created from a cryptographic hash of a blob's content. Digests are defined under the OCI Image Spec (Appendix 3).

## Pulling Blobs

### 1. Request URL
To pull a blob, perform a `GET` request to:

```text
/v2/<name>/blobs/<digest>
```

### 2. Blob Pull Rules
1. `<name>` is the namespace of the repository.
2. `<digest>` is the blob digest.
3. A `GET` request to an existing blob URL MUST return `200 OK` with the expected blob.
4. A successful response MUST contain `Docker-Content-Digest`.
5. If present, `Docker-Content-Digest` MUST match the response body digest.
6. Most clients MAY ignore this value, but if used, clients MUST verify it matches the response body.
7. Clients SHOULD verify the response body matches the requested digest.
8. If the blob is not found, the response code MUST be `404 Not Found`.
9. A registry SHOULD support the `Range` request header in accordance with RFC 9110.

## Checking If Content Exists In The Registry

### Registry Pseudocode
```text
HandleManifestHead(/v2/<name>/manifests/<reference>) {
    if !manifest.exists(name, reference) {
        return 404 Not Found.
    }
    manifest = uploaded.manifest.get(name, reference)
    Header.Docker-Content-Digest = manifest.digest
    Header.Content-Length = blob.Length
    return 200 OK
}

HandleBlobHead(/v2/<name>/blobs/<digest>) {
    if !blobs.exists(name, digest) {
        return 404 Not Found.
    }
    blob = uploaded.blobs.get(name, digest)
    Header.Docker-Content-Digest = blob.digest
    Header.Content-Length = blob.Length
    return 200 OK
}
```

### Client Pseudocode
```text
Check-If-ContentExists(registry, content) {
    # 1. Check if manifest exists

    # 2. Check if blob exists
}
```

### HEAD Rules
1. To verify a repository contains a manifest or blob, make a `HEAD` request to:

```text
/v2/<name>/manifests/<reference>  (for manifests)
/v2/<name>/blobs/<digest>         (for blobs)
```

2. A `HEAD` request to an existing blob or manifest URL MUST return `200 OK`.
3. A successful response MUST contain `Docker-Content-Digest` for the uploaded blob or manifest.
4. A successful response MUST contain `Content-Length` in bytes for the uploaded blob or manifest.
5. Implementers note: Clients may encounter registries implementing earlier spec versions that did not require `Docker-Content-Digest`.
6. In such cases, clients can reasonably assume `sha256`.
7. If the blob or manifest is not found, the response code MUST be `404 Not Found`.

## Appendix 3: Content Digests
Reference: https://github.com/opencontainers/image-spec/blob/v1.0.1/descriptor.md#digests

1. The `digest` property of a Descriptor acts as a content identifier.
2. This enables content addressability.
3. It uniquely identifies content by taking a collision-resistant hash of bytes.
4. If the digest is communicated securely, one can verify content from an insecure source by recalculating it independently.
5. The digest value is a string with:
   - an algorithm portion
   - an encoded portion
6. The algorithm specifies the cryptographic hash function and encoding; the encoded portion contains the encoded hash output.

## Error Codes

A `4XX` response from the registry MAY return a body in any format.
If the response body is JSON, it MUST have the following format:

```json
{
  "errors": [
    {
      "code": "<error identifier, see below>",
      "message": "<message describing condition>",
      "detail": "<unstructured>"
    }
  ]
}
```

1. The `code` field MUST be a unique identifier containing only uppercase alphabetic characters and underscores.
2. The `message` field is OPTIONAL; if present, it SHOULD be human-readable and MAY be empty.
3. The `detail` field is OPTIONAL and MAY contain arbitrary JSON data useful for resolving the issue.
4. The `code` field MUST be one of the following:

| ID | Code | Description |
| --- | --- | --- |
| code-1 | BLOB_UNKNOWN | blob unknown to registry |
| code-2 | BLOB_UPLOAD_INVALID | blob upload invalid |
| code-3 | BLOB_UPLOAD_UNKNOWN | blob upload unknown to registry |
| code-4 | DIGEST_INVALID | provided digest did not match uploaded content |
| code-5 | MANIFEST_BLOB_UNKNOWN | manifest references a manifest or blob unknown to registry |
| code-6 | MANIFEST_INVALID | manifest invalid |
| code-7 | MANIFEST_UNKNOWN | manifest unknown to registry |
| code-8 | NAME_INVALID | invalid repository name |
| code-9 | NAME_UNKNOWN | repository name not known to registry |
| code-10 | SIZE_INVALID | provided length did not match content length |
| code-11 | UNAUTHORIZED | authentication required |
| code-12 | DENIED | requested access to the resource is denied |
| code-13 | UNSUPPORTED | the operation is unsupported |
| code-14 | TOOMANYREQUESTS | too many requests |
