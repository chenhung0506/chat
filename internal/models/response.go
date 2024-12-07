package models

type Response struct {
	Value      string   `json:"value"`
	IsFinished bool     `json:"isFinished"`
	Data       []string `json:"data"`
	SubType    string   `json:"subType"`
	Type       string   `json:"type"`
}
