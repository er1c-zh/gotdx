package proto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"slices"
)

type GetMinuteTimeData struct {
	reqHeader  *ReqHeader
	respHeader *RespHeader
	request    *GetMinuteTimeDataRequest
	reply      *GetMinuteTimeDataReply

	contentHex string
}

type GetMinuteTimeDataRequest struct {
	Market uint16
	Code   [6]byte
	Date   uint32
}

type BillboardRow struct {
	DeltaFromCurrentPrice int
	Volume                int
}
type GetMinuteTimeDataReply struct {
	Count uint16
	List  []MinuteTimeData

	Reserved0 string
	Reserved1 int
	Reserved2 string

	CurrentPrice  int
	VolumeAfter   int
	VolumeAll     int
	VolumeCurrent int
	TradeAmmount  float32
	VolumeBuy     int
	VolumeSell    int
	Billboard     []BillboardRow
}

type MinuteTimeData struct {
	Price     int
	Volume    int
	Reserved0 int
}

func NewGetMinuteTimeData() *GetMinuteTimeData {
	obj := new(GetMinuteTimeData)
	obj.reqHeader = new(ReqHeader)
	obj.respHeader = new(RespHeader)
	obj.request = new(GetMinuteTimeDataRequest)
	obj.reply = new(GetMinuteTimeDataReply)

	obj.reqHeader.Zip = 0x0c
	obj.reqHeader.SeqID = GenSeqID()
	obj.reqHeader.PacketType = 0x00
	//obj.reqHeader.PkgLen1  =
	//obj.reqHeader.PkgLen2  =
	obj.reqHeader.Method = 0x051d
	//obj.reqHeader.Method = KMSG_MINUTETIMEDATA
	obj.contentHex = ""
	return obj
}

func (obj *GetMinuteTimeData) SetParams(req *GetMinuteTimeDataRequest) {
	obj.request = req
}

func (obj *GetMinuteTimeData) Serialize() ([]byte, error) {
	obj.reqHeader.PkgLen1 = 0x0e
	obj.reqHeader.PkgLen2 = 0x0e

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, obj.reqHeader)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, obj.request)
	if err != nil {
		return nil, err
	}
	b, err := hex.DecodeString(obj.contentHex)
	buf.Write(b)

	//b, err := hex.DecodeString(obj.contentHex)
	//buf.Write(b)

	//err = binary.Write(buf, binary.LittleEndian, uint16(len(obj.stocks)))

	return buf.Bytes(), err
}

// 结果数据都是\n,\t分隔的中文字符串，比如查询K线数据，返回的结果字符串就形如
// /“时间\t开盘价\t收盘价\t最高价\t最低价\t成交量\t成交额\n
// /20150519\t4.644000\t4.732000\t4.747000\t4.576000\t146667487\t683638848.000000\n
// /20150520\t4.756000\t4.850000\t4.960000\t4.756000\t353161092\t1722953216.000000”
func (obj *GetMinuteTimeData) UnSerialize(header interface{}, data []byte) error {
	obj.respHeader = header.(*RespHeader)

	pos := 0
	err := binary.Read(bytes.NewBuffer(data[pos:pos+2]), binary.LittleEndian, &obj.reply.Count)
	if err != nil {
		return err
	}
	pos += 2

	pos += 3 // reserved
	pos += 6 // code

	obj.reply.Reserved0 = fmt.Sprintf("%08b", data[pos:pos+2]) // reserved
	pos += 2

	obj.reply.CurrentPrice = getprice(data, &pos) // current price
	for i := 0; i < 5; i++ {
		_ = getprice(data, &pos) // reserved
	}
	obj.reply.VolumeAfter = getprice(data, &pos) // negative open? or 盘后量
	if obj.reply.VolumeAfter < 0 {
		obj.reply.VolumeAfter = -1
	}

	obj.reply.VolumeAll = getprice(data, &pos)     // total volume
	obj.reply.VolumeCurrent = getprice(data, &pos) // current volume b3fd0f a107

	tmp := uint32(0)
	err = binary.Read(bytes.NewBuffer(data[pos:pos+4]), binary.LittleEndian, &tmp)
	if err != nil {
		return err
	}
	pos += 4
	obj.reply.TradeAmmount = float32(getvolume(int(tmp))) // reserved 9d63504a e1c03750 成交额

	obj.reply.VolumeSell = ParseInt(data, &pos) // reserved 内盘 9bed9904 9aad01
	obj.reply.VolumeBuy = ParseInt(data, &pos)  // reserved 外盘 84cac702 8016

	pos += 1 // reserved 00

	obj.reply.Reserved1 = ParseInt(data, &pos) // reserved 9c8a9701 ad63 958f9705

	// 3 档盘口
	for i := 0; i < 3; i++ {
		d1, d2 := ParseInt(data, &pos), ParseInt(data, &pos)
		v1, v2 := ParseInt(data, &pos), ParseInt(data, &pos)
		if d1 != obj.reply.CurrentPrice && d1 != -obj.reply.CurrentPrice {
			obj.reply.Billboard = append(obj.reply.Billboard, BillboardRow{d1, v1})
		}
		if d2 != obj.reply.CurrentPrice && d2 != -obj.reply.CurrentPrice {
			obj.reply.Billboard = append(obj.reply.Billboard, BillboardRow{d2, v2})
		}
	}
	slices.SortFunc(obj.reply.Billboard, func(a, b BillboardRow) int {
		return a.DeltaFromCurrentPrice - b.DeltaFromCurrentPrice
	})

	// 2byte reserved
	obj.reply.Reserved2 = fmt.Sprintf("%08b", data[pos:pos+2]) // reserved
	pos += 2

	curPrice := 0
	for index := uint16(0); index < obj.reply.Count; index++ {
		curPrice += ParseInt(data, &pos)
		r0 := getprice(data, &pos)
		vol := ParseInt(data, &pos)
		ele := MinuteTimeData{
			Price:     curPrice,
			Volume:    vol,
			Reserved0: r0,
		}
		obj.reply.List = append(obj.reply.List, ele)
	}
	return nil
}

func (obj *GetMinuteTimeData) Reply() *GetMinuteTimeDataReply {
	return obj.reply
}
