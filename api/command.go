package api

import (
	"gotdx/models"
	"strings"
)

func (c *App) CommandMatch(s string) []models.StockMetaItem {
	result := []models.StockMetaItem{}
	if c.stockMeta == nil {
		return result
	}
	for _, item := range c.stockMeta.StockList {
		if strings.HasPrefix(item.Code, s) {
			result = append(result, item)
		} else if strings.HasPrefix(item.PinYinInitial, s) {
			result = append(result, item)
		}
		if len(result) > 5 {
			break
		}
	}
	return result
}
