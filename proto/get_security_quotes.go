package proto

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

type GetSecurityQuotes struct {
	reqHeader  *ReqHeader
	respHeader *RespHeader
	request    *GetSecurityQuotesRequest
	reply      *GetSecurityQuotesReply

	contentHex string
}

type Stock struct {
	Market uint8
	Code   string
}

type GetSecurityQuotesRequest struct {
	StockList []Stock
}

type GetSecurityQuotesReply struct {
	Count uint16
	List  []SecurityQuote
}

type SecurityQuote struct {
	Market uint8  // 市场
	Code   string // 代码
	// Active1        uint16  // 活跃度
	Price     int // 现价
	LastClose int // 昨收
	Open      int // 开盘
	High      int // 最高
	Low       int // 最低
	// ServerTime     string  // 时间
	ReversedInt0 int     // 保留(时间 ServerTime)
	ReversedInt1 int     // 保留
	Vol          int     // 总量
	CurVol       int     // 现量
	Amount       float64 // 总金额
	SVol         int     // 内盘
	BVol         int     // 外盘
	ReversedInt2 int     // 保留
	ReversedInt3 int     // 保留
	// ReversedBytes2 int     // 保留
	// ReversedBytes3 int     // 保留
	BidLevels []Level
	AskLevels []Level
	Reserved1 string
	// ReversedBytes4 uint16  // 保留
	ReversedInt5 int   // 保留
	ReversedInt6 int   // 保留
	ReversedInt7 int   // 保留
	ReversedInt8 int   // 保留
	DeltaRate    int16 // 涨速 in percent
	// Active2        uint16  // 活跃度
}

type Level struct {
	Price int
	Vol   int
}

func NewGetSecurityQuotes() *GetSecurityQuotes {
	obj := new(GetSecurityQuotes)
	obj.reqHeader = new(ReqHeader)
	obj.respHeader = new(RespHeader)
	obj.request = new(GetSecurityQuotesRequest)
	obj.reply = new(GetSecurityQuotesReply)

	obj.reqHeader.Zip = 0x0c
	obj.reqHeader.SeqID = GenSeqID()
	obj.reqHeader.PacketType = 0x01
	obj.reqHeader.Method = KMSG_SECURITYQUOTES
	obj.contentHex = "0500000000000000"
	return obj
}

func (obj *GetSecurityQuotes) SetParams(req *GetSecurityQuotesRequest) {
	obj.request = req
}

func (obj *GetSecurityQuotes) Serialize() ([]byte, error) {
	obj.reqHeader.PkgLen1 = 2 + uint16(len(obj.request.StockList)*7) + 10
	obj.reqHeader.PkgLen2 = 2 + uint16(len(obj.request.StockList)*7) + 10

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, obj.reqHeader)
	if err != nil {
		return nil, err
	}
	b, err := hex.DecodeString(obj.contentHex)
	if err != nil {
		return nil, err
	}
	buf.Write(b)

	err = binary.Write(buf, binary.LittleEndian, uint16(len(obj.request.StockList)))

	for _, stock := range obj.request.StockList {
		code := make([]byte, 6)
		copy(code, stock.Code)
		tmp := []byte{stock.Market}
		tmp = append(tmp, code...)
		buf.Write(tmp)
	}

	return buf.Bytes(), err
}

func (obj *GetSecurityQuotes) UnSerialize(header interface{}, data []byte) error {
	obj.respHeader = header.(*RespHeader)
	var err error

	pos := 0

	pos += 2 // 跳过两个字节

	obj.reply.Count, err = ReadAsInt(data, &pos, obj.reply.Count)
	if err != nil {
		return err
	}
	// binary.Read(bytes.NewBuffer(data[pos:pos+2]), binary.LittleEndian, &obj.reply.Count)
	// pos += 2

	for index := 0; index < int(obj.reply.Count); index++ {
		ele := SecurityQuote{}

		ele.Market, err = ReadAsInt(data, &pos, ele.Market)
		if err != nil {
			return err
		}
		// binary.Read(bytes.NewBuffer(data[pos:pos+1]), binary.LittleEndian, &ele.Market)
		// pos += 1

		var code [6]byte
		binary.Read(bytes.NewBuffer(data[pos:pos+6]), binary.LittleEndian, &code)
		ele.Code = Utf8ToGbk(code[:])
		pos += 6

		// unknown
		magicNumber := 0
		binary.Read(bytes.NewBuffer(data[pos:pos+2]), binary.LittleEndian, &magicNumber)
		pos += 2

		price := ParseInt(data, &pos)
		ele.Price = obj.getPrice(price, 0)
		ele.LastClose = obj.getPrice(price, ParseInt(data, &pos))
		ele.Open = obj.getPrice(price, ParseInt(data, &pos))
		ele.High = obj.getPrice(price, ParseInt(data, &pos))
		ele.Low = obj.getPrice(price, ParseInt(data, &pos))

		ele.ReversedInt0 = ParseInt(data, &pos)
		// ele.ServerTime = fmt.Sprintf("%d", ele.ReversedBytes0)
		ele.ReversedInt1 = ParseInt(data, &pos)

		ele.Vol = ParseInt(data, &pos)
		ele.CurVol = ParseInt(data, &pos)

		// var amountraw int32
		// binary.Read(bytes.NewBuffer(data[pos:pos+4]), binary.LittleEndian, &amountraw)
		// pos += 4
		// ele.Amount = ParseFloat(amountraw)
		ele.Amount, err = ReadFloat(data, &pos)
		if err != nil {
			return err
		}

		ele.SVol = ParseInt(data, &pos)
		ele.BVol = ParseInt(data, &pos)

		ele.ReversedInt2 = ParseInt(data, &pos)
		ele.ReversedInt3 = ParseInt(data, &pos)

		for i := 0; i < 5; i++ {
			bidele := Level{Price: obj.getPrice(ParseInt(data, &pos), price)}
			offerele := Level{Price: obj.getPrice(ParseInt(data, &pos), price)}
			bidele.Vol = ParseInt(data, &pos)
			offerele.Vol = ParseInt(data, &pos)
			ele.BidLevels = append(ele.BidLevels, bidele)
			ele.AskLevels = append(ele.AskLevels, offerele)
		}

		ele.Reserved1, err = ReadByteToHex(data, &pos, 2)
		if err != nil {
			return err
		}
		// binary.Read(bytes.NewBuffer(data[pos:pos+2]), binary.LittleEndian, &ele.ReversedBytes4)
		// pos += 2

		ele.ReversedInt5 = ParseInt(data, &pos)
		ele.ReversedInt6 = ParseInt(data, &pos)
		ele.ReversedInt7 = ParseInt(data, &pos)
		ele.ReversedInt8 = ParseInt(data, &pos)

		// var reversedbytes9 int16
		// binary.Read(bytes.NewBuffer(data[pos:pos+2]), binary.LittleEndian, &reversedbytes9)
		// pos += 2
		ele.DeltaRate, err = ReadAsInt(data, &pos, ele.DeltaRate)
		if err != nil {
			return err
		}

		magicNumberEnd := 0
		binary.Read(bytes.NewBuffer(data[pos:pos+2]), binary.LittleEndian, &magicNumberEnd)
		pos += 2

		if magicNumberEnd != magicNumber {
			return errors.New("magic number miss match")
		}

		obj.reply.List = append(obj.reply.List, ele)
	}
	return nil
}

func (obj *GetSecurityQuotes) Reply() *GetSecurityQuotesReply {
	return obj.reply
}

func (obj *GetSecurityQuotes) getPrice(price int, diff int) int {
	return price + diff
}
