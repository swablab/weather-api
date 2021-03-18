package storage

import (
	"math/rand"
	"time"
)

//WeatherStorage interface for different storage-implementations of weather data
type WeatherStorage interface {
	Save(WeatherData) error
	GetData() ([]*WeatherData, error)
	Close() error
}

//WeatherData type
type WeatherData struct {
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"airPressure"`
	Temperature float64   `json:"temperature"`
	CO2Level    float64   `json:"co2level"`
	Location    string    `json:"location"`
	TimeStamp   time.Time `json:"timestamp"`
}

//NewRandomWeatherData creates random WeatherData with given Location
func NewRandomWeatherData(location string) WeatherData {
	rand.Seed(time.Now().UnixNano())
	var data WeatherData
	data.Humidity = rand.Float64() * 100
	data.Pressure = rand.Float64()*80 + 960
	data.Temperature = rand.Float64()*40 - 5
	data.CO2Level = rand.Float64()*50 + 375
	data.Location = location
	data.TimeStamp = time.Now()
	return data
}
