package sign

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignatureHandler(t *testing.T) {
	signer := new(upperSigner)
	next := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "next handler"}`)
	}
	h := SignatureHandler(signer, http.HandlerFunc(next))
	// assert that:
	// - writes by 'next' handler(s) are buffered and signed
	// - headers set by 'next' handler(s) are ignored
	expectedBody := `{"MESSAGE": "NEXT HANDLER"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())
	assert.NotEqual(t, "application/json", w.Header().Get("Content-Type"))
}

func TestSignatureHandler_ErrorStatusCode(t *testing.T) {
	signer := new(upperSigner)
	next := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "missing some argument"}`)
	}
	h := SignatureHandler(signer, http.HandlerFunc(next))
	// assert that:
	// - error status code headers by 'next' handler(s) are propagated
	// - writes by 'next' handler(s) are buffered and signed
	expectedBody := `{"MESSAGE": "MISSING SOME ARGUMENT"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())
}

func TestSignatureHandler_SignatureError(t *testing.T) {
	errorMessage := "sign: error signing message"
	signer := &errorSigner{
		errorMessage: errorMessage,
	}
	next := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "next handler")
	}
	h := SignatureHandler(signer, http.HandlerFunc(next))
	// assert that:
	// - Flush signing errors return without partial signature writes
	// - SignatureHandler writes a 500 InternalServerError
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(w, req)
	assert.Equal(t, errorMessage+"\n", w.Body.String())
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// upperSigner "signs" messages by writing a signature that is the upper case
// form of the message body. For testing purposes only.
type upperSigner struct{}

func (s *upperSigner) Sign(w io.Writer, message io.Reader) error {
	b, err := io.ReadAll(message)
	if err != nil {
		return err
	}
	signature := strings.ToUpper(string(b))
	_, err = io.Copy(w, bytes.NewReader([]byte(signature)))
	return err
}

// errorSigner always returns an error message.
type errorSigner struct {
	errorMessage string
}

func (s *errorSigner) Sign(w io.Writer, message io.Reader) error {
	return errors.New(s.errorMessage)
}
