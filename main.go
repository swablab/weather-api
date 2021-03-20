package main

import (
	"os"
	"weather-data/api"
	"weather-data/config"
	"weather-data/storage"
	"weather-data/weathersource"
)

var sensorRegistry storage.SensorRegistry
var weatherStorage storage.WeatherStorage
var weatherSource weathersource.WeatherSource
var weatherAPI api.WeatherAPI

func main() {
	//setup new sensorRegistry -> InmemorySensorRegistry
	sensorRegistry = storage.NewInmemorySensorRegistry()
	defer sensorRegistry.Close()

	//setup a new weatherstorage -> InfluxDB
	var err error
	weatherStorage, err = storage.NewInfluxStorage(
		config.GetInfluxToken(),
		config.GetInfluxBucket(),
		config.GetInfluxOrganization(),
		config.GetInfluxUrl())

	if err != nil {
		os.Exit(1)
	}
	defer weatherStorage.Close()

	//setup new weatherData source -> mqtt
	weatherSource, err = weathersource.NewMqttSource(
		config.GetMqttUrl(),
		config.GetMqttTopic(),
		sensorRegistry)

	if err != nil {
		os.Exit(1)
	}
	defer weatherSource.Close()
	weatherSource.AddNewWeatherDataCallback(handleNewWeatherData)

	//setup a API -> REST
	weatherAPI = api.NewRestAPI(":10000", weatherStorage, sensorRegistry)
	defer weatherAPI.Close()
	weatherAPI.AddNewWeatherDataCallback(handleNewWeatherData)

	err = weatherAPI.Start()
	if err != nil {
		os.Exit(1)
	}
}

func handleNewWeatherData(wd storage.WeatherData) {
	weatherStorage.Save(wd)
}
