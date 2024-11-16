package v2

import (
	"bytes"
	"context"
	"encoding/binary"
)

func (c *Client) Realtime(stock []StockQuery) (*RealtimeResp, error) {
	realtime := &Realtime{}

	realtime.SetDebug(c.ctx)

	realtime.Req = &RealtimeReq{
		Size: uint16(len(stock)),
	}
	for _, s := range stock {
		realtime.Req.ItemList = append(realtime.Req.ItemList, RealtimeReqItem{
			Market: s.Market,
			Code:   [10]byte{s.Code[0], s.Code[1], s.Code[2], s.Code[3], s.Code[4], s.Code[5]},
		})
	}

	err := do(c, c.dataConn, realtime)
	if err != nil {
		return nil, err
	}

	return realtime.Resp, nil
}

type Realtime struct {
	BlankCodec
	Req  *RealtimeReq
	Resp *RealtimeResp
}

type RealtimeReq struct {
	Size     uint16
	ItemList []RealtimeReqItem
}

type RealtimeReqItem struct {
	Market uint8
	Code   [10]byte
}

type RealtimeResp struct {
	Data []byte
}

func (obj *Realtime) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.Method = 0x0547
	return nil
}

func (obj *Realtime) MarshalReqBody(ctx context.Context) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(nil)
	err = binary.Write(buf, binary.LittleEndian, obj.Req.Size)
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
	data = decryptSimpleXOR(data, 0x93)
	obj.Resp = &RealtimeResp{
		Data: data,
	}
	return nil
}
