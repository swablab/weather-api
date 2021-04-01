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
const mqttUser = "weather-api"
const mqttPassword = "weather-api"
const useAnonymousMqttAuthentication = false
const mqttPublishInterval = time.Second
const mqttMinDistToLastValue = 250 * time.Millisecond

//const mongodb stuff
const mongodbURL = "mongodb://default-address.com:27017"
const mongodbName = "weathersensors"
const mongodbCollection = "sensordata"

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

func GetMqttUser() string {
	return getVariableWithDefault("WEATHER-API-MQTT_USER", mqttTopic)
}

func GetMqttPassword() string {
	return getVariableWithDefault("WEATHER-API-MQTT_PASSWORD", mqttTopic)
}

func UseAnonymousMqttAuthentication() bool {
	return getVariableWithDefaultBool("WEATHER-API-ANONYMOUS_MQTT_AUTHENTICATION", useAnonymousMqttAuthentication)
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

//mongodb config
func GetMongodbURL() string {
	return getVariableWithDefault("WEATHER-API-MONGODB_URL", mongodbURL)
}

func GetMongodbName() string {
	return getVariableWithDefault("WEATHER-API-MONGODB_NAME", mongodbName)
}

func GetMongodbCollection() string {
	return getVariableWithDefault("WEATHER-API-MONGODB_COLLECTION", mongodbCollection)
}

//common config
func AllowUnregisteredSensors() bool {
	return getVariableWithDefaultBool("WEATHER-API-ALLOW_UNREGISTERED_SENSORS", allowUnregisteredSensors)
}

//helper
func getVariableWithDefault(variableKey, defaultValue string) string {
	variable := os.Getenv(variableKey)
	if len(variable) == 0 {
		return defaultValue
	}
	return variable
}

func getVariableWithDefaultBool(variableKey string, defaultValue bool) bool {
	ok, err := strconv.ParseBool(os.Getenv(variableKey))
	if err != nil {
		return defaultValue
	}
	return ok
}
