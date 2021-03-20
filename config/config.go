package config

import (
	"os"
	"strconv"
	"time"
)

// const influx stuff
const influxToken = "default-token"
const influxWeatherBucket = "default-bucket"
const influxOrganization = "default-org"
const influxURL = "https://influx.default-address.com"

//const mqtt stuff
const mqttURL = "tcp://default-address.com:1883"
const mqttTopic = "sensor/#"
const defaultLocation = "default-location"
const mqttPublishInterval = time.Second
const mqttMinDistToLastValue = 250 * time.Millisecond

//other config stuff
const allowUnregisteredSensors = false

//influx config
func GetInfluxUrl() string {
	return getVariableWithDefault("WEATHER-API-INFLUX_URL", influxURL)
}

func GetInfluxToken() string {
	return getVariableWithDefault("WEATHER-API-INFLUX_TOKEN", influxToken)
}

func GetInfluxOrganization() string {
	return getVariableWithDefault("WEATHER-API-INFLUX_ORG", influxOrganization)
}

func GetInfluxBucket() string {
	return getVariableWithDefault("WEATHER-API-INFLUX_BUCKET", influxWeatherBucket)
}

//mqtt config
func GetMqttUrl() string {
	return getVariableWithDefault("WEATHER-API-MQTT_URL", mqttURL)
}

func GetMqttTopic() string {
	return getVariableWithDefault("WEATHER-API-MQTT_TOPIC", mqttTopic)
}

func MqttPublishInterval() time.Duration {
	interval, err := strconv.ParseInt(os.Getenv("WEATHER-API-MQTT_PUBLISH_INTERVAL"), 10, 64)
	if err != nil {
		return mqttPublishInterval
	}
	return time.Millisecond * time.Duration(interval)
}

func MqttMinDistToLastValue() time.Duration {
	interval, err := strconv.ParseInt(os.Getenv("WEATHER-API-MQTT_MIN_DIST_TO_LAST_VALUE"), 10, 64)
	if err != nil {
		return mqttMinDistToLastValue
	}
	return time.Millisecond * time.Duration(interval)
}

//common config
func AllowUnregisteredSensors() bool {
	allow, err := strconv.ParseBool(os.Getenv("WEATHER-API-ALLOW_UNREGISTERED_SENSORS"))
	if err != nil {
		return allowUnregisteredSensors
	}
	return allow
}

//helper
func getVariableWithDefault(variableKey, defaultValue string) string {
	variable := os.Getenv(variableKey)
	if len(variable) == 0 {
		return defaultValue
	}
	return variable
}
