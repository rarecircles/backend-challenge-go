package eth

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

type Hex []byte

func MustNewHex(input string) Hex {
	return Hex(mustNewBytes("hex", input))
}

func NewHex(input string) (Hex, error) {
	out, err := newBytes("hex", input)
	if err != nil {
		return nil, err
	}

	return Hex(out), nil
}

func (h Hex) String() string                   { return bytes(h).String() }
func (h Hex) Pretty() string                   { return bytes(h).Pretty() }
func (h Hex) Bytes() []byte                    { return h[:] }
func (h Hex) MarshalText() ([]byte, error)     { return bytes(h).MarshalText() }
func (h Hex) ID() uint64                       { return bytes(h).ID() }
func (h Hex) MarshalJSON() ([]byte, error)     { return bytes(h).MarshalJSON() }
func (h Hex) MarshalJSONRPC() ([]byte, error)  { return bytes(h).MarshalJSONRPC() }
func (h *Hex) UnmarshalJSON(data []byte) error { return (*bytes)(h).UnmarshalJSON(data) }

type Hash []byte

func MustNewHash(input string) Hash {
	return Hash(mustNewBytes("hash", input))
}

func NewHash(input string) (Hash, error) {
	out, err := newBytes("hash", input)
	if err != nil {
		return nil, err
	}

	return Hash(out), nil
}

func (h Hash) String() string                   { return bytes(h).String() }
func (h Hash) Pretty() string                   { return bytes(h).Pretty() }
func (h Hash) Bytes() []byte                    { return h[:] }
func (h Hash) MarshalText() ([]byte, error)     { return bytes(h).MarshalText() }
func (h Hash) ID() uint64                       { return bytes(h).ID() }
func (h Hash) MarshalJSON() ([]byte, error)     { return bytes(h).MarshalJSON() }
func (h Hash) MarshalJSONRPC() ([]byte, error)  { return bytes(h).MarshalJSONRPC() }
func (h *Hash) UnmarshalJSON(data []byte) error { return (*bytes)(h).UnmarshalJSON(data) }

type Address []byte

func MustNewAddress(input string) Address {
	out, err := NewAddress(input)
	if err != nil {
		panic(err)
	}

	return out
}

func NewAddress(input string) (Address, error) {
	out, err := newBytes("address", input)
	if err != nil {
		return nil, err
	}

	byteCount := len(out)
	if byteCount > 20 {
		out = out[byteCount-20:]
	}

	return Address(out), nil
}

func (a Address) String() string                   { return bytes(a).String() }
func (a Address) Pretty() string                   { return bytes(a).Pretty() }
func (a Address) Bytes() []byte                    { return a[:] }
func (a Address) MarshalText() ([]byte, error)     { return bytes(a).MarshalText() }
func (a Address) ID() uint64                       { return bytes(a).ID() }
func (a Address) MarshalJSON() ([]byte, error)     { return bytes(a).MarshalJSON() }
func (a Address) MarshalJSONRPC() ([]byte, error)  { return bytes(a).MarshalJSONRPC() }
func (a *Address) UnmarshalJSON(data []byte) error { return (*bytes)(a).UnmarshalJSON(data) }

// PSQL Definition

func (a Address) Value() (driver.Value, error) {
	return a.Pretty(), nil
}
func (b *Address) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	sv, err := driver.String.ConvertValue(value)
	if err != nil {
		return err
	}

	bs, ok := sv.(string)
	if !ok {
		return fmt.Errorf("could not convert data to byte array")
	}

	addr, err := NewAddress(bs)
	if err != nil {
		return fmt.Errorf("could not convert to eth address")
	}
	*b = addr
	return nil
}

//----------------------------------------------------

type bytes []byte

func mustNewBytes(tag string, input string) bytes {
	out, err := newBytes(tag, input)
	if err != nil {
		panic(err)
	}

	return out
}

func newBytes(tag string, input string) (out bytes, err error) {
	bytes, err := hex.DecodeString(SanitizeHex(input))
	if err != nil {
		return out, fmt.Errorf("invalid %s %q: %w", tag, input, err)
	}

	return bytes, nil
}

func (b bytes) String() string {
	return hex.EncodeToString(b)
}

func (b bytes) Pretty() string {
	return "0x" + hex.EncodeToString(b)
}

func (b bytes) Bytes() []byte {
	return b
}

func (b bytes) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b bytes) ID() uint64 {
	return binary.LittleEndian.Uint64(b)
}

func (b bytes) MarshalJSON() ([]byte, error) {
	return []byte(`"` + hex.EncodeToString([]byte(b)) + `"`), nil
}

func (b bytes) MarshalJSONRPC() ([]byte, error) {
	return []byte(`"` + b.Pretty() + `"`), nil
}

func (b *bytes) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	var err error
	if *b, err = hex.DecodeString(strings.TrimPrefix(s, "0x")); err != nil {
		return err
	}

	return nil
}
