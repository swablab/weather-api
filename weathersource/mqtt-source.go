package weathersource

import (
	"strconv"
	"strings"
	"time"
	"weather-data/storage"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttWeatherSource struct {
	url           string
	topic         string
	mqttClient    mqtt.Client
	lastData      storage.WeatherData
	weatherSource WeatherSourceBase
}

//Close mqtt client
func (source *mqttWeatherSource) Close() {
	source.mqttClient.Disconnect(2)
}

//NewMqttSource Factory function for mqttWeatherSource
func NewMqttSource(url, topic, defaultLocation string) (*mqttWeatherSource, error) {
	source := new(mqttWeatherSource)
	source.url = url

	opts := mqtt.NewClientOptions().AddBroker(url)

	//mqtt
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(source.mqttMessageHandler())
	opts.SetPingTimeout(1 * time.Second)

	source.mqttClient = mqtt.NewClient(opts)
	source.lastData.Location = defaultLocation

	if token := source.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := source.mqttClient.Subscribe(topic, 2, nil); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return source, nil
}

//mqttMessageHandler returns a function that handles incoming mqtt-messages
func (source *mqttWeatherSource) mqttMessageHandler() mqtt.MessageHandler {

	return func(client mqtt.Client, msg mqtt.Message) {

		diff := time.Now().Sub(source.lastData.TimeStamp)
		if diff >= time.Second && diff < time.Hour*6 {
			source.newWeatherData(source.lastData)
		}

		if strings.HasSuffix(msg.Topic(), "pressure") {
			source.lastData.Pressure, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			source.lastData.TimeStamp = time.Now()
		}
		if strings.HasSuffix(msg.Topic(), "temp") {
			source.lastData.Temperature, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			source.lastData.TimeStamp = time.Now()
		}
		if strings.HasSuffix(msg.Topic(), "humidity") {
			source.lastData.Temperature, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			source.lastData.TimeStamp = time.Now()
		}
		if strings.HasSuffix(msg.Topic(), "co2level") {
			source.lastData.CO2Level, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			source.lastData.TimeStamp = time.Now()
		}
	}
}

//AddNewWeatherDataCallback adds a new callbackMethod for incoming weather data
func (source *mqttWeatherSource) AddNewWeatherDataCallback(callback NewWeatherDataCallbackFunc) {
	source.weatherSource.AddNewWeatherDataCallback(callback)
}

func (source *mqttWeatherSource) newWeatherData(datapoint storage.WeatherData) {
	for _, callback := range source.weatherSource.newWeatherDataCallbackFuncs {
		callback(datapoint)
	}
}
