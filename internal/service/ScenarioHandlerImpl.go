package service

import (
	"chat/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DefaultHandler struct{}

func (d *DefaultHandler) CreateOptions() models.Message {
	return models.Message{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (d *DefaultHandler) ExecuteJobs(mess string) models.Message {
	return models.Message{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

// WeeklyWeatherHandler 实现
type WeeklyWeatherHandler struct{}

func (w *WeeklyWeatherHandler) CreateOptions() models.Message {
	return models.Message{
		Value:      "請選擇城市:",
		Code:       2,
		IsFinished: false,
		Data:       models.GetCountryValue(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (w *WeeklyWeatherHandler) ExecuteJobs(mess string) models.Message {
	apiURL := "http://139.162.2.175:3001/weather?country=" + models.FromValue(mess).Code
	weatherResponse, err := ParseWeatherAPI(apiURL)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return models.Message{
		Value:      fmt.Sprintf("%v", weatherResponse.Result),
		Code:       2,
		IsFinished: true,
		SubType:    "text",
		Type:       "text",
	}
}

func DecodeUnicode(input string) string {
	// 將 Unicode 字符串轉換為可讀的文本
	var result string
	err := json.Unmarshal([]byte(fmt.Sprintf("\"%s\"", input)), &result)
	if err != nil {
		log.Printf("Failed to decode unicode: %v", err)
		return input
	}
	return result
}

func ParseWeatherAPI(apiURL string) (*models.WeatherResponse, error) {
	response, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var weatherResponse models.WeatherResponse
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// 解碼 Unicode 字符
	for i, weather := range weatherResponse.Result {
		weatherResponse.Result[i].Date = DecodeUnicode(weather.Date)
		weatherResponse.Result[i].Morning = DecodeUnicode(weather.Morning)
		weatherResponse.Result[i].Night = DecodeUnicode(weather.Night)
	}

	return &weatherResponse, nil
}

type TomorrowWeatherHandler struct{}

func (t *TomorrowWeatherHandler) CreateOptions() models.Message {
	return models.Message{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (t *TomorrowWeatherHandler) ExecuteJobs(mess string) models.Message {
	return models.Message{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

// LeaveMessageHandler 实现
type LeaveMessageHandler struct{}

func (l *LeaveMessageHandler) CreateOptions() models.Message {
	return models.Message{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (l *LeaveMessageHandler) ExecuteJobs(mess string) models.Message {
	return models.Message{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

// SoaredStocksHandler 实现
type SoaredStocksHandler struct{}

func (s *SoaredStocksHandler) CreateOptions() models.Message {
	return models.Message{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (s *SoaredStocksHandler) ExecuteJobs(mess string) models.Message {
	return models.Message{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}
