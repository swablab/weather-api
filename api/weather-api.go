package api

import (
	"weather-data/weathersource"

	"github.com/google/uuid"
)

//WeatherAPI is the common interface for different apis
type WeatherAPI interface {
	Start() error
	Close()
	weathersource.WeatherSource
}

//SensorRegistration is the data for a new Sensorregistration
type SensorRegistration struct {
	Name     string
	Id       uuid.UUID
	Location string
}
