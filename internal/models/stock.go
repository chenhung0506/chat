package models

var StockCategorys = []StockCategory{
	Semiconductor,
	BioTechnology,
	Shipping,
	ETF,
	Financial,
}

var (
	Semiconductor = StockCategory{Code: 1, Desc: "半導體", Value: []string{"2303.TW", "2330.TW", "2337.TW", "2408.TW", "2388.TW", "2449.TW", "3450.TW", "3661.TW", "4968.TW", "6533.TW", "6770.TW", "8150.TW"}}
	BioTechnology = StockCategory{Code: 2, Desc: "生技", Value: []string{"1707.TW", "1734.TW", "1762.TW", "3164.TW", "4142.TW", "4746.TW", "6431.TW", "1733.TW"}}
	Shipping      = StockCategory{Code: 3, Desc: "航運", Value: []string{"2603.TW", "2606.TW", "2607.TW", "2608.TW", "2609.TW", "2610.TW", "2615.TW", "2618.TW"}}
	ETF           = StockCategory{Code: 4, Desc: "ETF", Value: []string{"0050.TW", "0051.TW", "0056.TW", "00660.TW", "00878.TW", "00927.TW", "00929.TW", "00946.TW"}}
	Financial     = StockCategory{Code: 5, Desc: "金融", Value: []string{"2801.TW", "2838.TW", "2881.TW", "2884.TW", "2885.TW", "2886.TW", "2891.TW", "5876.TW", "5880.TW"}}
)

type StockCategory struct {
	Code  int
	Desc  string
	Value []string
}

var StockNames = map[string]string{
	"2303.TW":  "聯電",
	"2330.TW":  "台積電",
	"2337.TW":  "旺宏",
	"2408.TW":  "南亞科",
	"2388.TW":  "威盛",
	"2449.TW":  "京元電子",
	"3450.TW":  "聯鈞",
	"3661.TW":  "世芯-KY",
	"4968.TW":  "立積",
	"6533.TW":  "晶心科",
	"6770.TW":  "力積電",
	"8150.TW":  "南茂",
	"1707.TW":  "葡萄王",
	"1734.TW":  "杏輝",
	"1762.TW":  "中化生",
	"3164.TW":  "景岳",
	"4142.TW":  "國光生",
	"4746.TW":  "台耀",
	"6431.TW":  "光麗-KY",
	"1733.TW":  "五鼎",
	"2603.TW":  "長榮",
	"2606.TW":  "裕民",
	"2607.TW":  "榮運",
	"2608.TW":  "大榮",
	"2609.TW":  "陽明",
	"2610.TW":  "華航",
	"2615.TW":  "萬海",
	"2618.TW":  "長榮航",
	"0050.TW":  "元大台灣50",
	"0051.TW":  "元大中型100",
	"0056.TW":  "元大高股息",
	"00660.TW": "元大寶滬深",
	"00878.TW": "國泰永續高股息",
	"00927.TW": "富邦台灣優質高息30",
	"00929.TW": "元大台灣ESG永續",
	"00946.TW": "兆豐藍籌30",
	"2801.TW":  "彰銀",
	"2838.TW":  "聯邦銀",
	"2881.TW":  "富邦金",
	"2884.TW":  "玉山金",
	"2885.TW":  "元大金",
	"2886.TW":  "兆豐金",
	"2891.TW":  "中信金",
	"5876.TW":  "上海商銀",
	"5880.TW":  "合庫金",
}

func GetStockCategoryDesc() []string {
	var values []string
	for _, categorys := range StockCategorys {
		values = append(values, categorys.Desc)
	}
	return values
}

func GetStockCategoryValuesByDesc(desc string) []string {
	for _, category := range StockCategorys {
		if category.Desc == desc {
			return category.Value
		}
	}
	return nil
}
