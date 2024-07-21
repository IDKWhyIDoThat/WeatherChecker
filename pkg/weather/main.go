package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"weatherbot/pkg/additional/getsmth"
)

type WeatherData struct {
	Location struct {
		City string `json:"name"`
	} `json:"location"`
	Current struct {
		TempC     Celsium   `json:"temp_c"`
		TempF     Farenheit `json:"temp_f"`
		Condition struct {
			Status WeatherStatus `json:"text"`
		} `json:"condition"`
		Windk        WindSpeedInKPH        `json:"wind_kph"`
		Windm        WindSpeedInMPH        `json:"wind_mph"`
		Cloud        CloudsCoverPercentage `json:"cloud"`
		Pressure     PreesureInmBars       `json:"pressure_mb"`
		Humidity     HumidityPercent       `json:"humidity"`
		Visabilitykm VisibilityRangeInKm   `json:"vis_km"`
		Visibilitymi VisibilityRangeInMi   `json:"vis_miles"`
	} `json:"current"`
}

func GetWeather(city string) (*WeatherData, error) {
	apiKey, err := getsmth.GetAPIkey("http://api.weatherapi.com")

	if err != nil {
		return nil, fmt.Errorf("ошибка при получении токена: %s", err)
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении данных: %s", err)
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, fmt.Errorf("ошибка при форматировании .json данных: %s", err)
	}

	return &weatherData, err
}

func WeatherDataOutputFormat(weatherData WeatherData, outputformat int, valueformat int) (string, error) {
	if weatherData.Location.City == "" {
		return "", fmt.Errorf("unknown city name")
	}
	switch outputformat {
	case 1:
		return fmt.Sprintf("%s:\n\t-Temreture: %v\n\t-Humidity: %v\n\t-Wind: %v",
			weatherData.Location.City, weatherData.Temp(valueformat), weatherData.Current.Humidity,
			weatherData.Wind(valueformat)), nil
	case 2:
		return fmt.Sprintf("%s:\n\t-Weather: %v\n\t-Temreture: %v\n\t-Humidity: %v\n\t-Wind: %v\n\t-Pressure: %v\n\t-Visability: %v\n\t-Clouds: %v",
			weatherData.Location.City, weatherData.Current.Condition.Status, weatherData.Temp(valueformat), weatherData.Current.Humidity,
			weatherData.Wind(valueformat), weatherData.Current.Pressure, weatherData.Visibility(valueformat), weatherData.Current.Cloud), nil
	default:
		return "", fmt.Errorf("WeatherDataOutputFormat: wrong format (argument) number %d", outputformat)
	}
}
