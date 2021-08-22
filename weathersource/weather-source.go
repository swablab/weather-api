package weathersource

import "weather-data/storage"

//NewWeatherDataCallbackFunc Function-Signature for new weather data callback function
type NewWeatherDataCallbackFunc func(storage.WeatherData)

//WeatherSource is the interface for different weather-source implementations
type WeatherSource interface {
	AddNewWeatherDataCallback(NewWeatherDataCallbackFunc)
	Close()
}
