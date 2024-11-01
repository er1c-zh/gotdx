package proto

import (
	"bytes"
	"encoding/binary"
)

type Heartbeat struct {
	ReqHeader  *ReqHeader
	RespHeader *RespHeader
}

func NewHeartbeat() *Heartbeat {
	obj := new(Heartbeat)
	obj.ReqHeader = new(ReqHeader)
	obj.ReqHeader.Zip = 0x0c
	obj.ReqHeader.SeqID = seqID()
	obj.ReqHeader.PacketType = 0x02
	obj.ReqHeader.PkgLen1 = 0x02
	obj.ReqHeader.PkgLen2 = 0x02
	obj.ReqHeader.Method = KMSG_HEARTBEAT

	return obj
}

type HeartbeatResp struct {
	Reserved [6]byte
	Date     int32
}

func (obj *Heartbeat) Serialize() ([]byte, error) {
	var err error
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.LittleEndian, obj.ReqHeader)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (obj *Heartbeat) UnSerialize(header interface{}, data []byte) error {
	obj.RespHeader = header.(*RespHeader)
	return nil
}
