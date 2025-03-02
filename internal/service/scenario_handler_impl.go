package service

import (
	"chat/internal/models"
	"chat/internal/redis"
	"log"
)

type DefaultHandler struct{}

func (d *DefaultHandler) CreateOptions(userId string, redisClient *redis.Client) models.Message {
	return models.Message{
		UUID:       userId,
		Value:      "歡迎使用以下功能:",
		Code:       1,
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (d *DefaultHandler) ExecuteJobs(userId string, mess string) models.Message {
	return models.Message{
		UUID:       userId,
		Value:      "歡迎使用以下功能:",
		Code:       1,
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

type WeeklyWeatherHandler struct{}

func (w *WeeklyWeatherHandler) CreateOptions(userId string, redisClient *redis.Client) models.Message {
	return models.Message{
		UUID:       userId,
		Value:      "請選擇城市:",
		Code:       3,
		IsFinished: false,
		Data:       models.GetCountryValue(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (w *WeeklyWeatherHandler) ExecuteJobs(userId string, mess string) models.Message {
	apiURL := "http://139.162.2.175:3001/weather?country=" + models.FromValue(mess).Code
	weatherResponse, err := ParseWeatherAPI(apiURL)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	result := ""
	for _, weather := range weatherResponse.Result {
		result = result + weather.Date + " 白天:" + weather.MorningTemp + " 晚上:" + weather.NightTemp + "\n"
	}
	log.Println(weatherResponse.Result)
	return models.Message{
		UUID:       userId,
		Value:      result,
		Code:       2,
		IsFinished: true,
		SubType:    "text",
		Type:       "text",
	}
}

type TomorrowWeatherHandler struct{}

func (t *TomorrowWeatherHandler) CreateOptions(userId string, redisClient *redis.Client) models.Message {
	return models.Message{
		UUID:       userId,
		Value:      "請選擇城市:",
		Code:       3,
		IsFinished: false,
		Data:       models.GetCountryValue(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (t *TomorrowWeatherHandler) ExecuteJobs(userId string, mess string) models.Message {
	apiURL := "http://139.162.2.175:3001/weather?country=" + models.FromValue(mess).Code
	weatherResponse, err := ParseWeatherAPI(apiURL)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	result := ""
	result = weatherResponse.Result[0].Date + " 白天:" + weatherResponse.Result[0].Morning + " 晚上:" + weatherResponse.Result[0].Night

	log.Println(weatherResponse.Result)
	return models.Message{
		UUID:       userId,
		Value:      result,
		Code:       3,
		IsFinished: true,
		SubType:    "text",
		Type:       "text",
	}
}

type LeaveMessageHandler struct{}

func (l *LeaveMessageHandler) CreateOptions(userId string, redisClient *redis.Client) models.Message {
	return models.Message{
		UUID:       userId,
		Value:      "請輸入留言:",
		Code:       4,
		IsFinished: false,
		SubType:    "text",
		Type:       "text",
	}
}

func (l *LeaveMessageHandler) ExecuteJobs(userId string, mess string) models.Message {
	return models.Message{
		UUID:       userId,
		Value:      "已收到你的留言:" + mess,
		Code:       4,
		IsFinished: true,
		SubType:    "text",
		Type:       "text",
	}
}

type SoaredStocksHandler struct{}

func (s *SoaredStocksHandler) CreateOptions(userId string, redisClient *redis.Client) models.Message {
	return models.Message{
		UUID:       userId,
		Value:      "功能修復中:",
		Code:       5,
		IsFinished: true,
		SubType:    "text",
		Type:       "text",
	}
}

func (s *SoaredStocksHandler) ExecuteJobs(userId string, mess string) models.Message {
	return models.Message{
		UUID:       userId,
		Value:      "功能修復中...",
		Code:       5,
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

type TransferToServiceHandler struct{}

func (s *TransferToServiceHandler) CreateOptions(userId string, redisClient *redis.Client) models.Message {
	redisClient.AddWaitingQueue(userId)
	return models.Message{
		UUID:       userId,
		Value:      "已轉接人工客服",
		Code:       6,
		IsFinished: false,
		SubType:    "text",
		Type:       "text",
	}
}

func (s *TransferToServiceHandler) ExecuteJobs(userId string, mess string) models.Message {
	return models.Message{
		UUID:       userId,
		Value:      "功能修復中...",
		Code:       6,
		IsFinished: false,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}
