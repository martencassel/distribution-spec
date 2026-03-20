# 1. A registry hosts images

- A place where images live
- A place you can push images to
- A place you can pull images from

# 2. Users pull images using image addresses

The image address is the user-facing identifier

<registry>/<namespace>/<repository>:<tag>

# 3. The address uniquely identifies the image they want.

# 4. The registry is just the host part of that address.

# 5. The tag is the human‑friendly version selector.

# 6. User Entity Set

- Registry	A server that stores images.	They log in to it and pull from it
- Image		A runnable artifact		They run it with docker run or similiar
- Image Address A pointer to an image		They use it in docker pull

# 7. Deeper OCI model

- Image -> actually a manifest + config + layers
- Image Address -> actually a repository + tag -> manifest digest
- Registry -> actually a distribution API endpoint

# Concepts to OCI

Registry - A server that hosts images - OCI Distribution Endpoint.
Image - A runnable container image - Manifest + Config + Layers - An image is a graph of OCI objects.
Image Address - A string used in CLI - Registry + Repository + Tag -> Manifest - The address that resolves manifest via tag
Tag - A human friendly version label - Tag -> Manifest Descriptor - Tags are pointers to manifests
Digest - A unique ID for an image - Content Digest - Identify immutable content (manifest, layer, configs)
Image Layers - Filesystem layer - Blob Objects - Layers are stored as blobs, referenced by manifest
Image Config - Metadata about the image - Config Blob - Contains env vars, entrypoint etc
Image Variant / SBOM / Signatures - Extra metadata - OCI Artifact / Referer Manifest - Manifests that reference another 
Pulling an Image - Downloading an image - GET manifest -> GET referenced blobs - Client resolves tag, manifest, then blobs
Pushing an Image - Uploading an image - PUT blobs -> PUT manifest - Blobs are uploaded first, then manifest references them

#

Pull ghcr.io/my/app:1.2.3

OCI does:

1. Resolve repository ghcr.io/my/app
2. Resolve tag: 1.2.3 → manifest digest
3. Fetch manifest
4. Fetch config blob
5. Fetch layer blobs
6. Optionally fetch referrers (signatures, SBOMs, etc.)

