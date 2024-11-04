package v2

import (
	"context"
	"encoding/hex"
)

type Handshake struct {
	StaticCodec
	RespHex string
}

func (h *Handshake) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x000B
	header.PacketType = 1
	return nil
}

func (h *Handshake) UnmarshalResp(ctx context.Context, data []byte) error {
	h.RespHex = hex.EncodeToString(data)
	return nil
}

func (h *Handshake) MarshalReqBody(ctx context.Context) ([]byte, error) {
	h.SetContentHex(ctx, "e53878ee8bd8dbb8749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae2770035705c1ef813f7654ef62479faec20aca5b1c271aa7637d7a43749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357f881d622360d15a8fb5f8c6b71caf911749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357749933ae27700357a6d54b1b30a9535f")
	return h.StaticCodec.MarshalReqBody(ctx)
}
