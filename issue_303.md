## Formal Summary Document — Discussion on “Streamed Blob Upload not defined by spec”
**Repository:** opencontainers/distribution-spec
**Issue:** #303 — *Streamed Blob Upload not defined by spec*
**URL:** https://github.com/opencontainers/distribution-spec/issues/303
**Date of summary:** 2026-02-27
**Reported by:** martencassel

---

### 1. Purpose of the Discussion
This discussion raises an ambiguity in the **OCI Distribution Specification** regarding supported blob upload mechanisms—specifically whether a **“streamed” blob upload** method exists in the specification and, if so, how it is defined.

---

### 2. Background and Context
The issue author notes that the specification appears to define **two** blob upload modes:

1. **Monolithic upload**
2. **Chunked upload**

The spec’s chunked upload section describes performing uploads via **PATCH** requests and indicates the use of headers such as `Content-Range`.

However, the author observes that the **conformance test suite** appears to validate a **third** behavior that is not clearly described in the spec: a “streamed” upload.

---

### 3. Observations and Evidence

#### 3.1 Spec describes two upload modes
The author cites the spec section stating blob uploads are either monolithic or chunked:
- Reference: https://github.com/opencontainers/distribution-spec/blob/main/spec.md?plain=1#L197

#### 3.2 Chunked upload described as PATCH with Content-Range
The chunked upload narrative implies a PATCH request includes `Content-Range` (and other headers).

The author highlights text that lists expected headers for a PATCH upload chunk:

- `Content-Type: application/octet-stream`
- `Content-Range: <range>`
- `Content-Length: <length>`

But the spec text, as written, may constrain the **values/meaning** of these headers without explicitly stating whether each header is **mandatory** in all cases.

#### 3.3 Conformance tests imply a “streamed” upload mode
The author references a conformance test that performs:
- a **single PATCH** request **without** `Content-Range`, followed by
- a **finalizing PUT** request to close the upload

This appears different from:
- monolithic upload (typically a single PUT with digest), and
- chunked upload (PATCH with explicit chunk ranges)

Reference:
- https://github.com/opencontainers/distribution-spec/blob/main/conformance/02_push_test.go#L21

---

### 4. Core Questions Raised

#### 4.1 Is “streamed blob upload” part of OCI Distribution Spec v1.0.0?
The author asks whether this behavior is:
- an intended part of the spec but under-documented, or
- legacy behavior that should not be tested as conformant, or
- an implementation-permitted variant that the spec should explicitly describe.

#### 4.2 Why does the conformance suite test behavior not clearly defined in the spec?
If the spec does not define streamed uploads, the author questions:
- what the conformance test is asserting, and
- whether the conformance suite is effectively setting requirements beyond the written spec.

#### 4.3 Are chunked-upload headers mandatory, or only conditionally required?
The author points out potential ambiguity: the spec lists headers and constrains their content, but does not unambiguously say whether the headers are strictly required for chunked uploads in all cases.

---

### 5. Historical Note Mentioned
The author cites evidence that a streamed upload specification existed in the past but may have been removed:

- Reference commit: https://github.com/opencontainers/distribution-spec/commit/92e1994a4f13cff06c03f11d86b40ade4ed92730

This is presented as a possible explanation for why the conformance tests include streamed behavior while the current spec text does not.

---

### 6. Impact / Why This Matters
The ambiguity affects:
- **implementers**, who may not know whether streamed uploads are required or optional for conformance;
- **conformance testing**, which may fail otherwise-valid implementations depending on interpretation; and
- **spec clarity**, because normative requirements must be explicit to avoid divergent implementations.

---

### 7. Requested Outcomes / Next Steps Suggested by the Author
The author requests clarification from maintainers/community on:

1. Whether streamed uploads are intended to be supported in v1.0.0.
2. If yes, how streamed uploads should be conducted (normative definition needed).
3. If no, whether the conformance test should be updated to align strictly with the spec.
4. Regardless of streamed upload status, whether the chunked upload section should be clarified to explicitly state which headers are **MUST** vs **SHOULD** vs optional.

The author notes they cannot currently produce a pull request but are willing to help improve either the spec or conformance suite once the intended behavior is clarified.

---

### 8. Key References
- Spec section on blob upload modes:
  https://github.com/opencontainers/distribution-spec/blob/main/spec.md?plain=1#L197
- Conformance test indicating streamed behavior:
  https://github.com/opencontainers/distribution-spec/blob/main/conformance/02_push_test.go#L21
- Historical commit referencing streamed upload spec:
  https://github.com/opencontainers/distribution-spec/commit/92e1994a4f13cff06c03f11d86b40ade4ed92730

---

## Formal Summary — Comments in Issue #303 (“Streamed Blob Upload not defined by spec”)

**Repository:** opencontainers/distribution-spec
**Issue:** #303 — Streamed Blob Upload not defined by spec
**Scope of this summary:** The *comment thread* (responses after the issue description)

### 1. Main Themes Raised in the Comments
The comments converge on three related points:

1. **The “streamed upload” behavior may not be a distinct third mode**, but rather a practical variant of chunked upload used by common clients.
2. **The specification text is not sufficiently precise** about whether `PATCH` requests must include specific headers (notably `Content-Range` and `Content-Length`), especially given real-world client behavior.
3. **There is broader dissatisfaction with the current “chunked upload” design**, including questions about why it diverges from standard HTTP mechanisms such as `Transfer-Encoding: chunked`.

---

### 2. Comment-by-Comment Summary (Chronological)

**(a) jdolitsky (2021-10-20)**
- Agrees that the observed “streamed” behavior appears valid.
- Suggests it is essentially **chunked upload** and notes chunked upload is tested elsewhere.
- Asks maintainer input (tags @jonjohnsonjr) on whether the specific conformance test that triggered the concern should be removed.

**(b) jonjohnsonjr (2021-11-29)**
- States that the behavior appears common across implementations.
- Concludes that, because it is **ubiquitous**, the correct response may be to **document it explicitly** rather than remove it from tests.

**(c) mpreu (2022-12-01)**
- Confirms encountering the same ambiguity when reading the spec.
- Notes that widely used tooling (example: **podman**) sometimes sends:
  - `Transfer-Encoding: chunked`, and/or
  - `Content-Length` without `Content-Range`
- Argues that the spec wording should be **more explicit** on header requirements, particularly since conformance tests appear to allow behavior that the spec does not clearly describe.

**(d) haampie (2023-08-06)**
- Challenges the motivation behind OCI’s custom chunked upload mechanism:
  - It uses **non-standard** `Content-Range` semantics (relative to typical HTTP expectations).
  - It still effectively requires **ordered** uploading, limiting the benefits of “chunking.”
- Observes that in practice most implementations prefer a workflow resembling:
  - `POST` to start + **single monolithic `PATCH`** + `PUT` to finalize
  rather than multiple PATCH requests.
- Notes implementation complexity and performance drawbacks of multiple PATCH requests.
- Asks why `Transfer-Encoding: chunked` works in practice—suggesting servers/frameworks may transparently handle it without registry-specific logic.

---

### 3. Consolidated Conclusion from the Comment Thread
The commenters do not reach a final resolution, but the emerging consensus is:

- The conformance-tested “streamed” pattern is **common in practice** and likely should be **explicitly documented or clarified** in the specification.
- The spec’s requirements for **chunked upload request headers** (especially the presence/necessity of `Content-Range`) are **ambiguous** and should be tightened to match real-world interoperability.
- There is interest in revisiting the overall chunked upload approach in future spec revisions to better align with standard HTTP transfer mechanisms and implementation realities.
