package storage

import (
	"errors"

	"github.com/google/uuid"
)

type inmemorySensorRegistry struct {
	weatherSensors []*WeatherSensor
}

func NewInmemorySensorRegistry() *inmemorySensorRegistry {
	sensorRegistry := new(inmemorySensorRegistry)
	return sensorRegistry
}

func (registry *inmemorySensorRegistry) RegisterSensorByName(name string) (*WeatherSensor, error) {
	if registry.ExistSensorName(name) {
		return nil, errors.New("Sensorname already exists")
	}
	sensor := new(WeatherSensor)
	sensor.Name = name
	sensor.Id = uuid.New()
	registry.weatherSensors = append(registry.weatherSensors, sensor)
	return sensor, nil
}

func (registry *inmemorySensorRegistry) ExistSensorName(name string) bool {
	for _, s := range registry.weatherSensors {
		if s.Name == name {
			return true
		}
	}
	return false
}

func (registry *inmemorySensorRegistry) ExistSensorId(sensorId uuid.UUID) bool {
	for _, s := range registry.weatherSensors {
		if s.Id == sensorId {
			return true
		}
	}
	return false
}

func (registry *inmemorySensorRegistry) ExistSensor(sensor *WeatherSensor) bool {
	for _, s := range registry.weatherSensors {
		if s.Id == sensor.Id {
			return true
		}
	}
	return false
}

func (registry *inmemorySensorRegistry) GetSensors() []*WeatherSensor {
	return registry.weatherSensors
}

func (registry *inmemorySensorRegistry) Close() error {
	return nil
}
