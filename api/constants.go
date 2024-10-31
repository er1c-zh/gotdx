package api

type MsgKey string

const (
	MsgKeyProcessMsg       MsgKey = "processMsg"
	MsgKeyConnectionStatus MsgKey = "connectionStatus"
)

var ExportMsg = []struct {
	Value  MsgKey
	TSName string
}{
	{MsgKeyProcessMsg, string(MsgKeyProcessMsg)},
	{MsgKeyConnectionStatus, string(MsgKeyConnectionStatus)},
}
