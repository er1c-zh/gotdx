package v2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func GenerateCodeBytesArray(s string) ([6]byte, error) {
	if len(s) != 6 {
		return [6]byte{}, fmt.Errorf("GenerateCodeBytesArray error")
	}
	var b [6]byte
	copy(b[:], s)
	return b, nil
}

func ReadCode(b []byte, cursor *int) (string, error) {
	if len(b) < *cursor+6 {
		return "", errors.New("read code overflow")
	}
	defer func() {
		*cursor += 6
	}()
	return string(b[*cursor : *cursor+6]), nil
}

func ReadTDXInt(b []byte, cursor *int) (int, error) {
	value := 0

	shift := 6

	cur := b[*cursor]
	*cursor += 1

	value += int(cur & 0b00111111)
	negative := false
	if (cur & 0b01000000) > 0 {
		negative = true
	}

	for (cur & 0b10000000) > 0 {
		cur = b[*cursor]
		*cursor += 1

		value += (int(cur&0b01111111) << shift)

		shift += 7
	}

	if negative {
		value = -value
	}
	return value, nil
}
func ReadTDXFloat(data []byte, cursor *int) (float64, error) {
	if len(data) < *cursor+4 {
		return 0, errors.New("read float overflow")
	}
	lleax := int(data[*cursor]) //[0]
	*cursor += 1
	lheax := int(data[*cursor]) //[1]
	*cursor += 1
	hleax := int(data[*cursor]) // [2]
	*cursor += 1
	logpoint := int(data[*cursor])
	*cursor += 1

	dwEcx := logpoint*2 - 0x7f // 0b0111 1111
	dwEdx := logpoint*2 - 0x86 // 0b1000 0110
	dwEsi := logpoint*2 - 0x8e // 0b1000 1110
	dwEax := logpoint*2 - 0x96 // 0b1001 0110
	tmpEax := dwEcx
	if dwEcx < 0 {
		tmpEax = -dwEcx
	} else {
		tmpEax = dwEcx
	}

	dbl_xmm6 := 0.0
	dbl_xmm6 = math.Pow(2.0, float64(tmpEax))
	if dwEcx < 0 {
		dbl_xmm6 = 1.0 / dbl_xmm6
	}

	dbl_xmm4 := 0.0
	dbl_xmm0 := 0.0

	if hleax > 0x80 {
		tmpdbl_xmm3 := 0.0
		dwtmpeax := dwEdx + 1
		tmpdbl_xmm3 = math.Pow(2.0, float64(dwtmpeax))
		dbl_xmm0 = math.Pow(2.0, float64(dwEdx)) * 128.0
		dbl_xmm0 += float64(hleax&0x7f) * tmpdbl_xmm3
		dbl_xmm4 = dbl_xmm0
	} else {
		if dwEdx >= 0 {
			dbl_xmm0 = math.Pow(2.0, float64(dwEdx)) * float64(hleax)
		} else {
			dbl_xmm0 = (1 / math.Pow(2.0, float64(dwEdx))) * float64(hleax)
		}
		dbl_xmm4 = dbl_xmm0
	}

	dbl_xmm3 := math.Pow(2.0, float64(dwEsi)) * float64(lheax)
	dbl_xmm1 := math.Pow(2.0, float64(dwEax)) * float64(lleax)
	if (hleax & 0x80) > 0 {
		dbl_xmm3 *= 2.0
		dbl_xmm1 *= 2.0
	}
	return dbl_xmm6 + dbl_xmm4 + dbl_xmm3 + dbl_xmm1, nil
}

func ReadInt[T any](b []byte, cursor *int, _type T) (T, error) {
	var length int = 0
	switch any(_type).(type) {
	case int8, uint8:
		length = 1
	case int16, uint16:
		length = 2
	case int32, uint32, int, uint:
		length = 4
	case int64, uint64:
		length = 8
	default:
		return _type, fmt.Errorf("unsupported type:%T", _type)
	}
	err := binary.Read(bytes.NewBuffer(b[*cursor:*cursor+length]), binary.LittleEndian, &_type)
	if err != nil {
		return _type, fmt.Errorf("ReadAsInt error:%s", err)
	}
	*cursor += length
	return _type, nil
}

func ReadByteArray(b []byte, cursor *int, length int) ([]byte, error) {
	if *cursor+length > len(b) {
		return nil, errors.New("ReadAsByteArray overflow")
	}
	defer func() {
		*cursor += length
	}()
	return b[*cursor : *cursor+length], nil
}

func ReadTDXString(b []byte, cursor *int, fixedLength int) (string, error) {
	nameGBKBuf, err := ReadByteArray(b, cursor, 8)
	if err != nil {
		return "", err
	}
	nameUtf8Buf, err := simplifiedchinese.GBK.NewDecoder().Bytes(nameGBKBuf)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(nameUtf8Buf), "\x00"), nil
}
