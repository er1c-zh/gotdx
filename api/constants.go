package api

import (
	"gotdx/models"
)

type MsgKey string

const (
	MsgKeyInit             MsgKey = "init"
	MsgKeyProcessMsg       MsgKey = "processMsg"
	MsgKeyConnectionStatus MsgKey = "connectionStatus"
)

var ExportMsg = []struct {
	Value  MsgKey
	TSName string
}{
	{MsgKeyInit, string(MsgKeyInit)},
	{MsgKeyProcessMsg, string(MsgKeyProcessMsg)},
	{MsgKeyConnectionStatus, string(MsgKeyConnectionStatus)},
}

var ExportMarketType = []struct {
	Value  models.MarketType
	TSName string
}{
	{models.MarketSZ, models.MarketSZ.String()},
	{models.MarketSH, models.MarketSH.String()},
	{models.MarketBJ, models.MarketBJ.String()},
}
