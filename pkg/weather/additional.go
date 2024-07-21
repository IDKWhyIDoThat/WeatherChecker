package weather

import (
	"fmt"
	"strings"
)

type Celsium float64
type Farenheit float64
type HumidityPercent int
type WindSpeedInKPH float64
type WindSpeedInMPH float64
type WeatherStatus string
type VisibilityRangeInKm float64
type VisibilityRangeInMi float64
type CloudsCoverPercentage float64
type PreesureInmBars float64

func (x Celsium) String() string {
	return fmt.Sprintf("%.1f °C", x)
}
func (x Farenheit) String() string {
	return fmt.Sprintf("%.1f °F", x)
}

func (x HumidityPercent) String() string {
	return fmt.Sprintf("%d %%", x)
}

func (x WindSpeedInKPH) String() string {
	return fmt.Sprintf("%.1f km/h", x)
}
func (x WindSpeedInMPH) String() string {
	return fmt.Sprintf("%.1f miles/h", x)
}

func (x VisibilityRangeInKm) String() string {
	return fmt.Sprintf("%.1f km", x)
}
func (x VisibilityRangeInMi) String() string {
	return fmt.Sprintf("%.1f miles", x)
}

func (x CloudsCoverPercentage) String() string {
	return fmt.Sprintf("%.0f %%", x)
}

func (x PreesureInmBars) String() string {
	return fmt.Sprintf("%.3f B", x/100)
}

func (weatherData WeatherData) Temp(format int) string {
	switch format {
	case 1:
		return fmt.Sprint(weatherData.Current.TempC)
	case 2:
		return fmt.Sprint(weatherData.Current.TempF)
	default:
		return "ErrorTempValue"
	}
}

func (weatherData WeatherData) Wind(format int) string {
	switch format {
	case 1:
		return fmt.Sprint(weatherData.Current.Windk)
	case 2:
		return fmt.Sprint(weatherData.Current.Windm)
	default:
		return "ErrorWindValue"
	}
}

func (weatherData WeatherData) Visibility(format int) string {
	switch format {
	case 1:
		return fmt.Sprint(weatherData.Current.Visabilitykm)
	case 2:
		return fmt.Sprint(weatherData.Current.Visibilitymi)
	default:
		return "ErrorWindValue"
	}
}

func replaceSpaces(input string, replacement rune) string {
	return strings.Map(func(r rune) rune {
		if r == ' ' {
			return replacement
		}
		return r
	}, input)
}
