package main

import (
	"fmt"
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
		os.Exit(1)
	}
	defer sensorRegistry.Close()

	//setup a new weatherstorage -> InfluxDB
	if weatherStorage, err = storage.NewInfluxStorage(config.InfluxConfiguration); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer weatherStorage.Close()

	//setup new weatherData source -> mqtt
	if weatherSource, err = weathersource.NewMqttSource(config.MqttConfiguration); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer weatherSource.Close()
	weatherSource.AddNewWeatherDataCallback(handleNewWeatherData)

	//setup a API -> REST
	weatherAPI = api.NewRestAPI(":10000", weatherStorage, sensorRegistry, config.RestConfiguration)
	defer weatherAPI.Close()
	weatherAPI.AddNewWeatherDataCallback(handleNewWeatherData)

	log.Print("Application is running")
	err = weatherAPI.Start()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func handleNewWeatherData(wd storage.WeatherData) error {
	_, err := sensorRegistry.ResolveSensorById(wd.SensorId)
	if !config.AllowUnregisteredSensors && err != nil {
		log.Print("discarded invalid weatherdata")
		return fmt.Errorf("could not resolve sensor")
	}
	weatherStorage.Save(wd)
	return nil
}
