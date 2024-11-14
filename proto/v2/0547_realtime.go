package v2

import (
	"bytes"
	"context"
	"encoding/binary"
)

func (c *Client) Realtime(stock []StockQuery) error {
	realtime := &Realtime{}

	realtime.SetDebug(c.ctx)

	realtime.Req = &RealtimeReq{
		R0: 7,
		R1: 0,
	}
	for _, s := range stock {
		realtime.Req.ItemList = append(realtime.Req.ItemList, RealtimeReqItem{
			Market: s.Market,
			Code:   [6]byte{s.Code[0], s.Code[1], s.Code[2], s.Code[3], s.Code[4], s.Code[5]},
			// R0:     [4]uint8{0xE2, 0x2B, 0x02, 0x00}, // 0xE22B0200
			R0: [4]uint8{0x00, 0x00, 0x00, 0x00},
		})
	}

	err := do(c, c.dataConn, realtime)
	if err != nil {
		return err
	}
	return nil
}

type Realtime struct {
	BlankCodec
	Req *RealtimeReq
}

type RealtimeReq struct {
	R0       uint8
	R1       uint8
	ItemList []RealtimeReqItem
}

type RealtimeReqItem struct {
	Market uint8
	Code   [6]byte
	R0     [4]uint8
}

func (obj *Realtime) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x0547
	return nil
}

func (obj *Realtime) MarshalReqBody(ctx context.Context) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(nil)
	err = binary.Write(buf, binary.LittleEndian, obj.Req.R0)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, obj.Req.R1)
	if err != nil {
		return nil, err
	}
	for _, item := range obj.Req.ItemList {
		err = binary.Write(buf, binary.LittleEndian, item)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (obj *Realtime) UnmarshalResp(ctx context.Context, data []byte) error {
	return nil
}
