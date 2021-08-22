package weathersource

import (
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"
	"weather-data/config"
	"weather-data/storage"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var mqttTopicRegexPattern = "(^sensor)/([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})/(.*)"

var regexTopic *regexp.Regexp = regexp.MustCompile(mqttTopicRegexPattern)

var channelBufferSize = 10

type mqttWeatherSource struct {
	config                   config.MqttConfig
	mqttClient               mqtt.Client
	weatherSource            WeatherSourceBase
	activeSensorMeasurements map[uuid.UUID](chan map[storage.SensorValueType]float64)
	sensorMutex              sync.RWMutex
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
	opts.SetDefaultPublishHandler(source.mqttMessageHandler)
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

	source.activeSensorMeasurements = make(map[uuid.UUID]chan map[storage.SensorValueType]float64)
	source.sensorMutex = sync.RWMutex{}

	log.Print("successfully connected to mqtt-broker")
	return source, nil
}

//mqttMessageHandler returns a function that handles incoming mqtt-messages
func (source *mqttWeatherSource) mqttMessageHandler(client mqtt.Client, msg mqtt.Message) {
	if !regexTopic.MatchString(msg.Topic()) {
		return
	}

	sensorId, err := uuid.Parse(regexTopic.FindStringSubmatch(msg.Topic())[2])
	if err != nil {
		return
	}

	value, err := strconv.ParseFloat(string(msg.Payload()), 64)
	if err != nil {
		return
	}

	sensorValueType := storage.SensorValueType(regexTopic.FindStringSubmatch(msg.Topic())[3])

	dataValue := map[storage.SensorValueType]float64{
		sensorValueType: value,
	}

	source.sensorMutex.RLock()
	dataChannel, exists := source.activeSensorMeasurements[sensorId]
	if !exists {
		dataChannel = make(chan map[storage.SensorValueType]float64, channelBufferSize)
	}
	dataChannel <- dataValue
	source.sensorMutex.RUnlock()

	if !exists {
		go source.publishSensorMeasurement(sensorId, dataChannel)
		go source.cleanupSensorMeasurement(sensorId, dataChannel)

		source.sensorMutex.Lock()
		source.activeSensorMeasurements[sensorId] = dataChannel
		source.sensorMutex.Unlock()
	}
}

func (source *mqttWeatherSource) cleanupSensorMeasurement(sensorId uuid.UUID, channel chan<- map[storage.SensorValueType]float64) {
	time.Sleep(source.config.PublishDelay)

	source.sensorMutex.Lock()
	delete(source.activeSensorMeasurements, sensorId)
	source.sensorMutex.Unlock()

	close(channel)
}

func (source *mqttWeatherSource) publishSensorMeasurement(sensorId uuid.UUID, channel <-chan map[storage.SensorValueType]float64) {
	weatherData := storage.NewWeatherData()
	weatherData.TimeStamp = time.Now()
	weatherData.SensorId = sensorId

	for values := range channel {
		for k, v := range values {
			weatherData.Values[k] = v
		}
	}

	source.newWeatherData(*weatherData)
}

//AddNewWeatherDataCallback adds a new callbackMethod for incoming weather data
func (source *mqttWeatherSource) AddNewWeatherDataCallback(callback NewWeatherDataCallbackFunc) {
	source.weatherSource.AddNewWeatherDataCallback(callback)
}

func (source *mqttWeatherSource) newWeatherData(datapoint storage.WeatherData) {
	source.weatherSource.NewWeatherData(datapoint)
}
