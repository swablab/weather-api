package weathersource

import "weather-data/storage"

//NewWeatherDataFunc Function-Signature for new weather data
type NewWeatherDataFunc func(*storage.WeatherData)

//WeatherSource is the interface for different weather-source implementations
type WeatherSource interface {
	OnNewWeatherData(callback NewWeatherDataFunc)
	Close()
}

//WeatherSourceBase is the lowlevel-implementation of the WeatherSource interface, intended to used by highlevel-implementations
type WeatherSourceBase struct {
	onNewWeatherDataFunctions []NewWeatherDataFunc
}

//OnNewWeatherData add a function executed on NewWeatherData called
func (source *WeatherSourceBase) OnNewWeatherData(callback NewWeatherDataFunc) {
	source.onNewWeatherDataFunctions = append(source.onNewWeatherDataFunctions, callback)
}

//NewWeatherData executes all NewWeatherDataFunc for the weatherData
func (source *WeatherSourceBase) NewWeatherData(weatherData *storage.WeatherData) {
	for _, function := range source.onNewWeatherDataFunctions {
		function(weatherData)
	}
}
