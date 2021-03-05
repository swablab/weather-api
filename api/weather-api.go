package api

import "weather-data/weathersource"

//WeatherAPI is the common interface for different apis
type WeatherAPI interface {
	Start() error
	Close()
	weathersource.WeatherSource
}
