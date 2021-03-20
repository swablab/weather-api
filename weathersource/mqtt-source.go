package weathersource

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"weather-data/storage"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var uuidRegexPattern = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
var mqttTopicRegexPattern = fmt.Sprintf("^sensor/%s/(temp|pressure|humidity|co2level)$", uuidRegexPattern)

var regexTopic *regexp.Regexp = regexp.MustCompile(mqttTopicRegexPattern)
var regexUuid *regexp.Regexp = regexp.MustCompile(uuidRegexPattern)

type mqttWeatherSource struct {
	url                   string
	topic                 string
	mqttClient            mqtt.Client
	lastWeatherDataPoints []*storage.WeatherData
	weatherSource         WeatherSourceBase
}

//Close mqtt client
func (source *mqttWeatherSource) Close() {
	source.mqttClient.Disconnect(2)
}

//NewMqttSource Factory function for mqttWeatherSource
func NewMqttSource(url, topic string) (*mqttWeatherSource, error) {
	source := new(mqttWeatherSource)
	source.url = url

	opts := mqtt.NewClientOptions().AddBroker(url)

	//mqtt
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(source.mqttMessageHandler())
	opts.SetPingTimeout(1 * time.Second)

	source.mqttClient = mqtt.NewClient(opts)

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
		if !regexTopic.MatchString(msg.Topic()) {
			return
		}

		sensorId, err := uuid.Parse(regexUuid.FindAllString(msg.Topic(), 1)[0])
		if err != nil {
			return
		}

		lastWeatherData, found := source.getUnwrittenDatapoints(sensorId)

		if !found {
			lastWeatherData = new(storage.WeatherData)
			lastWeatherData.SensorId = sensorId
			source.lastWeatherDataPoints = append(source.lastWeatherDataPoints, lastWeatherData)
		}

		diff := time.Now().Sub(lastWeatherData.TimeStamp)
		if diff >= time.Second && diff < time.Hour*6 {
			source.newWeatherData(*lastWeatherData)
		}

		if strings.HasSuffix(msg.Topic(), "pressure") {
			lastWeatherData.Pressure, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			lastWeatherData.TimeStamp = time.Now()
		}
		if strings.HasSuffix(msg.Topic(), "temp") {
			lastWeatherData.Temperature, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			lastWeatherData.TimeStamp = time.Now()
		}
		if strings.HasSuffix(msg.Topic(), "humidity") {
			lastWeatherData.Temperature, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			lastWeatherData.TimeStamp = time.Now()
		}
		if strings.HasSuffix(msg.Topic(), "co2level") {
			lastWeatherData.CO2Level, _ = strconv.ParseFloat(string(msg.Payload()), 64)
			lastWeatherData.TimeStamp = time.Now()
		}
	}
}

func (source *mqttWeatherSource) getUnwrittenDatapoints(sensorId uuid.UUID) (*storage.WeatherData, bool) {
	for _, data := range source.lastWeatherDataPoints {
		if data.SensorId == sensorId {
			return data, true
		}
	}
	return nil, false
}

//AddNewWeatherDataCallback adds a new callbackMethod for incoming weather data
func (source *mqttWeatherSource) AddNewWeatherDataCallback(callback NewWeatherDataCallbackFunc) {
	source.weatherSource.AddNewWeatherDataCallback(callback)
}

func (source *mqttWeatherSource) newWeatherData(datapoint storage.WeatherData) error {
	return source.weatherSource.NewWeatherData(datapoint)
}
