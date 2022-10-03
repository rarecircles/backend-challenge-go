package eth

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"go.uber.org/zap"
)

type Decoder struct {
	buffer []byte
	offset uint64
	total  uint64
}

func NewDecoderFromString(input string) (*Decoder, error) {
	data, err := NewHex(input)
	if err != nil {
		return nil, fmt.Errorf("unable to decode hex input %q: %w", input, err)
	}

	return NewDecoder(data), nil
}

func NewDecoder(input []byte) *Decoder {
	return &Decoder{
		buffer: input,
		offset: 0,
		total:  uint64(len(input)),
	}
}

func (d *Decoder) String() string {
	return fmt.Sprintf("offset %d, total: %d", d.offset, d.total)
}

func (d *Decoder) SetBytes(input []byte) *Decoder {
	d.buffer = input
	d.offset = 0
	d.total = uint64(len(input))

	return d
}

func (d *Decoder) ReadMethodCall() (*MethodCall, error) {
	methodSignature, err := d.ReadMethod()
	if err != nil {
		return nil, err
	}

	methodDef, err := NewMethodDef(methodSignature)
	if err != nil {
		return nil, err
	}

	// Method offset of 4 since all offset jump must take into accounts the first 4 bytes of the input
	parameters, err := d.readParameters(methodDef.Parameters, 4)
	if err != nil {
		return nil, fmt.Errorf("read parameters: %w", err)
	}

	return methodDef.NewCall(parameters...), nil
}

func (d *Decoder) ReadOutput(parameters []*MethodParameter) (out []interface{}, err error) {
	return d.readParameters(parameters, 0)
}

func (d *Decoder) readParameters(parameters []*MethodParameter, methodOffset uint64) (out []interface{}, err error) {
	out = make([]interface{}, len(parameters))
	for i, param := range parameters {
		var currentOffset uint64

		isOffset := isOffsetType(param.TypeName)
		if isOffset {
			currentOffset = d.offset
			offset, err := d.read("uint256")
			if err != nil {
				return nil, fmt.Errorf("read offset for type %q (element #%d) at offset %d: %w", param.TypeName, i, d.offset, err)
			}

			jumpToOffset := offset.(*big.Int).Uint64() + methodOffset

			// The minus 32 is to ensure that offset hits a location where at least 32 bytes can be read
			if jumpToOffset > d.total-32 {
				return nil, NewErrDecoding("invalid offset value %d (max possible value %d) for type %q (element #%d) at offset %d", jumpToOffset, d.total-32, param.TypeName, i, d.offset)
			}

			d.offset = jumpToOffset
		}

		value, err := d.Read(param.TypeName)
		if err != nil {
			return nil, fmt.Errorf("read type %q (element #%d) at offset %d: %w", param.TypeName, i, d.offset, err)
		}

		if isOffset {
			d.offset = (currentOffset + 32)
		}

		out[i] = value
	}
	return
}

func (d *Decoder) Read(typeName string) (interface{}, error) {
	var isAnArray bool
	isAnArray, resolvedTypeName := isArray(typeName)
	if !isAnArray {
		return d.read(resolvedTypeName)
	}

	length, err := d.read("uint256")
	if err != nil {
		return nil, fmt.Errorf("cannot read slice %s size: %w", typeName, err)
	}

	size := length.(*big.Int).Uint64()

	arr, err := newArray(resolvedTypeName, size)
	if err != nil {
		return nil, fmt.Errorf("cannot setup new array: %w", err)
	}

	for i := uint64(0); i < size; i++ {
		out, err := d.read(resolvedTypeName)
		if err != nil {
			return nil, fmt.Errorf("cannot read item from slice %s.%d: %w", typeName, i, err)
		}
		arr.At(i, out)
	}
	return arr, nil
}

func (d *Decoder) read(typeName string) (out interface{}, err error) {
	switch typeName {
	case "bool":
		return d.ReadBool()
	case "uint8":
		v, err := d.ReadUint64()
		if err != nil {
			return nil, err
		}
		return uint8(v), nil
	case "uint16":
		v, err := d.ReadUint64()
		if err != nil {
			return nil, err
		}
		return uint16(v), nil
	case "uint24":
		v, err := d.ReadUint64()
		if err != nil {
			return nil, err
		}
		return uint32(v), nil
	case "uint32":
		v, err := d.ReadUint64()
		if err != nil {
			return nil, err
		}
		return uint32(v), nil
	case "uint40":
		v, err := d.ReadUint64()
		if err != nil {
			return nil, err
		}
		return v, nil
	case "uint48":
		v, err := d.ReadUint64()
		if err != nil {
			return nil, err
		}
		return v, nil
	case "uint56":
		v, err := d.ReadUint64()
		if err != nil {
			return nil, err
		}
		return v, nil
	case "uint64":
		v, err := d.ReadUint64()
		if err != nil {
			return nil, err
		}
		return v, nil
	case "uint72", "uint80", "uint88", "uint96", "uint104", "uint112", "uint120", "uint128", "uint136", "uint144", "uint152", "uint160", "uint168", "uint176", "uint184", "uint192", "uint200", "uint208", "uint216", "uint224", "uint232", "uint240", "uint248", "uint256":
		return d.ReadBigInt()
	case "method":
		return d.ReadMethod()
	case "address":
		return d.ReadAddress()
	case "string":
		return d.ReadString()
	case "bytes":
		return d.ReadBytes()
	}

	return nil, NewErrDecoding("type %q is not handled right now", typeName)
}

