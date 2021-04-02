package storage

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

//WeatherStorage interface for different storage-implementations of weather data
type WeatherStorage interface {
	Save(WeatherData) error
	GetData(*WeatherQuery) ([]*WeatherData, error)
	Close() error
}

type SensorRegistry interface {
	RegisterSensorByName(string) (*WeatherSensor, error)
	ExistSensor(*WeatherSensor) (bool, error)
	ResolveSensorById(uuid.UUID) (*WeatherSensor, error)
	GetSensors() ([]*WeatherSensor, error)
	Close() error
}

//WeatherData type
type WeatherData struct {
	Humidity    float64   `json:"humidity"`
	Pressure    float64   `json:"pressure"`
	Temperature float64   `json:"temperature"`
	CO2Level    float64   `json:"co2level"`
	SensorId    uuid.UUID `json:"sensorId"`
	TimeStamp   time.Time `json:"timestamp"`
}

func (data *WeatherData) GetQueriedValues(query *WeatherQuery) map[string]string {
	result := map[string]string{
		"sensorId":  data.SensorId.String(),
		"timestamp": data.TimeStamp.String(),
	}
	if query.Temperature {
		result["temperature"] = strconv.FormatFloat(data.Temperature, 'f', -1, 32)
	}
	if query.Pressure {
		result["pressure"] = strconv.FormatFloat(data.Pressure, 'f', -1, 32)
	}
	if query.Co2Level {
		result["co2level"] = strconv.FormatFloat(data.CO2Level, 'f', -1, 32)
	}
	if query.Humidity {
		result["humidity"] = strconv.FormatFloat(data.Humidity, 'f', -1, 32)
	}
	return result
}

func GetOnlyQueriedFields(dataPoints []*WeatherData, query *WeatherQuery) []map[string]string {
	var result []map[string]string
	for _, data := range dataPoints {
		result = append(result, data.GetQueriedValues(query))
	}
	return result
}

//WeatherSensor is the data for a new Sensorregistration
type WeatherSensor struct {
	Name      string
	Id        uuid.UUID
	Location  string
	Longitude float64
	Latitude  float64
}

type WeatherQuery struct {
	Start       time.Time
	End         time.Time
	SensorId    uuid.UUID
	Temperature bool
	Humidity    bool
	Pressure    bool
	Co2Level    bool
}

func (data *WeatherQuery) Init() {
	data.Start = time.Now().Add(-1 * time.Hour * 24 * 14)
	data.End = time.Now()
	data.SensorId = uuid.Nil
	data.Temperature = true
	data.Humidity = true
	data.Pressure = true
	data.Co2Level = true
}

func ParseFromUrlQuery(query url.Values) (*WeatherQuery, error) {
	result := new(WeatherQuery)
	result.Init()

	start := query.Get("start")
	end := query.Get("end")
	temperature := query.Get("temperature")
	humidity := query.Get("humidity")
	pressure := query.Get("pressure")
	co2level := query.Get("co2level")

	if len(start) != 0 {
		if tval, err := time.Parse(time.RFC3339, start); err == nil {
			result.Start = tval
		} else if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	if len(end) != 0 {
		if tval, err := time.Parse(time.RFC3339, end); err == nil {
			result.End = tval
		} else if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	if bval, err := strconv.ParseBool(temperature); err == nil {
		result.Temperature = bval
	}

	if bval, err := strconv.ParseBool(humidity); err == nil {
		result.Humidity = bval
	}

	if bval, err := strconv.ParseBool(pressure); err == nil {
		result.Pressure = bval
	}

	if bval, err := strconv.ParseBool(co2level); err == nil {
		result.Co2Level = bval
	}

	return result, nil
}

//NewRandomWeatherData creates random WeatherData with given Location
func NewRandomWeatherData(sensorId uuid.UUID) WeatherData {
	rand.Seed(time.Now().UnixNano())
	var data WeatherData
	data.Humidity = rand.Float64() * 100
	data.Pressure = rand.Float64()*80 + 960
	data.Temperature = rand.Float64()*40 - 5
	data.CO2Level = rand.Float64()*50 + 375
	data.SensorId = sensorId
	data.TimeStamp = time.Now()
	return data
}
