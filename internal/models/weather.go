package models

type WeatherResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Result  []Weather `json:"result"`
}

type Weather struct {
	Date        string `json:"date"`
	Morning     string `json:"morning"`
	Night       string `json:"night"`
	MorningTemp string `json:"morning-temp"`
	NightTemp   string `json:"night-temp"`
}
