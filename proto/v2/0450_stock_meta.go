package v2

import (
	"bytes"
	"context"
	"encoding/binary"
	"os"
)

func (c *Client) StockMeta(market uint16, offset uint16) (*StockMetaResp, error) {
	var err error
	s := StockMeta{}
	s.SetDebug(c.ctx)
	s.Req = &StockMetaReq{}
	s.Req.Market = market
	s.Req.Offset = offset
	err = do(c, c.dataConn, &s)
	if err != nil {
		return nil, err
	}
	return s.Resp, nil
}

type StockMeta struct {
	BlankCodec
	Req  *StockMetaReq
	Resp *StockMetaResp
}

type StockMetaReq struct {
	Market uint16
	Offset uint16
}

type StockMetaResp struct {
	Count uint16
	List  []StockMetaItem
}

type StockMetaItem struct {
	Code         string
	VolUnit      uint16
	Reserved1    uint32
	DecimalPoint int8
	Name         string
	PreClose     float64
	Reserved2    uint32
}

func (StockMeta) FillReqHeader(ctx context.Context, header *ReqHeader) error {
	header.MagicNumber = 0x0C
	header.PacketType = 1
	header.Method = 0x0450
	return nil
}

func (obj *StockMeta) MarshalReqBody(ctx context.Context) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, obj.Req)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (obj *StockMeta) UnmarshalResp(ctx context.Context, data []byte) error {

	{
		f, err := os.OpenFile("stock_meta.list", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			f.Close()
			return err
		}
		f.Write(data)
		f.Close()
	}

	var err error
	cursor := 0
	obj.Resp = &StockMetaResp{}
	obj.Resp.Count, err = ReadInt(data, &cursor, obj.Resp.Count)
	if err != nil {
		return err
	}

	for i := 0; i < int(obj.Resp.Count); i++ {
		item := StockMetaItem{}
		item.Code, err = ReadCode(data, &cursor)
		if err != nil {
			return err
		}
		item.VolUnit, err = ReadInt(data, &cursor, item.VolUnit)
		if err != nil {
			return err
		}

		item.Name, err = ReadTDXString(data, &cursor, 8)
		if err != nil {
			return err
		}

		item.Reserved1, err = ReadInt(data, &cursor, item.Reserved1)
		if err != nil {
			return err
		}
		item.DecimalPoint, err = ReadInt(data, &cursor, item.DecimalPoint)
		if err != nil {
			return err
		}
		item.PreClose, err = ReadTDXFloat(data, &cursor)
		if err != nil {
			return err
		}
		item.Reserved2, err = ReadInt(data, &cursor, item.Reserved2)
		if err != nil {
			return err
		}
		obj.Resp.List = append(obj.Resp.List, item)
	}
	return nil
}
