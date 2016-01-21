package sign

import (
	"errors"
	"io"
	"os"

	"golang.org/x/crypto/openpgp"
)

var (
	errEmptyKeyring      = errors.New("sign: provided key ring file contained no keys")
	errMissingPassphrase = errors.New("sign: missing passphrase for encrypted private key")
)

// A Signer signs messages and writes detached signatures to w.
type Signer interface {
	Sign(w io.Writer, message io.Reader) error
}

// gpgSigner reads messages and writes ascii armored OpenPGP signatures.
type gpgSigner struct {
	signer *openpgp.Entity
}

// Sign signs the given message and writes the ascii armored OpenPGP signature
// to w.
func (s *gpgSigner) Sign(w io.Writer, message io.Reader) error {
	return openpgp.ArmoredDetachSignText(w, s.signer, message, nil)
}

// NewGPGSigner returns a new Signer that reads messages and writes ascii
// armored OpenPGP signatures.
func NewGPGSigner(signer *openpgp.Entity) Signer {
	return &gpgSigner{
		signer: signer,
	}
}

// LoadGPGSigner loads a key ring file, unlocks the first key with the given
// passphrase, and returns a new Signer that reads messages and writes ascii
// armored OpenPGP signatures.
func LoadGPGSigner(keyRingPath, passphrase string) (Signer, error) {
	kring, err := os.Open(keyRingPath)
	if err != nil {
		return nil, err
	}
	defer kring.Close()
	entity, err := unlockKeyRingEntity(kring, passphrase)
	if err != nil {
		return nil, err
	}
	return NewGPGSigner(entity), nil
}

// unlockKeyRingEntity loads a key ring file and returns the first Entity. The
// given passphrase is used to unlock the entity if it has an encrypted private
// key.
func unlockKeyRingEntity(ring io.Reader, passphrase string) (*openpgp.Entity, error) {
	entities, err := openpgp.ReadKeyRing(ring)
	if err != nil {
		return nil, err
	}
	if len(entities) < 1 {
		return nil, errEmptyKeyring
	}
	entity := entities[0]
	if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
		if passphrase == "" {
			return nil, errMissingPassphrase
		}
		if err := entity.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
			return nil, err
		}
	}
	return entity, nil
}
