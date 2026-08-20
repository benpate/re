# re — Notes for AI Agents

- **`Reader.Read` advances a cursor and rewinds on EOF.** Each `Read` returns the next bytes (`n, nil`) and signals `0, io.EOF` only when the content is exhausted — at which point the cursor rewinds so the *next* read cycle starts over. This is what makes the same `Reader` re-readable. It also means re-reading depends on draining to EOF; to rewind early (mid-stream), call `Reset()`. A cursor-less implementation that returns EOF on every call silently truncates large inputs and repeats the first chunk — don't reintroduce one.

- **`NewReader` and `NewReaderFromBytes` return `*Reader` (a pointer).** The cursor is mutable state, so the methods use a pointer receiver. Pass the `*Reader` around; copying the value would share the buffer but is not the intended use.

- **`NewReaderFromBytes` retains the slice, it does not copy it.** Mutating the slice after the call mutates the `Reader`'s contents. Pass a slice you own, or copy first.

- **The `http.Body` helpers replace the body with a fresh `NopCloser`.** `ReadRequestBody`/`ReadResponseBody` read the body to completion and swap in a re-readable reader, so downstream handlers can read it again. They guard nil request/response and nil bodies (a nil body is legal, e.g. on a HEAD response) and return an empty slice rather than panicking.

- **`CloneResponse` returns `(http.Response, error)`.** It reads (and restores) the original's body to give the clone an independent copy; a failed read is returned rather than silently producing an empty clone.

- **`CloneResponse` deep-copies `Header`, `Trailer`, and `TransferEncoding`.** `clone := *original` duplicates the scalar fields but only shares a reference to the maps and slices, so setting a header on the "copy" used to rewrite the original's header too — silently, with no error. Those three explicit copies are load-bearing; don't collapse them back into the struct assignment. `http.Header.Clone()` and `append([]string(nil), ...)` both preserve nil, which matters because net/http treats a nil `Header` differently from an empty one. `TestCloneResponse_Independent*` and `TestCloneResponse_NilHeadersStayNil` pin all of it. Note that `Request` and `TLS` remain shared pointers, which is intended — they point at the same transaction.
