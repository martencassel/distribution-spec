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

# Push

1. Pushing an object typically works in the opposite order as a pull:
2. The blobs making up the object are uploaded first, and the manifest last. A useful diagram is provided here.
3. A registry MUST initially accept an otherwise valid manifest with a subject field that references a manifest that does not exist in the repository, allowing clients to push a manifest and referrers to that manifest in either order.
4. A registry MAY reject a manifest uploaded to the manifest endpoint with descriptors in other fields that reference a manifest or blob that does not exist in the registry.
5. When a manifest is rejected for this reason, it MUST result in one or more MANIFEST_BLOB_UNKNOWN errors code-1.

## Pushing Blobs

BLOB PUSH METHODS = "MONO METHODS" + "CHUNKED" 
MONO METHODS      = "POST then PUT" + "Single POST"

# POST then PUT
To push a blob monolithiclly by using a POST request followed by a PUT request, there are two steps:a

1. Obtain a session id (upload URL)
2. Upload the blob to said URL

# Single POST

1. 
To push a blob monolithically by using a single POST request, perform a POST request to a URL in the following form, and with the following headers and body:

POST /v2/<name>/blobs/uploads/?digest=<digest> end-4b
Content-Lngth: <length>
Content-Type: application/octet-streame
<upload bte stream>

Registries that do not support single request monolithic uploads SHOULD return a 202 Accepted status code and Location header and clients SHOULD proceed with a subsequent PUT request, as described by the POST then PUT upload method.

Successful completion of the request MUST return a 201 Created and MUST include the following header:
Location: <blob-location>


# CHUNKED (Pushing blobs in chunks)


# 1. 
Obtain a session ID (upload URL)

# 2. 
Upload the chunks (PATCH)

# 3. 
Close the session (PUT)


# Cancel a blob upload

During a blob upload, the session may be canceled with a DELETE request:
URL Path: <location>
Content-Length: 0

# Pushing Manifests with Subject

# 1.
When processing a request for an image manifest with the subject field, a registry implementation that supports the referrers API MUST respond with the response header OCI-Subject: <subject digest> to indicate to the client that the registry processed the request's subject.

Client Request:
	Image Manifest = { .subject = ... }

