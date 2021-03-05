package main

import (
	"os"
	"weather-data/api"
	"weather-data/storage"
	"weather-data/weathersource"
)

// const influx stuff
const influxToken = "Pg34RXv4QE488ayCeY6JX4p3EwcoNhLu-zPQDn9zxirFmc0og9DCgamf02jrVEAN9mS4mT05nprGUkSrKQAUjA=="
const influxWeatherBucket = "weatherdata"
const influxOrganization = "weather-org"
const influxURL = "https://influx.gamlo-cloud.de"

//const mqtt stuff
const mqttURL = "tcp://gamlo-cloud.de:1883"
const mqttTopic = "sensor/#"
const defaultLocation = "default location"

//const api stuff
const apiAddress = ":10000"

func main() {
	//setup a new weatherstorage -> InfluxDB
	var weatherStorage storage.WeatherStorage
	weatherStorage, err := storage.NewInfluxStorage(influxToken, influxWeatherBucket, influxOrganization, influxURL)
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
	weatherSource, err = weathersource.NewMqttSource(mqttURL, mqttTopic, defaultLocation)
	if err != nil {
		os.Exit(1)
	}
	defer weatherSource.Close()
	weatherSource.AddNewWeatherDataCallback(newWeatherDataHandler)

	//setup a API -> REST
	var weatherAPI api.WeatherAPI
	weatherAPI = api.NewRestAPI(apiAddress, weatherStorage)

	err = weatherAPI.Start()
	if err != nil {
		os.Exit(1)
	}
	defer weatherAPI.Close()
}
