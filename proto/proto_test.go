package proto

import (
	"encoding/hex"
	"testing"
)

func TestGetPrice(t *testing.T) {
	bytes, err := hex.DecodeString("d104")
	if err != nil {
		t.Error(err)
	}
	t.Log(bytes)
}

func TestGetprice(t *testing.T) {
	// bytes, err := hex.DecodeString("d104") -273
	// bytes, err := hex.DecodeString("e41b") // -1764

	// bytes, err := hex.DecodeString("a1c9cf0e") // 15331937
	// bytes, err := hex.DecodeString("90dfcf0e") // 15333328

	// bytes, err := hex.DecodeString("b3fd0f") // 130931
	// bytes, err := hex.DecodeString("e1c037") // -454689
	// bytes, err := hex.DecodeString("9bed9904") // -454689
	// bytes, err := hex.DecodeString("84cac702") // -454689

	// bytes, err := hex.DecodeString("a107") // 481
	// bytes, err := hex.DecodeString("9d63") // 6365
	// bytes, err := hex.DecodeString("50") // -16
	// bytes, err := hex.DecodeString("4a") // -10
	// bytes, err := hex.DecodeString("84cac702") // -454689
	// bytes, err := hex.DecodeString("9aad01") // 11098
	// bytes, err := hex.DecodeString("8016") // 1408
	// bytes, err := hex.DecodeString("9c8a9701") // 1237660
	// bytes, err := hex.DecodeString("ad63") // 6381
	for _, s := range []string{
		"8f0f",
		"91f901",
		"958f9705",
		"b9f0e701", "9ef0b701",
		"938209", "90dfcf0e", "96e19f03",
		"9cc709", "a808", "8b11",
		"b593cf0e", "e41b", "d104", "89a719", "92de22", "aac902", "bf6e", "41", "a1c9cf0e", "90dfcf0e",
		"e80b", "ab00", "e002", "4d", "a070", "d405"} {
		bytes, err := hex.DecodeString(s) // 6381
		if err != nil {
			return
		}
		var cursor = 0

		t.Logf("%s %d", s, ParseInt(bytes, &cursor))
	}
}

func TestGetVolume(t *testing.T) {
	t.Logf("%f", getvolume(1355337109))
	t.Logf("%f", getvolume(77))
}
