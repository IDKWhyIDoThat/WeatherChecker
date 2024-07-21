package weather

func GetCityWeatherData(city string, outputformat int, valueformat int) (string, error) {
	weatherData, err := GetWeather(replaceSpaces(city, '-'))
	if err != nil {
		return "Something gone terribly wrong", err
	}
	answer, err := WeatherDataOutputFormat(*weatherData, outputformat, valueformat)
	if err != nil {
		return "Unknown city name", err
	}
	return answer, err
}