func (d *Decoder) ReadMethod() (out string, err error) {
	data, err := d.ReadBuffer(4)
	if err != nil {
		return out, err
	}
	signatureID := hex.EncodeToString(data)
	return "", NewErrDecoding("method signature not found for %s", signatureID)
}

func (d *Decoder) ReadBool() (out bool, err error) {
	data, err := d.ReadBuffer(32)
	if err != nil {
		return out, err
	}
	return (data[31] == byte(0x01)), nil
}

func (d *Decoder) ReadString() (out string, err error) {
	size, err := d.ReadBigInt()
	if err != nil {
		return out, err
	}

	remaining := 32 - (size.Uint64() % 32)
	data, err := d.ReadBuffer(size.Uint64())
	if err != nil {
		return out, err
	}

	out = strings.ToValidUTF8(string(data), "�")
	d.offset += remaining

	return
}

func (d *Decoder) ReadBytes() ([]byte, error) {
	size, err := d.ReadBigInt()
	if err != nil {
		return nil, err
	}

	data, err := d.ReadBuffer(size.Uint64())
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (d *Decoder) ReadAddress() (out Address, err error) {
	data, err := d.ReadBuffer(32)
	if err != nil {
		return out, err
	}

	address := Address(data[12:])
	zlog.Debug("read address", zap.Stringer("value", address))

	return address, nil
}

func (d *Decoder) ReadUint64() (out uint64, err error) {
	data, err := d.ReadBuffer(32)
	if err != nil {
		return out, err
	}
	return binary.BigEndian.Uint64(data[24:]), nil
}

func (d *Decoder) ReadBigInt() (out *big.Int, err error) {
	data, err := d.ReadBuffer(32)
	if err != nil {
		return out, err
	}

	return new(big.Int).SetBytes(data[:]), nil
}

func (d *Decoder) ReadBuffer(byteCount uint64) ([]byte, error) {
	zlog.Debug("trying to read bytes", zap.Uint64("byte_count", byteCount), zap.Uint64("remaining", d.total-d.offset))

	if d.total-d.offset < byteCount {
		return nil, NewErrDecoding("not enough bytes to read %d bytes, only %d remaining", byteCount, d.total-d.offset)
	}

	out := d.buffer[d.offset : d.offset+byteCount]
	zlog.Debug("read bytes", zap.Uint64("byte_count", byteCount), zap.String("data", hex.EncodeToString(out)))

	d.offset += byteCount

	return out, nil
}

type decodedArray interface {
	At(index uint64, value interface{})
}

func newArray(typeName string, count uint64) (decodedArray, error) {
	switch typeName {
	case "bool":
		return BoolArray(make([]bool, count)), nil
	case "uint8":
		return Uint8Array(make([]uint8, count)), nil
	case "uint16":
		return Uint16Array(make([]uint16, count)), nil
	case "uint24", "uint32":
		return Uint32Array(make([]uint32, count)), nil
	case "uint40", "uint48", "uint56", "uint64":
		return Uint64Array(make([]uint64, count)), nil
	case "uint72", "uint80", "uint88", "uint96", "uint104", "uint112", "uint120", "uint128", "uint136", "uint144", "uint152", "uint160", "uint168", "uint176", "uint184", "uint192", "uint200", "uint208", "uint216", "uint224", "uint232", "uint240", "uint248", "uint256":
		return BigIntArray(make([]*big.Int, count)), nil
	case "address":
		return AddressArray(make([]Address, count)), nil
	case "string":
		return StringArray(make([]string, count)), nil
	}

	return nil, NewErrDecoding("array of type %q is not handled right now", typeName)
}

type BoolArray []bool

func (a BoolArray) At(index uint64, value interface{}) {
	([]bool)(a)[index] = value.(bool)
}

type StringArray []string

func (a StringArray) At(index uint64, value interface{}) {
	([]string)(a)[index] = value.(string)
}

type AddressArray []Address

func (a AddressArray) At(index uint64, value interface{}) {
	([]Address)(a)[index] = value.(Address)
}

type BigIntArray []*big.Int

func (a BigIntArray) At(index uint64, value interface{}) {
	([]*big.Int)(a)[index] = value.(*big.Int)
}

type Uint8Array []uint8

func (a Uint8Array) At(index uint64, value interface{}) {
	([]uint8)(a)[index] = value.(uint8)
}

type Uint16Array []uint16

func (a Uint16Array) At(index uint64, value interface{}) {
	([]uint16)(a)[index] = value.(uint16)
}

type Uint32Array []uint32

func (a Uint32Array) At(index uint64, value interface{}) {
	([]uint32)(a)[index] = value.(uint32)
}

type Uint64Array []uint64

func (a Uint64Array) At(index uint64, value interface{}) {
	([]uint64)(a)[index] = value.(uint64)
}
