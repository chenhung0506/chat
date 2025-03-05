package service

import (
	"chat/internal/models"
	"math/rand"
	"sort"
	"time"
)

// Stock 表示股票的基本數據
type Stock struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	PrevClose float64 `json:"prev_close"`
}

// Gainer 表示漲幅最高的股票
type Gainer struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Change    float64 `json:"change"`
	ChangePct float64 `json:"change_pct"`
}

// 獲取股票數據
func getStockData(category string) []Stock {
	rand.Seed(time.Now().UnixNano())
	stockSymbols := models.GetStockCategoryValuesByDesc(category)
	var stocks []Stock
	for _, symbol := range stockSymbols {
		stocks = append(stocks, Stock{
			Symbol:    symbol,
			Price:     rand.Float64()*200 + 100,
			PrevClose: rand.Float64()*200 + 100,
		})
	}
	return stocks
}

func getTopGainer(stocks []Stock) Gainer {
	var gainers []Gainer
	for _, stock := range stocks {
		change := stock.Price - stock.PrevClose
		changePct := (change / stock.PrevClose) * 100
		gainers = append(gainers, Gainer{
			Symbol:    stock.Symbol,
			Price:     stock.Price,
			Change:    change,
			ChangePct: changePct,
		})
	}

	sort.Slice(gainers, func(i, j int) bool {
		return gainers[i].ChangePct > gainers[j].ChangePct
	})

	return gainers[0]
}

// API 處理函數
func TopGainerHandler(desc string) Gainer {
	stocks := getStockData(desc)
	topGainer := getTopGainer(stocks)
	return topGainer
}
