package api

import (
	"fmt"
	"gotdx/models"
	v2 "gotdx/proto/v2"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type QuoteSubscripition struct {
	app    *App
	rwMu   sync.RWMutex
	cli    *v2.Client
	m      map[string]*SubscribeReq
	ticker *time.Ticker
}

type SubscribeReq struct {
	Group     string
	Code      []string
	QuoteType string
}

func NewQuoteSubscripition(app *App) *QuoteSubscripition {
	return &QuoteSubscripition{
		app:    app,
		cli:    app.cli,
		m:      make(map[string]*SubscribeReq),
		ticker: time.NewTicker(time.Second * 3),
	}
}

func (a *QuoteSubscripition) Subscribe(req *SubscribeReq) {
	a.rwMu.Lock()
	defer a.rwMu.Unlock()
	old, ok := a.m[req.Group]
	if ok {
		old.Code = append(old.Code, req.Code...)
	} else {
		a.m[req.Group] = req
	}
}

func (a *QuoteSubscripition) Unsubscribe(req *SubscribeReq) {
	a.rwMu.Lock()
	defer a.rwMu.Unlock()
	old, ok := a.m[req.Group]
	if !ok {
		return
	}
	if len(req.Code) == 0 {
		delete(a.m, req.Group)
		return
	}
	b := make([]string, len(old.Code)-1)
	for _, code := range old.Code {
		for _, c := range req.Code {
			if code == c {
				goto next
			}
		}
		b = append(b, code)
	next:
	}
	old.Code = b
}

func (a *QuoteSubscripition) Start() {
	go a.worker()
}

func (a *QuoteSubscripition) worker() {
	for {
		<-a.ticker.C
		a.rwMu.RLock()
		for _, req := range a.m {
			realtimeReq := make([]v2.StockQuery, 0, len(req.Code))
			for _, code := range req.Code {
				if a.app.stockMetaMap[code] == nil {
					continue
				}
				realtimeReq = append(realtimeReq,
					v2.StockQuery{Market: uint8(a.app.stockMetaMap[code].Market), Code: code})
			}
			go func(req []v2.StockQuery, group string) {
				resp, err := a.cli.Realtime(req)
				if err != nil {
					a.app.LogProcessError(models.ProcessInfo{Msg: fmt.Sprintf("realtime subscribe failed: %s", err.Error())})
					return
				}
				runtime.EventsEmit(a.app.ctx, string(MsgKeySubscribeBroadcast), group, resp.ItemList)
			}(realtimeReq, req.Group)
		}

		a.rwMu.RUnlock()
	}
}

func (a *App) Subscribe(req *SubscribeReq) {
	a.qs.Subscribe(req)
}

func (a *App) Unsubscribe(req *SubscribeReq) {
	a.qs.Unsubscribe(req)
}
