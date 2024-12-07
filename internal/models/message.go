package models

type Message struct {
	UUID       string   `json:"uuid"`
	SessionId  string   `json:"sessionId"`
	Mess       string   `json:"mess"`
	Domain     string   `json:"domain"`
	Code       int      `json:"code"`
	IsFinished bool     `json:"isFinished"`
	Value      string   `json:"value"`
	Data       []string `json:"data"`
	SubType    string   `json:"subType"`
	Type       string   `json:"type"`
}

type Messages struct {
	Messages []Message
}

func NewInitialMessages(userId string) Messages {
	return Messages{
		Messages: []Message{
			{
				UUID:       userId,
				SessionId:  "init-session",
				Mess:       "功能列表",
				Domain:     "example.com",
				Code:       1,
				IsFinished: true,
			},
		},
	}
}

func (m *Messages) AddMessage(newMessage Message) {
	m.Messages = append(m.Messages, newMessage)
}

func (m *Messages) GetPreviousMessage() Message {
	if len(m.Messages) > 1 {
		return m.Messages[len(m.Messages)-1]
	}
	return Message{}
}
