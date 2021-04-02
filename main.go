package main

import (
	"errors"
	"fmt"
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
	//setup new sensorRegistry -> MongodbSensorRegistry
	var err error
	sensorRegistry, err = storage.NewMongodbSensorRegistry(
		config.GetMongodbURL(),
		config.GetMongodbName(),
		config.GetMongodbCollection(),
		config.GetMongodbUserName(),
		config.GetMongodbPassword())

	if err != nil {
		os.Exit(1)
	}
	defer sensorRegistry.Close()

	//setup a new weatherstorage -> InfluxDB
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
	if config.UseAnonymousMqttAuthentication() {
		weatherSource, err = weathersource.NewAnonymousMqttSource(
			config.GetMqttUrl(),
			config.GetMqttTopic())
	} else {
		weatherSource, err = weathersource.NewMqttSource(
			config.GetMqttUrl(),
			config.GetMqttTopic(),
			config.GetMqttUser(),
			config.GetMqttPassword())
	}
	if err != nil {
		fmt.Println("Could not connect to mqtt:", err.Error())
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

func handleNewWeatherData(wd storage.WeatherData) error {
	_, couldResolve := sensorRegistry.ResolveSensorById(wd.SensorId)
	if !config.AllowUnregisteredSensors() && !couldResolve {
		return errors.New("sensor have to be registered")
	}
	weatherStorage.Save(wd)
	return nil
}
