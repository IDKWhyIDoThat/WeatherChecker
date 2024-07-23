package weather

func GetCityWeatherData(city string, outputformat int, valueformat int) (string, error) {
	weatherData, err := GetWeather(replaceSpaces(city, '-'))
	if err != nil {
		return "Что-то пошло не по плану. Вероятно, сервер не отвечает. Повторите попытку позже.", err
	}
	response, err := WeatherDataOutputFormat(*weatherData, outputformat, valueformat)
	if err != nil {
		return "Неизвестный город/местность", err
	}
	return response, err
}
