package api

import (
	"fmt"

	"gotdx/proto"
)

type RealtimeData struct {
	Data []proto.MinuteTimeData
	Meta proto.Security
}

func (a *App) FetchRealtimeData(code string) RealtimeData {
	if a.indexInfo.StockMap == nil {
		a.EmitProcessInfo(ProcessInfo{Msg: "Stock list not found"})
		return RealtimeData{}
	}
	meta, ok := a.indexInfo.StockMap[code]
	if !ok {
		a.EmitProcessInfo(ProcessInfo{Msg: fmt.Sprintf("Stock %s not found", code)})
		return RealtimeData{}
	}
	reply, err := a.cli.GetMinuteTimeData(meta.Market, meta.Code)
	if err != nil {
		a.EmitProcessInfo(ProcessInfo{Msg: fmt.Sprintf("Stock %s not found, %s", code, err.Error())})
		return RealtimeData{}
	}
	resp := RealtimeData{
		Data: reply.List,
		Meta: meta.Meta,
	}
	a.EmitProcessInfo(ProcessInfo{Msg: fmt.Sprintf("Get realtime data: %s %d", code, len(resp.Data))})
	return resp
}
