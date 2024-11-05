package v2

import (
	"bytes"
	"context"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type ServerInfo struct {
	BlankCodec
	StaticCodec
	Resp ServerInfoResp
}

type ServerInfoResp struct {
	Name string
}

func (c *Client) ServerInfo() (*ServerInfo, error) {
	var err error
	s := ServerInfo{}
	s.SetDebug(c.ctx)
	s.SetContentHex(c.ctx, "01")
	err = do(c, c.dataConn, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *ServerInfo) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x000D
	header.PacketType = 1
	return nil
}

func (s *ServerInfo) MarshalReqBody(ctx context.Context) ([]byte, error) {
	return s.StaticCodec.MarshalReqBody(ctx)
}

func (s *ServerInfo) UnmarshalResp(ctx context.Context, data []byte) error {
	cursor := 68
	buf := bytes.NewBuffer(nil)
	for cursor < len(data) && data[cursor] != '\000' {
		buf.WriteByte(data[cursor])
		cursor += 1
	}
	utf8Buf, err := simplifiedchinese.GBK.NewDecoder().Bytes(buf.Bytes())
	if err != nil {
		return err
	}
	s.Resp.Name = string(utf8Buf)
	return nil
}
