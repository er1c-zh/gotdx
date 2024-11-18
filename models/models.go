package models

type ProcessInfoType string

const (
	ProcessInfoTypeMsg ProcessInfoType = "msg"
	ProcessInfoTypeErr ProcessInfoType = "err"
)

type ProcessInfo struct {
	Type ProcessInfoType
	Msg  string
}

type MarketType uint16

const (
	MarketSZ MarketType = 0
	MarketSH MarketType = 1
	MarketBJ MarketType = 2
)

func (m MarketType) String() string {
	switch m {
	case MarketSZ:
		return "深圳"
	case MarketSH:
		return "上海"
	case MarketBJ:
		return "北京"
	default:
		return "unknown"
	}
}

type StockMetaAll struct {
	StockList []StockMetaItem
}

type StockMetaItem struct {
	Market        MarketType
	Code          string
	Desc          string
	PinYinInitial string
}
