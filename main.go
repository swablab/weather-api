package main

import (
	"log"
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
	log.SetOutput(os.Stdout)

	//setup new sensorRegistry -> MongodbSensorRegistry
	var err error
	if sensorRegistry, err = storage.NewMongodbSensorRegistry(config.MongoConfiguration); err != nil {
		log.Fatal(err)
	}
	defer sensorRegistry.Close()

	//setup a new weatherstorage -> InfluxDB
	if weatherStorage, err = storage.NewInfluxStorage(config.InfluxConfiguration); err != nil {
		log.Fatal(err)
	}
	defer weatherStorage.Close()

	//setup new weatherData source -> mqtt
	if weatherSource, err = weathersource.NewMqttSource(config.MqttConfiguration); err != nil {
		log.Fatal(err)
	}
	defer weatherSource.Close()
	weatherSource.OnNewWeatherData(handleNewWeatherData)

	//setup a API -> REST
	weatherAPI = api.NewRestAPI(":10000", weatherStorage, sensorRegistry, config.RestConfiguration)
	defer weatherAPI.Close()
	weatherAPI.OnNewWeatherData(handleNewWeatherData)

	log.Print("Application is running")
	err = weatherAPI.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func handleNewWeatherData(wd *storage.WeatherData) {
	if config.AllowUnregisteredSensors {
		weatherStorage.Save(wd)
	} else if exist, err := sensorRegistry.ExistSensor(wd.SensorId); err == nil && exist {
		weatherStorage.Save(wd)
	}
}
