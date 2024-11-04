package v2

import (
	"bytes"
	"testing"
)

func TestGenerateCodeBytesArray(t *testing.T) {
	b, err := GenerateCodeBytesArray("600000")
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(b[:], []byte{'6', '0', '0', '0', '0', '0'}) {
		t.Error("GenerateCodeBytesArray error")
	}
}
