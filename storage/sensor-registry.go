package storage

import "github.com/google/uuid"

type SensorRegistry interface {
	RegisterSensorByName(string) (*WeatherSensor, error)
	ExistSensor(sensorId uuid.UUID) (bool, error)
	ExistSensorName(name string) (bool, error)
	ResolveSensorById(uuid.UUID) (*WeatherSensor, error)
	DeleteSensor(uuid.UUID) error
	UpdateSensor(*WeatherSensor) error
	GetSensors() ([]*WeatherSensor, error)
	GetSensorsOfUser(userId string) ([]*WeatherSensor, error)
	Close() error
}

//WeatherSensor is the data for a new Sensorregistration
type WeatherSensor struct {
	Name      string
	Id        uuid.UUID
	UserId    string
	Location  string
	Longitude float64
	Latitude  float64
}
