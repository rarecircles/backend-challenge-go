package eth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptoRand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"

	"github.com/btcsuite/btcd/btcec"
	"golang.org/x/crypto/sha3"
)

// KeyBag holds private keys in memory, for signing transactions.
type KeyBag struct {
	Keys []*PrivateKey `json:"keys"`
}

func NewKeyBag() *KeyBag {
	return &KeyBag{
		Keys: make([]*PrivateKey, 0),
	}
}

type PublicKey struct {
	inner ecdsa.PublicKey
}

func (p PublicKey) Address() Address {
	return pubkeyToAddress(p.inner)
}

type PrivateKey struct {
	inner *ecdsa.PrivateKey
}

func NewRandomPrivateKey() (*PrivateKey, error) {
	return generatePrivateKey(cryptoRand.Reader)
}

func generatePrivateKey(random io.Reader) (*PrivateKey, error) {
	if random == nil {
		random = cryptoRand.Reader
	}

	privateKey, err := ecdsa.GenerateKey(btcec.S256(), random)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{inner: privateKey}, nil
}

func NewPrivateKey(rawPrivateKey string) (*PrivateKey, error) {
	keyBytes, err := hex.DecodeString(rawPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	return privateKeyFromRawBytes(keyBytes)
}

func privateKeyFromRawBytes(privateKeyBytes []byte) (*PrivateKey, error) {
	if len(privateKeyBytes) != btcec.PrivKeyBytesLen {
		return nil, fmt.Errorf("not enough bytes, got %d bytes but secp256k1 private key must have %d bytes",
			len(privateKeyBytes), btcec.PrivKeyBytesLen)
	}

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	return &PrivateKey{inner: (*ecdsa.PrivateKey)(privKey)}, nil
}

func (p *PrivateKey) String() string {
	return hex.EncodeToString(p.Bytes())
}

func (p *PrivateKey) Bytes() (out []byte) {
	byteCount := p.inner.Params().BitSize / 8
	if p.inner.D.BitLen()/8 >= byteCount {
		return p.inner.D.Bytes()
	}

	return p.inner.D.FillBytes(make([]byte, byteCount))
}

func (p *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return p.inner
}

func (p *PrivateKey) MarshalJSON() ([]byte, error) {
	// The `p.String()` is guaranteed to returns only hex characters, so it's safe to wrap directly with `"` symbols
	return []byte(`"` + p.String() + `"`), nil
}

func (p *PrivateKey) UnmarshalJSON(v []byte) (err error) {
	var s string
	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}

	newPrivKey, err := NewPrivateKey(s)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	*p = *newPrivKey
	return
}

func (p *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{inner: p.inner.PublicKey}
}

type keccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := sha3.NewLegacyKeccak256().(keccakState)
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

func pubkeyToAddress(p ecdsa.PublicKey) Address {
	if p.X == nil || p.Y == nil {
		return nil
	}

	pubBytes := elliptic.Marshal(btcec.S256(), p.X, p.Y)
	return Address(Keccak256(pubBytes[1:])[12:])
}
