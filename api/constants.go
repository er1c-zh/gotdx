package api

import (
	"gotdx/models"

	v2 "gotdx/proto/v2"
)

type MsgKey string

const (
	MsgKeyInit               MsgKey = "init"
	MsgKeyProcessMsg         MsgKey = "processMsg"
	MsgKeyServerStatus       MsgKey = "serverStatus"
	MsgKeySubscribeBroadcast MsgKey = "subscribeBroadcast"
)

var ExportMsg = []struct {
	Value  MsgKey
	TSName string
}{
	{MsgKeyInit, string(MsgKeyInit)},
	{MsgKeyProcessMsg, string(MsgKeyProcessMsg)},
	{MsgKeyServerStatus, string(MsgKeyServerStatus)},
	{MsgKeySubscribeBroadcast, string(MsgKeySubscribeBroadcast)},
}

var ExportMarketType = []struct {
	Value  models.MarketType
	TSName string
}{
	{models.MarketSZ, models.MarketSZ.String()},
	{models.MarketSH, models.MarketSH.String()},
	{models.MarketBJ, models.MarketBJ.String()},
}

var ExportCandleStickPeriodType = []struct {
	Value  v2.CandleStickPeriodType
	TSName string
}{
	{v2.CandleStickPeriodType_1Min, "CandleStickPeriodType1Min"},
	{v2.CandleStickPeriodType_5Min, "CandleStickPeriodType5Min"},
	{v2.CandleStickPeriodType_15Min, "CandleStickPeriodType15Min"},
	{v2.CandleStickPeriodType_30Min, "CandleStickPeriodType30Min"},
	{v2.CandleStickPeriodType_1Hour, "CandleStickPeriodType60Min"},
	{v2.CandleStickPeriodType_Day, "CandleStickPeriodType1Day"},
	{v2.CandleStickPeriodType_Week, "CandleStickPeriodType1Week"},
	{v2.CandleStickPeriodType_Month, "CandleStickPeriodType1Month"},
}
