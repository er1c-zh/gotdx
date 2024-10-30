package main

import (
	"context"
	"gotdx/tdx"

	"github.com/er1c-zh/gotdx"
)

// App struct
type App struct {
	ctx context.Context

	// tdx
	cli *gotdx.Client
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
	IsConnected bool
	Msg         string

	// stock
	StockCount int
}

func (a *App) FetchStatus() IndexInfo {
	info := IndexInfo{
		IsConnected: a.cli != nil && a.cli.IsConnected(),
		Msg:         "",
	}
	if a.cli == nil || !a.cli.IsConnected() {
		info.Msg = "未连接"
		return info
	}
	resp, err := a.cli.GetSecurityCount(tdx.MarketSh)
	if err != nil {
		info.Msg = err.Error()
		return info
	}
	if resp != nil {
		info.StockCount = int(resp.Count)
	} else {
		info.Msg = "未知错误"
		return info
	}
	return info
}

func (a *App) Connect(host string) string {
	if host == "" {
		host = "124.71.187.122:7709"
	}
	if a.cli == nil {
		a.cli = gotdx.New(gotdx.WithTCPAddress(host))
	}
	_, err := a.cli.Connect()
	if err != nil {
		return err.Error()
	}
	return "success"
}
