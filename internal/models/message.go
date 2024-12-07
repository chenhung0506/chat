package models

type Message struct {
	UUID       string `json:"uuid"`
	SessionId  string `json:"sessionId"`
	Mess       string `json:"mess"`
	Domain     string `json:"domain"`
	Code       int    `json:"code"`
	IsFinished bool   `json:"isFinished"`
}

type Messages struct {
	Messages []Message
}
