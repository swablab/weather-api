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
	exist, err := registry.ExistSensorName(name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("sensorname already exists")
	}
	sensor := new(WeatherSensor)
	sensor.Name = name
	sensor.Id = uuid.New()
	registry.weatherSensors = append(registry.weatherSensors, sensor)
	return sensor, nil
}

func (registry *inmemorySensorRegistry) ExistSensorName(name string) (bool, error) {
	for _, s := range registry.weatherSensors {
		if s.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (registry *inmemorySensorRegistry) ResolveSensorById(sensorId uuid.UUID) (*WeatherSensor, error) {
	for _, s := range registry.weatherSensors {
		if s.Id == sensorId {
			return s, nil
		}
	}
	return nil, errors.New("sensor does not exist")
}

func (registry *inmemorySensorRegistry) ExistSensor(sensorId uuid.UUID) (bool, error) {
	for _, s := range registry.weatherSensors {
		if s.Id == sensorId {
			return true, nil
		}
	}
	return false, nil
}

func (registry *inmemorySensorRegistry) DeleteSensor(sensorId uuid.UUID) error {
	for i, s := range registry.weatherSensors {
		if s.Id == sensorId {
			registry.weatherSensors = remove(registry.weatherSensors, i)
			return nil
		}
	}
	return nil
}

func (registry *inmemorySensorRegistry) UpdateSensor(sensor *WeatherSensor) error {
	for i, s := range registry.weatherSensors {
		if s.Id == sensor.Id {
			registry.weatherSensors[i] = sensor
		}
	}

	return nil
}

func (registry *inmemorySensorRegistry) GetSensors() ([]*WeatherSensor, error) {
	return registry.weatherSensors, nil
}

func (registry *inmemorySensorRegistry) Close() error {
	return nil
}

func remove(s []*WeatherSensor, i int) []*WeatherSensor {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
