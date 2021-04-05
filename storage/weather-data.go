package storage

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SensorValueType string

const (
	Temperature SensorValueType = "temperature"
	Pressure    SensorValueType = "pressure"
	Humidity    SensorValueType = "humidity"
	Co2Level    SensorValueType = "co2level"
)

const (
	SensorId  string = "sensorId"
	TimeStamp string = "timeStamp"
)

func GetSensorValueTypes() []SensorValueType {
	return []SensorValueType{Temperature, Pressure, Humidity, Co2Level}
}

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
	Values    map[SensorValueType]float64
	SensorId  uuid.UUID `json:"sensorId"`
	TimeStamp time.Time `json:"timestamp"`
}

func (data *WeatherData) OnlyQueriedValues(query *WeatherQuery) *WeatherData {
	for _, sensorValueType := range GetSensorValueTypes() {
		if !query.Values[sensorValueType] {
			delete(data.Values, sensorValueType)
		}
	}
	return data
}

func (data *WeatherData) ToMap() map[string]interface{} {
	mappedData := map[string]interface{}{
		"sensorId":  data.SensorId.String(),
		"timeStamp": data.TimeStamp.String(),
	}

	for sensorValueType, value := range data.Values {
		mappedData[string(sensorValueType)] = value //strconv.FormatFloat(value, 'f', -1, 64)
	}

	return mappedData
}

func FromMap(value map[string]interface{}) (*WeatherData, error) {
	var data = new(WeatherData)
	data.Values = make(map[SensorValueType]float64)
	var err error

	copy := make(map[string]interface{})
	for key, value := range value {
		copy[key] = value
	}

	idString, ok := copy[SensorId].(string)
	if !ok {
		return nil, fmt.Errorf("sensorId must be of type string")
	}
	data.SensorId, err = uuid.Parse(idString)
	if err != nil {
		return nil, err
	}
	delete(copy, SensorId)

	timeStampString, ok := copy[TimeStamp].(string)
	if !ok {
		return nil, fmt.Errorf("timeStamp must be of type string")
	}
	data.TimeStamp, err = time.Parse(time.RFC3339, timeStampString)
	if err != nil {
		return nil, err
	}
	delete(copy, TimeStamp)

	for key, val := range copy {
		switch v := val.(type) {
		case float64:
			data.Values[SensorValueType(key)] = float64(v)
		default:
		}
	}

	return data, nil
}

func GetOnlyQueriedFields(dataPoints []*WeatherData, query *WeatherQuery) []*WeatherData {
	for _, data := range dataPoints {
		data.OnlyQueriedValues(query)
	}
	return dataPoints
}

func ToMap(dataPoints []*WeatherData) []map[string]interface{} {
	var result = make([]map[string]interface{}, 0)
	for _, data := range dataPoints {
		result = append(result, data.ToMap())
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
	Start    time.Time
	End      time.Time
	SensorId uuid.UUID
	Values   map[SensorValueType]bool
}

func (query *WeatherQuery) Init() {
	query.Start = time.Now().Add(-1 * time.Hour * 24 * 14)
	query.End = time.Now()
	query.SensorId = uuid.Nil
	query.Values = make(map[SensorValueType]bool)
	for _, sensorValueType := range GetSensorValueTypes() {
		query.Values[sensorValueType] = true
	}
}

func ParseFromUrlQuery(query url.Values) (*WeatherQuery, error) {
	result := new(WeatherQuery)
	result.Init()

	start := query.Get("start")
	end := query.Get("end")

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

	for _, sensorValueType := range GetSensorValueTypes() {
		queryParam := query.Get(string(sensorValueType))
		if bval, err := strconv.ParseBool(queryParam); err == nil {
			result.Values[sensorValueType] = bval
		}
	}

	return result, nil
}

//NewRandomWeatherData creates random WeatherData with given Location
func NewRandomWeatherData(sensorId uuid.UUID) WeatherData {
	rand.Seed(time.Now().UnixNano())
	var data WeatherData
	data.Values[Humidity] = rand.Float64() * 100
	data.Values[Pressure] = rand.Float64()*80 + 960
	data.Values[Temperature] = rand.Float64()*40 - 5
	data.Values[Co2Level] = rand.Float64()*50 + 375
	data.SensorId = sensorId
	data.TimeStamp = time.Now()
	return data
}
