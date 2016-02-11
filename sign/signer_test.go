package sign

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/openpgp"
)

func TestLoadGPGEntity(t *testing.T) {
	entity, err := LoadGPGEntity("fixtures/secring.gpg", "test")
	assert.Nil(t, err)
	assert.NotNil(t, entity)
}

func TestLoadGPGEntity_MissingKeyring(t *testing.T) {
	_, err := LoadGPGEntity("", "")
	assert.NotNil(t, err)
}

func TestLoadGPGEntity_ReadKeyringError(t *testing.T) {
	_, err := LoadGPGEntity("fixtures/mangled.gpg", "test")
	assert.NotNil(t, err)
}

func TestLoadGPGEntity_EmptyKeyring(t *testing.T) {
	_, err := LoadGPGEntity("fixtures/empty.gpg", "")
	assert.Equal(t, errEmptyKeyring, err)
}

func TestLoadGPGEntity_MissingPassphrase(t *testing.T) {
	_, err := LoadGPGEntity("fixtures/secring.gpg", "")
	assert.Equal(t, errMissingPassphrase, err)
}

func TestLoadGPGEntity_IncorrectPassphrase(t *testing.T) {
	_, err := LoadGPGEntity("fixtures/secring.gpg", "incorrect")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "private key checksum failure")
	}
}

func TestGPGSigner(t *testing.T) {
	entity, err := LoadGPGEntity("fixtures/secring.gpg", "test")
	assert.Nil(t, err)
	// assert that:
	// - fixture private key is read from a key ring file
	// - fixture encrypted private key is decrypted by passphrase
	// - gpgSigner creates a signature which can be verified
	signer := NewGPGSigner(entity)
	expectedMessage := "Hello World!"
	signature := new(bytes.Buffer)
	err = signer.Sign(signature, strings.NewReader(expectedMessage))
	assert.Nil(t, err)
	// valid signature
	// gpg --homedir sign/fixtures --verify sig msg
	kring, err := os.Open("fixtures/secring.gpg")
	assert.Nil(t, err)
	defer kring.Close()
	entities, err := openpgp.ReadKeyRing(kring)
	assert.Nil(t, err)
	_, err = openpgp.CheckDetachedSignature(entities, strings.NewReader(expectedMessage), signature)
	assert.Nil(t, err)
}

func TestArmoredGPGSigner(t *testing.T) {
	entity, err := LoadGPGEntity("fixtures/secring.gpg", "test")
	assert.Nil(t, err)
	// assert that:
	// - fixture private key is read from a key ring file
	// - fixture encrypted private key is decrypted by passphrase
	// - armoredGPGSigner creates an armored signature which can be verified
	signer := NewArmoredGPGSigner(entity)
	expectedMessage := "Hello World!"
	signature := new(bytes.Buffer)
	err = signer.Sign(signature, strings.NewReader(expectedMessage))
	assert.Nil(t, err)
	// valid signature
	// gpg --homedir sign/fixtures --verify sig msg
	kring, err := os.Open("fixtures/secring.gpg")
	assert.Nil(t, err)
	defer kring.Close()
	entities, err := openpgp.ReadKeyRing(kring)
	assert.Nil(t, err)
	_, err = openpgp.CheckArmoredDetachedSignature(entities, strings.NewReader(expectedMessage), signature)
	assert.Nil(t, err)
}
