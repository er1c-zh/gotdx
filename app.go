package main

import (
	"context"
	"gotdx/tdx"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context

	// tdx
	cli       *tdx.Client
	mu        sync.Mutex
	indexInfo IndexInfo
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(_ctx context.Context) {
	if a.cli != nil {
		a.cli.Disconnect()
	}
}

type IndexInfo struct {
	isInit bool

	IsConnected bool
	Msg         string

	// stock
	StockCount int
	AllStock   []StockMarketList
}

type StockMarketList struct {
	Market    uint8
	MarketStr string
	Count     int
	StockList []StockMeta
	StockMap  map[string]StockMeta
}

type StockMeta struct {
	Market uint8
	Code   string
	Desc   string
}

func (a *App) Connect(host string) string {
	if host == "" {
		host = "124.71.187.122:7709"
	}
	if a.cli == nil {
		a.cli = tdx.New(tdx.WithTCPAddress(host))
	}
	_, err := a.cli.Connect()
	if err != nil {
		return err.Error()
	}
	return "success"
}

func (a *App) FetchStatus() IndexInfo {
	a.getBasicInfo()
	return a.indexInfo
}

func (a *App) getBasicInfo() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.indexInfo.isInit {
		return
	}
	runtime.LogInfof(a.ctx, "get basic info")
	for _, market := range []uint8{tdx.MarketSh, tdx.MarketSz, tdx.MarketBj} {
		runtime.LogInfof(a.ctx, "market: %d", market)
		marketItem := StockMarketList{
			Market:    market,
			MarketStr: tdx.MarketStrMap[market],
			Count:     -1,
			StockList: []StockMeta{},
			StockMap:  map[string]StockMeta{},
		}
		countResp, err := a.cli.GetSecurityCount(uint16(market))
		if err != nil {
			runtime.LogErrorf(a.ctx, "get security count error: %s", err.Error())
			continue
		}
		if countResp == nil {
			continue
		}

		runtime.LogInfof(a.ctx, "count: %d", countResp.Count)

		marketItem.Count = int(countResp.Count)
		marketItem.StockList = make([]StockMeta, marketItem.Count)
		marketItem.StockMap = make(map[string]StockMeta, marketItem.Count)
		listResp, err := a.cli.GetSecurityList(market, 0)
		if err != nil {
			runtime.LogErrorf(a.ctx, "get security list error: %s", err.Error())
			continue
		}
		if listResp == nil {
			continue
		}
		runtime.LogInfof(a.ctx, "list count: %d", len(listResp.List))
		for i := 0; i < marketItem.Count && i < len(listResp.List); i++ {
			meta := listResp.List[i]
			stockMeta := StockMeta{
				Market: market,
				Code:   meta.Code,
				Desc:   meta.Name,
			}
			marketItem.StockList[i] = stockMeta
			marketItem.StockMap[stockMeta.Code] = stockMeta
		}

		a.indexInfo.AllStock = append(a.indexInfo.AllStock, marketItem)
	}
	a.indexInfo.isInit = true
}
