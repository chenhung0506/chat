package service

import (
	"chat/internal/models"
)

type ScenarioHandler interface {
	CreateOptions() models.Message
	ExecuteJobs(mess string) models.Message
}

func GetScenarioHandler(scenario ScenarioEnum) ScenarioHandler {
	switch scenario.Code {
	case WeeklyWeatherScenario.Code:
		return &WeeklyWeatherHandler{}
	case TomorrowWeatherScenario.Code:
		return &TomorrowWeatherHandler{}
	case LeaveMessageScenario.Code:
		return &LeaveMessageHandler{}
	case SoaredStocks.Code:
		return &SoaredStocksHandler{}
	default:
		return &DefaultHandler{}
	}
}
