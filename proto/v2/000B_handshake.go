package v2

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
)

type Handshake struct {
	StaticCodec
	RespHex string
	Req     *HandshakeReq
}

type HandshakeReq struct {
	// dumb padding      offset
	// Method uint16   // 0x0
	R0         uint8    // 0x2 0x0
	LoginType  [50]byte // 0x3 guest
	R1         [50]byte // 0x35
	R2         [10]byte // 0x67
	Padding0   byte     // 0x71
	Flag0      uint8    // 0x72 8 - 1 ConnectCfgOtherWhichAutoUpInfoMinusOne
	Flag1      uint8    // 0x73 0x1
	Flag2      uint16   // 0x74 0x0104
	Flag3      uint8    // 0x76 0x1
	Padding1   byte     // 0x77
	Version    [4]byte  // 0x78 7.66 => 0xB81EF540
	MinVersion [4]byte  // 0x7C 6.04 => 0xAE47C140
	Flag4      uint8    // 0x80 0x1 from param
	Flag5      [2]uint8 // 0x81 0x0101
	Padding2   [48]byte // 0x83 padding
	SomeKey    [12]byte // 0xB3 000C29F019C2
	Padding3   [87]byte // 0xBE padding
	Flag6      uint32   // 0x116 0x000056C2
}

func (c *Client) TDXHandshake() (string, error) {
	handshake := &Handshake{}
	handshake.Req = &HandshakeReq{}
	handshake.Req.R0 = 0x0
	handshake.Req.LoginType = [50]byte{'g', 'u', 'e', 's', 't'}
	handshake.Req.Flag0 = 0x7
	handshake.Req.Flag1 = 0x1
	handshake.Req.Flag2 = 0x0104
	handshake.Req.Flag3 = 0x0
	handshake.Req.Version = [4]byte{0xB8, 0x1E, 0xF5, 0x40}
	handshake.Req.MinVersion = [4]byte{0xAE, 0x47, 0xC1, 0x40}
	handshake.Req.Flag4 = 0x1
	handshake.Req.Flag5 = [2]uint8{0x1, 0x1}
	handshake.Req.SomeKey = [12]byte{'0', '0', '0', 'C', '2', '9', 'F', '0', '1', '9', 'C', '2'}
	handshake.Req.Flag6 = 0x000056C2

	handshake.SetDebug(c.ctx)

	err := do(c, c.dataConn, handshake)
	if err != nil {
		return "", err
	}
	return handshake.RespHex, nil
}

func (h *Handshake) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x000B
	header.PacketType = 1
	return nil
}

func (h *Handshake) UnmarshalResp(ctx context.Context, data []byte) error {
	h.RespHex = hex.Dump(data)
	return nil
}

func (h *Handshake) MarshalReqBody(ctx context.Context) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, h.Req)
	return buf.Bytes(), nil
}
