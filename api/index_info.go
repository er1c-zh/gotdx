package api

import (
	"gotdx/proto"
	"gotdx/tdx"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type IndexInfo struct {
	IsConnected bool
	Msg         string

	// stock
	StockCount int
	StockList  []StockMeta
	StockMap   map[string]StockMeta
}

type StockMeta struct {
	Market uint8
	Code   string
	Desc   string
	Meta   proto.Security
}

func (a *App) InitBasicInfo() {
	a.once.Do(a.initBasicInfo)
}
func (a *App) initBasicInfo() {
	runtime.LogInfof(a.ctx, "get basic info")
	a.indexInfo.StockCount = 0
	for _, market := range []uint8{tdx.MarketSh, tdx.MarketSz} {
		runtime.LogInfof(a.ctx, "market: %d", market)
		countResp, err := a.cli.GetSecurityCount(uint16(market))
		if err != nil {
			runtime.LogErrorf(a.ctx, "get security count error: %s", err.Error())
			continue
		}
		if countResp == nil {
			continue
		}

		runtime.LogInfof(a.ctx, "count: %d", countResp.Count)

		a.indexInfo.StockCount += int(countResp.Count)
	}

	a.indexInfo.StockList = make([]StockMeta, 0, a.indexInfo.StockCount)
	a.indexInfo.StockMap = make(map[string]StockMeta, a.indexInfo.StockCount)

	runtime.LogInfof(a.ctx, "get stock list")

	for _, market := range []uint8{tdx.MarketSh, tdx.MarketSz} {
		cursor := 0
		for {
			runtime.LogInfof(a.ctx, "market: %d, cursor: %d", market, cursor)
			listResp, err := a.cli.GetSecurityList(market, uint16(cursor))
			if err != nil {
				runtime.LogErrorf(a.ctx, "get security list error: %s", err.Error())
				break
			}
			if listResp == nil {
				break
			}
			for _, meta := range listResp.List {
				stockMeta := StockMeta{
					Market: market,
					Code:   meta.Code,
					Desc:   meta.Name,
					Meta:   meta,
				}
				a.indexInfo.StockList = append(a.indexInfo.StockList, stockMeta)
				a.indexInfo.StockMap[meta.Code] = stockMeta
			}
			if len(listResp.List) < 1000 {
				break
			}
			cursor += 1000
		}
	}
	runtime.LogInfof(a.ctx, "stock count: %d", len(a.indexInfo.StockList))
}
