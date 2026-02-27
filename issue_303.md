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

