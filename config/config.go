package config

import (
	"os"
	"strconv"
	"time"
)

type MongoConfig struct {
	Host       string
	Database   string
	Username   string
	Password   string
	Collection string
}

type InfluxConfig struct {
	Host         string
	Token        string
	Organization string
	Bucket       string
}

type MqttConfig struct {
	Host                         string
	Topic                        string
	Username                     string
	Password                     string
	PublishDelay                 time.Duration
	AllowAnonymousAuthentication bool
}

type RestConfig struct {
	AccessControlAllowOriginHeader string
	UseTokenAuthorization          bool
	ValidateTokenUrl               string
}

var MongoConfiguration = MongoConfig{
	Host:       getEnv("MONGO_HOST", "localhost:27017"),
	Database:   getEnv("MONGO_DB", "weathersensors"),
	Username:   getEnv("MONGO_USER", "admin"),
	Password:   getEnv("MONGO_PASS", "admin"),
	Collection: getEnv("MONGO_COLLECTION", "sensors"),
}

var InfluxConfiguration = InfluxConfig{
	Host:         getEnv("INFLUX_HOST", "localhost:8086"),
	Token:        getEnv("INFLUX_TOKEN", "token"),
	Organization: getEnv("INFLUX_ORG", "org_name"),
	Bucket:       getEnv("INFLUX_BUCKET", "bucket_name"),
}

var MqttConfiguration = MqttConfig{
	Host:                         getEnv("MQTT_HOST", "localhost:1883"),
	Topic:                        getEnv("MQTT_TOPIC", "sensor/#"),
	Username:                     getEnv("MQTT_USER", "mqtt"),
	Password:                     getEnv("MQTT_PASS", "mqtt"),
	PublishDelay:                 getEnvDuration("MQTT_PUBLISH_DELAY", time.Second),
	AllowAnonymousAuthentication: getEnvBool("MQTT_ANONYMOUS", false),
}

var RestConfiguration = RestConfig{
	AccessControlAllowOriginHeader: getEnv("ACCESS_CONTROL_ALLOW_ORIGIN_HEADER", "*"),
	UseTokenAuthorization:          getEnvBool("USE_TOKEN_AUTHORIZATION", false),
	ValidateTokenUrl:               getEnv("JWT_TOKEN_VALIDATION_URL", "https://api.swablab.de/ldap/validateToken"),
}

var AllowUnregisteredSensors = getEnvBool("ALLOW_UNREGISTERED_SENSORS", false)

//helper
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvBool(key string, fallback bool) bool {

	if value, ok := os.LookupEnv(key); ok {
		if bValue, err := strconv.ParseBool(value); err == nil {
			return bValue
		}
	}

	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		if iValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return iValue
		}
	}

	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if iValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Millisecond * time.Duration(iValue)
		}
	}

	return fallback
}
