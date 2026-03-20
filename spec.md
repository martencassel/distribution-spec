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

# 1.5
Many clients impose a limit of 255 characters on the length of the concatenation of the registry hostname
(and optional port), /, and <name> value

# 1.6
If the registry name is registry.example.org:5000, those clients would be limited to a <name> of 229 characters (255 minus 25 for the registry hostname and port and minus 1 for a / separator).

# 1.7
For compatibility with those clients, registries should avoid values of <name> that would cause this limit to be exceeded.

# 1.8
Throughout this document, <reference> as a tag MUST be at most 128 characters in length and MUST match the following regular expression: [a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}

# 1.9
The client SHOULD include an Accept header indicating which manifest content types it supports.
In a successful response, the Content-Type header will indicate the type of the returned manifest.

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
