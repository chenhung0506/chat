package service

import (
	"chat/internal/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func DecodeUnicode(input string) string {
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
