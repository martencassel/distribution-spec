# 303

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
