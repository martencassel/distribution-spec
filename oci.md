Object
	A piece of stored content in a registry, retrievable by digest, optionally referenced by other objects.

Object storage 
	A storage system where:
	1. data is stored as object
	2. each object has a unique identifier
	3. objects are immutable
	4. obects are retrieved by ID, not by path
	5. objects have no hiearchical filesystem structure

Content-Addressed Object Store
	1. It stores immutable objects
	2. Each identifier by digest of its byte
	3. Retrievable via Distribution APi
	4. It only stores and serves bytes

Object Types

	1. Blob. Stored as: opaque bytes. Retrieved by: Digest. Meaning: Layer, configs, SBOM, signatures
	2. Manifest. Stores as: JSON. Retrieved by: Digest. Meaning: Describes an artifact instance.
	3. Index. Stored as: JSON. Retrieved by: Digest. Meaning: Groups manifests (multi-arch).

Filesystem
	1. Path identifies file.		Digest identifies object
	2. Files can change.			Objects are immutable
	3. Hiearchical				Flat namespace
	3. Metadata is external			Metada is part of the obejct

Git
1. Git stores by hash(content)
2. Objects stores by digest(content)

Gives us:
- immutability
- deduplication
- integrity checking
- caching by digest
- reproducible builds

Merkle DAG
Git objects form a DAG
commit -> tree -> blobs

OCI objects form a merkle DAG
manifest -> config blob + layer blobs
index -> manifests
referrers -> manifests

# Immutable objects  + mutable references

Git:
objects = immutable
refs/tags = mutable pointers

OCI:
objects (blobs, manfests, indexes) = immutable
tags = mutable pointersi

# Dedup by content

Git deduplicates identical objects automatically because the hash is the identity.
OCI does the same:
identical layers across images are stored once
identical configs are stored once
identical SBOMs/signatures are stored once

# Integrity and verification
Git verifies object integrity by recomputing the hash.

OCI does the same:
clients verify manifests and blobs by digest
signatures bind to digests
SBOMs reference digests
provenance references digests

# Reproducable
Git guarantees that a commit hash always resolves to the same content.
OCI guarantees that a digest always resolves to the same artifact.
This is essential for:
secure supply chain
reproducible builds
deterministic deployments
signature verification

#

ID(x) = SHA256(x)


Concept	Git	OCI	Notes
Object identity	sha1(content)	sha256(content)	Both content‑addressed
Object types	blob, tree, commit	blob, manifest, index	Structural categories
Reference	tree entry	descriptor	Typed pointer to another object
Pointer metadata	mode, name	mediaType, annotations	Both describe the target
Graph	Merkle DAG	Merkle DAG	Identical structure
Mutable pointer	ref	tag	Name → immutable object

