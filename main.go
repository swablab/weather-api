package main

import (
	"os"
	"weather-data/api"
	"weather-data/config"
	"weather-data/storage"
	"weather-data/weathersource"
)

func main() {
	//setup a new weatherstorage -> InfluxDB
	var weatherStorage storage.WeatherStorage
	weatherStorage, err := storage.NewInfluxStorage(
		config.GetInfluxToken(),
		config.GetInfluxBucket(),
		config.GetInfluxOrganization(),
		config.GetInfluxUrl())

	if err != nil {
		os.Exit(1)
	}
	defer weatherStorage.Close()

	var newWeatherDataHandler weathersource.NewWeatherDataCallbackFunc
	newWeatherDataHandler = func(wd storage.WeatherData) {
		weatherStorage.Save(wd)
	}

	//add a new weatherData source -> mqtt
	var weatherSource weathersource.WeatherSource
	weatherSource, err = weathersource.NewMqttSource(
		config.GetMqttUrl(),
		config.GetMqttTopic(),
		config.GetMqttLocation())

	if err != nil {
		os.Exit(1)
	}
	defer weatherSource.Close()

	weatherSource.AddNewWeatherDataCallback(newWeatherDataHandler)

	//setup a API -> REST
	var weatherAPI api.WeatherAPI
	weatherAPI = api.NewRestAPI(":10000", weatherStorage)
	defer weatherAPI.Close()

	err = weatherAPI.Start()
	if err != nil {
		os.Exit(1)
	}

}
