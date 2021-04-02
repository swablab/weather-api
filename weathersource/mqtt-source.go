package weathersource

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"weather-data/config"
	"weather-data/storage"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var mqttTopicRegexPattern = "(^sensor/)([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})(/(temp|pressure|humidity|co2level)$)"

var regexTopic *regexp.Regexp = regexp.MustCompile(mqttTopicRegexPattern)

type mqttWeatherSource struct {
	config                config.MqttConfig
	mqttClient            mqtt.Client
	lastWeatherDataPoints []*storage.WeatherData
	weatherSource         WeatherSourceBase
}

//Close mqtt client
func (source *mqttWeatherSource) Close() {
	source.mqttClient.Disconnect(2)
}

//NewMqttSource Factory function for mqttWeatherSource with authentication
func NewMqttSource(cfg config.MqttConfig) (*mqttWeatherSource, error) {
	source := new(mqttWeatherSource)
	source.config = cfg

	opts := mqtt.NewClientOptions().AddBroker(cfg.Host)

	//mqtt
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(source.mqttMessageHandler())
	opts.SetPingTimeout(1 * time.Second)

	if !cfg.AllowAnonymousAuthentication {
		opts.Username = cfg.Username
		opts.Password = cfg.Password
	}

	source.mqttClient = mqtt.NewClient(opts)

	if token := source.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	if token := source.mqttClient.Subscribe(cfg.Topic, 2, nil); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	go source.publishDataValues()

	return source, nil
}

//mqttMessageHandler returns a function that handles incoming mqtt-messages
func (source *mqttWeatherSource) mqttMessageHandler() mqtt.MessageHandler {

	return func(client mqtt.Client, msg mqtt.Message) {
		if !regexTopic.MatchString(msg.Topic()) {
			return
		}

		sensorId, err := uuid.Parse(regexTopic.FindStringSubmatch(msg.Topic())[2])
		if err != nil {
			return
		}

		lastWeatherData, found := source.getUnwrittenDatapoints(sensorId)

		if !found {
			lastWeatherData = new(storage.WeatherData)
			lastWeatherData.SensorId = sensorId
			source.lastWeatherDataPoints = append(source.lastWeatherDataPoints, lastWeatherData)
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

func (source *mqttWeatherSource) publishDataValues() {
	for {
		for len(source.lastWeatherDataPoints) != 0 {
			current := *source.lastWeatherDataPoints[0]
			diff := time.Now().Sub(current.TimeStamp)
			if diff >= source.config.MinDistToLastValue {
				source.newWeatherData(current)
				source.lastWeatherDataPoints = source.lastWeatherDataPoints[1:]
			}

		}
		time.Sleep(source.config.PublishInterval)
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
