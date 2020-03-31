package id

import (
	"crypto/ed25519"
	"encoding/base64"
)

func NewKeys() *Keys {
	return &Keys{}
}

type Keys struct{}

func (k *Keys) GenerateKeypair() (PublicKey, PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, err
	}

	return PublicKey(pub), PrivateKey(priv), nil
}

type PublicKey ed25519.PublicKey

// Encode returns a human-readble public identifier for the key
// CrevIDs are encoded as URL safe without padding:
// https://github.com/crev-dev/cargo-crev/blob/69cbc57d9c903663c484ba03ad29c2950058b5ef/crev-common/src/lib.rs#L51
func (p PublicKey) Encode() string {
	return base64.RawURLEncoding.EncodeToString(p)
}

type PrivateKey ed25519.PrivateKey
