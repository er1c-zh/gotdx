package proto

import (
	"bytes"
	"encoding/binary"
)

type GetHistoryMinuteTimeData struct {
	reqHeader  *ReqHeader
	respHeader *RespHeader
	request    *GetHistoryMinuteTimeDataRequest
	reply      *GetHistoryMinuteTimeDataReply

	contentHex string
}

type GetHistoryMinuteTimeDataRequest struct {
	Date   uint32
	Market uint8
	Code   [6]byte
}

type GetHistoryMinuteTimeDataReply struct {
	Count     uint16
	PriceUnit int
	List      []HistoryMinuteTimeData
}

type HistoryMinuteTimeData struct {
	Price int
	Vol   int
}

func NewGetHistoryMinuteTimeData() *GetHistoryMinuteTimeData {
	obj := new(GetHistoryMinuteTimeData)
	obj.reqHeader = new(ReqHeader)
	obj.respHeader = new(RespHeader)
	obj.request = new(GetHistoryMinuteTimeDataRequest)
	obj.reply = new(GetHistoryMinuteTimeDataReply)

	obj.reqHeader.Zip = 0x0c
	obj.reqHeader.SeqID = GenSeqID()
	obj.reqHeader.PacketType = 0x00
	obj.reqHeader.Method = KMSG_HISTORYMINUTETIMEDATE
	obj.contentHex = ""
	return obj
}

func (obj *GetHistoryMinuteTimeData) SetParams(req *GetHistoryMinuteTimeDataRequest) {
	obj.request = req
}

func (obj *GetHistoryMinuteTimeData) Serialize() ([]byte, error) {
	obj.reqHeader.PkgLen1 = 0x0d
	obj.reqHeader.PkgLen2 = 0x0d

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, obj.reqHeader)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buf, binary.LittleEndian, obj.request)
	if err != nil {
		return nil, err
	}
	// b, err := hex.DecodeString(obj.contentHex)
	// buf.Write(b)

	return buf.Bytes(), err
}

// 结果数据都是\n,\t分隔的中文字符串，比如查询K线数据，返回的结果字符串就形如
// /“时间\t开盘价\t收盘价\t最高价\t最低价\t成交量\t成交额\n
// /20150519\t4.644000\t4.732000\t4.747000\t4.576000\t146667487\t683638848.000000\n
// /20150520\t4.756000\t4.850000\t4.960000\t4.756000\t353161092\t1722953216.000000”
func (obj *GetHistoryMinuteTimeData) UnSerialize(header interface{}, data []byte) error {
	obj.respHeader = header.(*RespHeader)

	pos := 0
	err := binary.Read(bytes.NewBuffer(data[pos:pos+2]), binary.LittleEndian, &obj.reply.Count)
	if err != nil {
		return err
	}
	pos += 2
	// 跳过4个字节 功能未解析
	_, _, _, _ = data[pos], data[pos+1], data[pos+2], data[pos+3]
	pos += 4

	obj.reply.PriceUnit = 100

	curPrice := 0
	for index := uint16(0); index < obj.reply.Count; index++ {
		curPrice += ParseInt(data, &pos)
		_ = ParseInt(data, &pos)
		vol := ParseInt(data, &pos)

		ele := HistoryMinuteTimeData{Price: curPrice, Vol: vol}
		obj.reply.List = append(obj.reply.List, ele)
	}
	return err
}

func (obj *GetHistoryMinuteTimeData) Reply() *GetHistoryMinuteTimeDataReply {
	return obj.reply
}
