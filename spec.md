# Pull

# PULL.1.1 The process of pulling an object centers around retrieving two components: the manifest and one or more blobs.
# PULL.1.2 Typically, the first step in pulling an object is to retrieve the manifest.
# PULL:1.3 However, you MAY retrieve content from the registry in any order.

# 1. Pulling manifests (https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests)

# 1.1
To pull a manifest, perform a GET request to a URL in the following form:
/v2/<name>/manifests/<reference> end-3

# 1.2
<name> refers to the namespace of the repository.

# 1.3
<reference> MUST be either (a) the digest of the manifest or (b) a tag
he <reference> MUST NOT be in any other format.

# 1.4
Throughout this document, <name> MUST match the following regular expression:

[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*(\/[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*)*

[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*
(\/
[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*
)*

<name> = match path-segment (\/ path-segment )*

# 1.5
Many clients impose a limit of 255 characters on the length of the concatenation of the registry hostname
(and optional port), /, and <name> value

registry-1.docker.io:443/v2/library/alpine
registry-1.docker.io:443/v2/one/two/three/four/.../manifests/latest

# 1.6
If the registry name is registry.example.org:5000, those clients would be limited to a <name> of 229 characters (255 minus 25 for the registry hostname and port and minus 1 for a / separator).

# 1.7
For compatibility with those clients, registries should avoid values of <name> that would cause this limit to be exceeded.

# 1.8
Throughout this document, <reference> as a tag MUST be at most 128 characters in length and MUST match the following regular expression: [a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}

latest
oijsdaijoasdoijsdajiodsajiosadoisdajoisdajoioisjasajdoi
aA1bB90aA_aA9._-912921391239

# 1.9
The client SHOULD include an Accept header indicating which manifest content types it supports.
In a successful response, the Content-Type header will indicate the type of the returned manifest.

HTTP GET /v2/library/alpine
Accept:

# 1.10
The registry SHOULD NOT include parameters on the Content-Type header.

# 1.11
The client SHOULD ignore parameters on the Content-Type header.

# 1.12
The Content-Type header SHOULD match what the client pushed as the manifest's Content-Type.

# 1.13
If the manifest has a mediaType field, clients SHOULD reject unless the mediaType field's
value matches the type specified by the Content-Type header.

# 1.14
For more information on the use of Accept headers and content negotiation, please see Content Negotiation and RFC7231.

# 1.15
A GET request to an existing manifest URL MUST provide the expected manifest, with a response code that MUST be 200 OK.

# 1.16
A successful response MUST contain the digest of the uploaded blob in the header Docker-Content-Digest.


# 1.17
The Docker-Content-Digest header, if present on the response, returns the digest of the uploaded blob which MAY differ from the provided digest.

# 1.18
If the digest does differ, it MAY be the case that the hashing algorithms used do not match.

# 1.19
See Content Digests apdx-3 for information on how to detect the hashing algorithm in use.

# 1.20
Most clients MAY ignore the value, but if it is used, the client MUST verify the value matches the returned manifest.

# 1.21
If the <reference> part of a manifest request is a digest,
clients SHOULD verify the returned manifest matches this digest.

# 1.22
If the manifest is not found in the repository, the response code MUST be 404 Not Found.

# Exercises

# 1.1
curl http://registry.local:8080/v2/<name>/manifests/<reference>
curl https://registry-1.docker.io/v2/library/alpine/manifests/latest

# 1.2
library/alpine
one/two/three/four/repo
alpine

# 1.3
Reference = digest of the manifest OR or a tag

# Definitions

# 1.1 Digest
Digest: a unique identifier created from a cryptographic hash of a Blob's content. Digests are defined under the OCI Image Spec apdx-3

# Pulling blobs

# 1.1
To pull a blob, perform a GET request to a URL in the following form: /v2/<name>/blobs/<digest> (end-2)

# 1.2
<name> is the namespace of the repository, and <digest> is the blob's digest.

# 1.3
A GET request to an existing blob URL MUST provide the expected blob, with a response code that MUST be 200 OK.

# 1.4
A successful response MUST contain the digest of the uploaded blob in the header Docker-Content-Digest.

# 1.5
If present, the value of this header MUST be a digest matching that of the response body.

# 1.6
Most clients MAY ignore the value, but if it is used, the client MUST verify the value matches the returned response body.

# 1.7
Clients SHOULD verify that the response body matches the requested digest.

# 1.8
If the blob is not found in the repository, the response code MUST be 404 Not Found.

# 1.9
A registry SHOULD support the Range request header in accordance with RFC 9110.

# Checking if content exists in the registry

Registry:

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

Client:

Check-If-ContentExists(registry, content) {
    # 1. Check if manifest exists

    # 2. Check if blob exists
}

# 1.1
In order to verify that a repository contains a given manifest or blob, make a HEAD request to a URL in the following form:

# 1.2
/v2/<name>/manifests/<reference> end-3 (for manifests), or

# 1.3
/v2/<name>/blobs/<digest> end-2 (for blobs).

# 1.4
A HEAD request to an existing blob or manifest URL MUST return 200 OK.

# 1.5
A successful response MUST contain the digest of the uploaded blob or manifest in the header Docker-Content-Digest

# 1.6
A successful response MUST contain the size in bytes of the uploaded blob or manifest in the header Content-Length.

# 1.7
Implementers note: Clients may encounter registries implementing earlier spec versions which did not require the Docker-Content-Digest header.

# 1.8
In such cases, the clients can reasonably assume the digest algorithm used is sha256.

# 1.9
If the blob or manifest is not found in the repository, the response code MUST be 404 Not Found.




# apdx-3
https://github.com/opencontainers/image-spec/blob/v1.0.1/descriptor.md#digests

# 1.1
The digest property of a Descriptor acts as a content identifier.

# 1.2
The digest property  of 1.1 enables content addressability.

# 1.3
It uniquely identifies content by taking a collision-resistant hash of the bytes.

# 1.4
If the digest can be communicated in a secure manner,
one can verify content from an insecure source by recalculating the digest independently, ensuring the content has not been modified.

# 1.5
The value of the digest property is a string consisting of an algorithm portion and an encoded portion.

# 1.6
The algorithm specifies the cryptographic hash function and encoding used for the digest;
the encoded portion contains the encoded result of the hash function.

# Error Codes

A 4XX response code from the registry MAY return a body in any format.
If the response body is in JSON format, it MUST have the following format:

```json
    {
        "errors": [
            {
                "code": "<error identifier, see below>",
                "message": "<message describing condition>",
                "detail": "<unstructured>"
            },
            ...
        ]
    }
```
1.1 The code field MUST be a unique identifier, containing only uppercase alphabetic characters and underscores
1.2 The message field is OPTIONAL, and if present, it SHOULD be a human readable string or MAY be empty.
1.3 The detail field is OPTIONAL and MAY contain arbitrary JSON data providing information the client can use to resolve the issue.
1.4 The code field MUST be one of the following:

1.5
ID = code-1
Code = BLOB_UNKNOWN
Description = blob unknown to registry

1.6
ID = code-2
Code = BLOB_UPLOAD_INVALID
Description = blob upload invalid

1.7
ID = code-3
Code = BLOB_UPLOAD_UNKNOWN
Description = blob upload unknown to registry

1.8
ID = code-4
Code = DIGEST_INVALID
Description = provided digest did not match uploaded content

1.9
ID = code-5
Code = MANIFEST_BLOB_UNKNOWN
Description = manifest references a manifest or blob unknown to registry

1.10
ID = code-6
Code = MANIFEST_INVALID
Description = manifest invalid

1.11
ID = code-7
Code = MANIFEST_UNKNOWN
Description = manifest unknown to registry

1.12
ID = code-8
Code = NAME_INVALID
Description = invalid repository name

1.13
ID = code-9
Code = NAME_UNKNOWN
Description = repository name not known to registry

1.14
ID = code-10
Code = SIZE_INVALID
Description = provided length did not match content length

1.15
ID = code-11
Code = UNAUTHORIZED
Description = authentication required

1.16
ID = code-12
Code = DENIED
Description = requested access to the resource is denied

1.17
ID = code-13
Code = UNSUPPORTED
Description = the operation is unsupported

1.18
ID = code-14
Code = TOOMANYREQUESTS
Description = too many requests
