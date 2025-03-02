package service

import (
	"chat/internal/models"
	"chat/internal/redis"
)

type ScenarioHandler interface {
	CreateOptions(userId string, redisClient *redis.Client) models.Message
	ExecuteJobs(userId string, mess string) models.Message
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
	case TransferToService.Code:
		return &TransferToServiceHandler{}
	default:
		return &DefaultHandler{}
	}
}
