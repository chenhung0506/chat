package service

import (
	"chat/internal/models"
)

type DefaultHandler struct{}

func (d *DefaultHandler) CreateOptions() models.Response {
	return models.Response{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (d *DefaultHandler) ExecuteJobs() string {
	return "Executing default jobs..."
}

// WeeklyWeatherHandler 实现
type WeeklyWeatherHandler struct{}

func (w *WeeklyWeatherHandler) CreateOptions() models.Response {
	return models.Response{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       models.GetCountryValue(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (w *WeeklyWeatherHandler) ExecuteJobs() string {
	return "Fetching weekly weather forecast..."
}

// TomorrowWeatherHandler 实现
type TomorrowWeatherHandler struct{}

func (t *TomorrowWeatherHandler) CreateOptions() models.Response {
	return models.Response{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (t *TomorrowWeatherHandler) ExecuteJobs() string {
	return "Fetching tomorrow's weather forecast..."
}

// LeaveMessageHandler 实现
type LeaveMessageHandler struct{}

func (l *LeaveMessageHandler) CreateOptions() models.Response {
	return models.Response{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (l *LeaveMessageHandler) ExecuteJobs() string {
	return "Handling leave message functionality..."
}

// SoaredStocksHandler 实现
type SoaredStocksHandler struct{}

func (s *SoaredStocksHandler) CreateOptions() models.Response {
	return models.Response{
		Value:      "歡迎使用以下功能:",
		IsFinished: true,
		Data:       GetScenarioValues(),
		SubType:    "relatelist",
		Type:       "text",
	}
}

func (s *SoaredStocksHandler) ExecuteJobs() string {
	return "Fetching recent soared stocks data..."
}
