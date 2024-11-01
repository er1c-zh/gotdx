package api

import (
	"context"
	"gotdx/models"
	"gotdx/tdx"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func (a *App) EmitProcessInfo(i models.ProcessInfo) {
	runtime.EventsEmit(a.ctx, string(MsgKeyProcessMsg), i)
}

func (a *App) Connect(host string) string {
	if host == "" {
		host = "124.71.187.122:7709"
	}
	if a.cli == nil {
		a.cli = tdx.New(tdx.DefaultOption.WithTCPAddress(host).WithMsgCallback(a.EmitProcessInfo))
	}
	reply, err := a.cli.Connect()
	if err != nil {
		return err.Error()
	}
	runtime.EventsEmit(a.ctx, string(MsgKeyConnectionStatus), 1)
	go a.InitBasicInfo()
	return reply.Info
}

func (a *App) FetchStatus() IndexInfo {
	a.InitBasicInfo()
	return a.indexInfo
}

func (a *App) FetchStockMeta(code string) StockMeta {
	a.InitBasicInfo()
	return a.indexInfo.StockMap[code]
}
