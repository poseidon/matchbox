package sign

import (
	"bytes"
	"net/http"
)

// signatureResponseWriter buffers response writes, flushes data to create a
// detached signature, and writes to an underlying ResponseWriter according to
// the http.ResponseWriter contract. Wrap an http.ResponseWriter to capture
// subsequent writes and flush to reply with a Signer signature and code.
//
// Note: Header and WriteHeader(http.StatusOK) calls are ignored. The response
// is tranformed by signing and has different StatusOK conditions; this is
// transparent to handlers which write buffered data.
type signatureResponseWriter struct {
	w      http.ResponseWriter
	signer Signer
	header http.Header
	buf    *bytes.Buffer
}

// newSignatureResponseWriter returns an http.ResponseWriter which buffers
// response writes into a detached signature response using the given Signer.
func newSignatureResponseWriter(w http.ResponseWriter, signer Signer) *signatureResponseWriter {
	return &signatureResponseWriter{
		w:      w,
		signer: signer,
		header: make(http.Header),
		buf:    new(bytes.Buffer),
	}
}

// Header returns a Header map which is not used when responding to the HTTP
// connection since buffered data will be transformed into a signature.
func (rw *signatureResponseWriter) Header() http.Header {
	return rw.header
}

// Write buffers data to be signed as part of an HTTP reply.
func (rw *signatureResponseWriter) Write(data []byte) (int, error) {
	// data is buffered, not written. Do not call WriteHeader, the underlying
	// ResponseWriter will do that on first Write if not yet called.
	return rw.buf.Write(data)
}

// WriteHeader sends an HTTP response header with status code. Propagate
// calls with error codes so that signatures of error responses are apparent.
// Ignore StatusOK codes, the signer determines the status code of the
// signature reply.
func (rw *signatureResponseWriter) WriteHeader(code int) {
	// success is determined at the time of signing, ignore StatusOK since the
	// header can only be written once.
	if code != http.StatusOK {
		rw.w.WriteHeader(code)
	}
}

// Flush signs the buffered data and writes the signature to the underlying
// http.ResponseWriter if signing succeeds. Any sign or write errors are
// returned.
func (rw *signatureResponseWriter) Flush() error {
	buf := new(bytes.Buffer)
	err := rw.signer.Sign(buf, rw.buf)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(rw.w)
	return err
}

// SignatureHandler wraps an http.Handler and responds with the signature of the
// response from the wrapped handler using the given signer.
func SignatureHandler(signer Signer, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		srw := newSignatureResponseWriter(w, signer)
		// capture response writes from next Handler(s)
		next.ServeHTTP(srw, req)
		// flush buffered data through Signer to underlying http.ResponseWriter
		if err := srw.Flush(); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
	return http.HandlerFunc(fn)
}
