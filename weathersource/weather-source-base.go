package weathersource

import "weather-data/storage"

//WeatherSourceBase is the lowlevel-implementation of the WeatherSource interface, intended to used by highlevel-implementations
type WeatherSourceBase struct {
	newWeatherDataCallbackFuncs []NewWeatherDataCallbackFunc
}

//AddNewWeatherDataCallback adds a new callbackMethod for incoming weather data
func (source *WeatherSourceBase) AddNewWeatherDataCallback(callback NewWeatherDataCallbackFunc) {
	source.newWeatherDataCallbackFuncs = append(source.newWeatherDataCallbackFuncs, callback)
}

//NewWeatherData executes all newWeatherDataCallbackFuncs for this datapoint
func (source *WeatherSourceBase) NewWeatherData(datapoint storage.WeatherData) {
	for _, callback := range source.newWeatherDataCallbackFuncs {
		callback(datapoint)
	}
}