OCI Server (that supports referrer api):
	Header: OCI-Subject: <subject digest> (We ack that we processed image manifest.subjecg

# 2. 
When pushing a manifest with the subject field and the OCI-Subject header was not set, the client MUST:

# 3.
Pull the current referrers list using the referrers tag schema.

# 4.
If that pull returns a manifest other than the expected image index, the client SHOULD report a failure and skip the remaining steps.

# 5.
If the tag returns a 404, the client MUST begin with an empty image index.

# 6
Verify the descriptor for the manifest is not already in the referrers list (duplicate entries SHOULD NOT be created).

# 7
Append a descriptor for the pushed manifest to the manifests in the referrers list. 

# 8
The value of the artifactType MUST be set to the artifactType value in the pushed manifest, if present. 

# 9
If the artifactType is empty or missing in a pushed image manifest, the value of artifactType MUST be set to the config descriptor mediaType value. 

# 10
All annotations from the pushed manifest MUST be copied to this descriptor.

# 11
Push the updated referrers list using the same referrers tag schema. 

# 12
The client MAY use conditional HTTP requests to prevent overwriting a referrers list that has changed since it was first pulled.

# Content Discovery

# Listing Tags
To fetch the list of tags, perform a GET request to a path in the following format: /v2/<name>/tags/list end-8a

# Listing Referrers
To fetch the list of referrers, perform a GET request to a path in the following format: /v2/<name>/referrers/<digest> end-12a.

# 1.1
Upon success, the response MUST be a JSON body with an image index containing a list of descriptors. 

# 1.2
The Content-Type header MUST be set to application/vnd.oci.image.index.v1+json. 

# 1.3
Each descriptor is of an image manifest or index in the same <name> namespace with a subject field that specifies the value of <digest>. 

# 1.4
The descriptors MUST include an artifactType field that is set to the value of the artifactType in the image manifest or index, if present. 

# 1.5
If the artifactType is empty or missing in the image manifest, the value of artifactType MUST be set to the config descriptor mediaType value. 

# 1.6
If the artifactType is empty or missing in an index, the artifactType MUST be omitted. 

# 1.7
The descriptors MUST include annotations from the image manifest or index. 

# 1.8
If a query results in no matching referrers, an empty manifest list MUST be returned.

```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 1234,
      "digest": "sha256:a1a1a1...",
      "artifactType": "application/vnd.example.sbom.v1",
      "annotations": {
        "org.opencontainers.image.created": "2022-01-01T14:42:55Z",
        "org.example.sbom.format": "json"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 1234,
      "digest": "sha256:a2a2a2...",
      "artifactType": "application/vnd.example.signature.v1",
      "annotations": {
        "org.opencontainers.image.created": "2022-01-01T07:21:33Z",
        "org.example.signature.fingerprint": "abcd"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.index.v1+json",
      "size": 1234,
      "digest": "sha256:a3a3a3...",
      "annotations": {
        "org.opencontainers.image.created": "2023-01-01T07:21:33Z",
      }
    }
  ]
}
```

# SUPPORT filtering on artifactType

The registry SHOULD support filtering on artifactType. To fetch the list of referrers with a filter, perform a GET request to a path in the following format: /v2/<name>/referrers/<digest>?artifactType=<artifactType> end-12b. If filtering is requested and applied, the response MUST include a header OCI-Filters-Applied: artifactType denoting that an artifactType filter was applied. If multiple filters are applied, the header MUST contain a comma separated list of applied filters.

Example request with filtering:

GET /v2/<name>/referrers/<digest>?artifactType=application/vnd.example.sbom.v1


# Fallback

If the referrers API returns a 404, the client MUST fallback to pulling the referrers tag schema. The response SHOULD be an image index with the same content that would be expected from the referrers API. If the response to the referrers API is a 404, and the tag schema does not return a valid image index, the client SHOULD assume there are no referrers to the manifest.

# Deleting a Manifest

To delete a manifest, perform a DELETE request to a path in the following format: /v2/<name>/manifests/<digest> end-9

<name> is the namespace of the repository, and <digest> is the digest of the manifest to be deleted. Upon success, the registry MUST respond with a 202 Accepted code. If the repository does not exist, the response MUST return 404 Not Found. If manifest deletion is disabled, the registry MUST respond with either a 400 Bad Request or a 405 Method Not Allowed.

Once deleted, a GET to /v2/<name>/manifests/<digest> and any tag pointing to that digest will return a 404.

When deleting an image manifest that contains a subject field, and the referrers API returns a 404, clients SHOULD:

1. Pull the referrers list using the referrers tag schema.
2. Remove the descriptor entry from the array of manifests that references the deleted manifest.
3. Push the updated referrers list using the same referrers tag schema. The client MAY use conditional HTTP requests to prevent overwriting an referrers list that has changed since it was first pulled.
4. When deleting a manifest that has an associated referrers tag schema, clients MAY also delete the referrers tag when it returns a valid image index.



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

- A registry MAY reject a manifest uploaded to the manifest endpoint with descriptors in other fields
  that reference a manifest or blob that does not exist in the registry.


| code-2 | BLOB_UPLOAD_INVALID | blob upload invalid |


| code-3 | BLOB_UPLOAD_UNKNOWN | blob upload unknown to registry |

- end-5: PATCH blob upload 
- end-6: PUT blob upload
- end-13: GET blob upload
- end-14: DELETE blob upload

| code-4 | DIGEST_INVALID | provided digest did not match uploaded content |
| code-5 | MANIFEST_BLOB_UNKNOWN | manifest references a manifest or blob unknown to registry |
| code-6 | MANIFEST_INVALID | manifest invalid |
| code-7 | MANIFEST_UNKNOWN | manifest unknown to registry |

- end-3, end-7, end-9

| code-8 | NAME_INVALID | invalid repository name |
| code-9 | NAME_UNKNOWN | repository name not known to registry |

| code-10 | SIZE_INVALID | provided length did not match content length |
| code-11 | UNAUTHORIZED | authentication required |

| code-12 | DENIED | requested access to the resource is denied |
| code-13 | UNSUPPORTED | the operation is unsupported |
| code-14 | TOOMANYREQUESTS | too many requests |

# Entity Relationships

## Entities

Repository
	Name

Manifests
	Reference
	Media Type
	Size
	Subject

Tags
	Key: Repository Name
	Key: Name

Referrers
	Repository Name
	Digest
	ArtifactType
	
Blobs
	Repository Name
	Digest
	Size 

## Relationships

Repository HAS Manifests		/v2/<name>/manifests/
Repository HAS Blobs			/v2/<name>/blobs/
Repository HAS Tags			/v2/<name>/tags/
Repository HAS Referrers		/v2/<name>/referrers/

alpine:latest
	SBOM -> alpine:latest	Manifest(SBOM, ArtifactType=SBOM).Subject = alpine.latest
	SIG  -> alpine:latest	Manifest(SIG, ArtifactType=SIG).Subject  = alpine.latest

/<name>/referrers/<digest-of-alpine:latest>
	 "mediaType": "application/vnd.oci.image.index.v1+json",
	 manifests:
	 	{  Manifest(SBOM, ArtifactType=SBOM) }
		{   Manifest(SIG, ArtifactType=SIG).Subject }

## Types

Reference = Digest or Tag

