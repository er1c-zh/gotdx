package tdx

import "errors"

const (
	MarketSz = 0 // 深圳
	MarketSh = 1 // 上海
	MarketBj = 2 // 北京
)

var (
	MarketStrMap = map[uint8]string{
		MarketSz: "SZ",
		MarketSh: "SH",
		MarketBj: "BJ",
	}
)

const (
	KLINE_TYPE_5MIN      = 0  // 5 分钟K 线
	KLINE_TYPE_15MIN     = 1  // 15 分钟K 线
	KLINE_TYPE_30MIN     = 2  // 30 分钟K 线
	KLINE_TYPE_1HOUR     = 3  // 1 小时K 线
	KLINE_TYPE_DAILY     = 4  // 日K 线
	KLINE_TYPE_WEEKLY    = 5  // 周K 线
	KLINE_TYPE_MONTHLY   = 6  // 月K 线
	KLINE_TYPE_EXHQ_1MIN = 7  //  1 分钟
	KLINE_TYPE_1MIN      = 8  // 1 分钟K 线
	KLINE_TYPE_RI_K      = 9  // 日K 线
	KLINE_TYPE_3MONTH    = 10 // 季K 线
	KLINE_TYPE_YEARLY    = 11 // 年K 线
)

var (
	ErrBadData = errors.New("more than 8M data")
)
