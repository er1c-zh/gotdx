package api

import (
	"context"
	"fmt"
	"gotdx/tdx"
	"sync"
)

// App struct
type App struct {
	ctx context.Context

	// tdx
	cli       *tdx.Client
	once      sync.Once
	indexInfo IndexInfo
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Shutdown(_ctx context.Context) {
	if a.cli != nil {
		a.cli.Disconnect()
	}
}

func (a *App) Connect(host string) string {
	if host == "" {
		host = "124.71.187.122:7709"
	}
	if a.cli == nil {
		a.cli = tdx.New(tdx.WithTCPAddress(host))
	}
	reply, err := a.cli.Connect()
	if err != nil {
		return err.Error()
	}
	go a.InitBasicInfo()
	return fmt.Sprintf("server: %s", reply.Info)
}

func (a *App) FetchStatus() IndexInfo {
	a.InitBasicInfo()
	return a.indexInfo
}

func (a *App) FetchStockMeta(code string) StockMeta {
	a.InitBasicInfo()
	return a.indexInfo.StockMap[code]
}
