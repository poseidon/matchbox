package sign

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/openpgp"
)

func TestLoadGPGSigner(t *testing.T) {
	signer, err := LoadGPGSigner("fixtures/secring.gpg", "test")
	assert.Nil(t, err)
	// assert that:
	// - fixture private key is read from a key ring file
	// - fixture encrypted private key is decrypted by passphrase
	// - Signer signs a message which can be verified by OpenPGP
	expectedMessage := "Hello World!"
	signature := new(bytes.Buffer)
	fmt.Println(signature)
	err = signer.Sign(signature, strings.NewReader(expectedMessage))
	assert.Nil(t, err)
	// valid signature
	// gpg --no-default-keyring --secret-keyring fixtures/secring.gpg --verify sig msg
	kring, err := os.Open("fixtures/secring.gpg")
	assert.Nil(t, err)
	defer kring.Close()
	entities, err := openpgp.ReadKeyRing(kring)
	assert.Nil(t, err)
	_, err = openpgp.CheckArmoredDetachedSignature(entities, strings.NewReader(expectedMessage), signature)
	assert.Nil(t, err)
}

func TestLoadGPGSigner_MissingKeyRing(t *testing.T) {
	_, err := LoadGPGSigner("", "")
	assert.NotNil(t, err)
}

func TestLoadGPGSigner_MissingPassphrase(t *testing.T) {
	_, err := LoadGPGSigner("fixtures/secring.gpg", "")
	assert.Equal(t, errMissingPassphrase, err)
}

func TestLoadGPGSigner_IncorrectPassphrase(t *testing.T) {
	_, err := LoadGPGSigner("fixtures/secring.gpg", "incorrect")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "private key checksum failure")
	}
}

// upperSigner "signs" messages by writing a signature that is the upper case
// form of the message body. For testing purposes only.
type upperSigner struct{}

func (s *upperSigner) Sign(w io.Writer, message io.Reader) error {
	b, err := ioutil.ReadAll(message)
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
