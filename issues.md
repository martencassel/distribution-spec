# 303 - Streamed Upload

1. The spec declares that there are two modes of blob upload: monolithic and chunked.

2. It goes on to describe that a chunked upload is done via PATCH.
It states that the PATCH will include the Content-Range header.

3. However, the conformance tests seem to suggest a third type of blob upload: streamed [1]
This test case issues a single PATCH without Content-Range, followed by a closing PUT

4. It is unclear what this conformance test is for.
I see the existence of a "streamed" upload specification previously[2], but I am not sure when it was removed.
If "streamed" upload is indeed part of v1.0.0, the spec is not clear how it should be conducted.

5. podman also just send Transfer-Encoding: chunked requests (at least in some cases)

6. or requests with Content-Length but withouth Content-Range.

7. I don't understand why the OCI spec created its own "chunked" upload,
using non-standard (at least to the HTTP RFCs) Content-Range values and still requiring ordered upload
What's the point then, given that Transfer-Encoding: chunked exists... I wonder if this can be simplified in a future version of the specification.

8. In practice, barely anybody uses chunked uploads? I see mostly POST + monolithic PATCH + PUT.

9. Doing multiple PATCH requests is slower and less robust. Also more work to implement

10. e.g. python's builtin urllib abstracts Transfer-Encoding: chunked away -- it happens transparently on a request with data, whereas this custom OCI chunked upload has to be implemented by hand (typically at the cost of performance, since urllib doesn't keep the connection alive between requests).

# 485 - Blob Uploads

1. We're proposing a new header that could be useful for registries that are behind proxies
2. Those proxies usually buffer requests so for big layers it degrades performance and memory usage of the server.
3. It would be nice if the registry, either on upload creation (POST /v2/:name/blobs/uploads/) or in another way to hint the client the preferred max chunk size.
4. For example, I've seen places that document the header OCI-Chunk-Min-Length, the purpose is hinting a minimum chunk size for clients. It would be amazing to see the counter-part OCI-Chunk-Max-Length.

5. Does this affect Docker's streaming method for a blob push

6. Streaming method: (not defined by OCI yet, but since it's used by docker it's well supported by registries),
   non-chunked pushes (used by a lot of other tools), or the final chunk size?

7. This affects mainly the POST and PATCH method on chunked uploads. It would work the same as OCI-Chunk-Min-Length, but by defining a maximum chunk size to send in PATCH to not overwhelm the server.
I sent a PR to the distribution repository showcasing how this change would look like.

8. Docker's streaming API uses a patch request without a content-length value. It's effectively the same as a post/put, but the digest is not known in advance, so they push the full blob with a single patch.
How do proxies deal with very large blobs uploaded with a post/put?

9. For proxies or platforms like Workers they buffer the request body. So it would be nice to tell the client (like Docker) to limit the request body to a certain size, as there might be a limit in the server on how many bytes it can buffer.

10. I believe ordered by popularity, there's:
    1. POST/PUT
    2. Docker's streaming PATCH
    3. Chunked POST/PATCH/PUT
    4. Single POST

11. The only method covered by this PR would be 3.
12. OCI hasn't standardized 2 so I don't see how this would apply to that.
13. And if we did define the streaming PATCH, I believe it would have to be excluded from this by its definition.

14. For reference, the streaming patch is defined in the distribution project
15. In practice, this was implemented as a POST / single PATCH (no length or range header) / PUT.
16. I believe this is done for packing the layer content on the fly, where the output of the tar + compress of the files is not known in advance and they didn't want to write to a temp file or memory.

# 578 - What's the expected "GET chunked blob" when the uploaded content is empty

1. A chunked blob upload is accomplished in three phases:
- 1. Obtain a session ID (upload URL) (POST)
     * Content-Length: 0
- 2. Upload the chunks (PATCH)
- 3. Close the session (PUT)

To get the current status after a 416 error, issue a GET request to a URL <location>
MUST have the following headers:
  Location: <location>
  Range: 0-<end-of-range>
