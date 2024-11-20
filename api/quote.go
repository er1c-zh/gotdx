package api

import (
	"fmt"
	"gotdx/models"
	v2 "gotdx/proto/v2"
)

func (a *App) CandleStick(code string, period v2.CandleStickPeriodType, cursor uint16) *v2.CandleStickResp {
	meta, ok := a.stockMetaMap[code]
	if !ok {
		return nil
	}
	resp, err := a.cli.CandleStick(meta.Market, code, period, cursor)
	if err != nil {
		a.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("candle stick failed: %s", err.Error())})
		return nil
	}
	return resp
}
