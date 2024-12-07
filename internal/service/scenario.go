package service

type ScenarioEnum struct {
	Code  int
	Desc  string
	Value string
}

var (
	DefaultScenario         = ScenarioEnum{Code: 1, Desc: "Default", Value: "功能列表"}
	WeeklyWeatherScenario   = ScenarioEnum{Code: 2, Desc: "Weekly Weather", Value: "一週天氣"}
	TomorrowWeatherScenario = ScenarioEnum{Code: 3, Desc: "Tomorrow Weather", Value: "明日天氣預報"}
	LeaveMessageScenario    = ScenarioEnum{Code: 4, Desc: "Leave message", Value: "留言"}
	SoaredStocks            = ScenarioEnum{Code: 5, Desc: "Soared Stocks", Value: "近期飆股"}
)

var Scenarios = []ScenarioEnum{
	DefaultScenario,
	WeeklyWeatherScenario,
	TomorrowWeatherScenario,
	LeaveMessageScenario,
	SoaredStocks,
}

func GetScenarioByDesc(desc string) (ScenarioEnum, bool) {
	for _, scenario := range Scenarios {
		if scenario.Value == desc {
			return scenario, true
		}
	}
	return DefaultScenario, false
}

func GetScenarioValues() []string {
	var values []string
	for _, scenario := range Scenarios {
		if scenario.Code != 1 {
			values = append(values, scenario.Value)
		}
	}
	return values
}
