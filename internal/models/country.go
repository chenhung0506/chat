package models

import (
	"strings"
)

// Country 定義結構體
type Country struct {
	Code    string
	Desc    string
	Value   string
	Station string
}

// 定義所有 Country
var Countries = []Country{
	{Code: "63", Desc: "臺北市", Value: "台北", Station: "中山區"},
	{Code: "68", Desc: "桃園市", Value: "桃園", Station: "桃園區"},
	{Code: "10018", Desc: "新竹市", Value: "新竹", Station: "竹北市"},
	{Code: "66", Desc: "臺中市", Value: "台中", Station: "西區"},
	{Code: "10002", Desc: "宜蘭縣", Value: "宜蘭", Station: "宜蘭市"},
	{Code: "67", Desc: "臺南市", Value: "台南", Station: "東區"},
}

// FromCode 根據代碼查找 Country
func FromCode(code string) (Country, bool) {
	for _, country := range Countries {
		if country.Code == code {
			return country, true
		}
	}
	return Country{}, false
}

// FromDesc 根據描述查找 Country
func FromDesc(desc string) (Country, bool) {
	for _, country := range Countries {
		if strings.EqualFold(country.Desc, desc) {
			return country, true
		}
	}
	return Country{}, false
}

// FromValue 根據值查找 Country
func FromValue(value string) Country {
	for _, country := range Countries {
		if strings.EqualFold(country.Value, value) {
			return country
		}
	}
	return Country{}
}

// FromStation 根據 station 查找 Country
func FromStation(station string) (Country, bool) {
	for _, country := range Countries {
		if strings.EqualFold(country.Station, station) {
			return country, true
		}
	}
	return Country{}, false
}

func GetCountryValue() []string {
	var values []string
	for _, country := range Countries {
		values = append(values, country.Value)
	}
	return values
}

// GetEnumList 返回所有 Country 的列表
func GetEnumList() []Country {
	return Countries
}
