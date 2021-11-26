package storage

import (
	"math"
	"math/rand"
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

//WeatherData type
type WeatherData struct {
	Values    map[SensorValueType]float64
	SensorId  uuid.UUID
	TimeStamp time.Time
}

//NewRandomWeatherData creates random WeatherData
func NewRandomWeatherData() *WeatherData {
	var data = new(WeatherData)
	data.Values = make(map[SensorValueType]float64)
	data.Values[Humidity] = rand.Float64() * 100
	data.Values[Pressure] = rand.Float64()*80 + 960
	data.Values[Temperature] = rand.Float64()*40 - 5
	data.Values[Co2Level] = rand.Float64()*50 + 375
	data.SensorId = uuid.New()
	data.TimeStamp = time.Now()
	return data
}

//NewRandomWeatherData creates random WeatherData
func NewWeatherData() *WeatherData {
	var data = new(WeatherData)
	data.Values = make(map[SensorValueType]float64)
	return data
}

//OnlyQueriedValues remove all values not contained by the WeatherQuery
func (data *WeatherData) OnlyQueriedValues(query *WeatherQuery) *WeatherData {
	for sensorValueType, value := range query.Values {
		if !value {
			delete(data.Values, sensorValueType)
		}
	}
	return data
}

//ToMap converts WeatherData to a map[string]interface{}
func (data *WeatherData) ToMap() map[string]interface{} {
	mappedData := map[string]interface{}{
		SensorId:  data.SensorId.String(),
		TimeStamp: data.TimeStamp.Format(time.RFC3339),
	}

	for sensorValueType, value := range data.Values {
		mappedData[string(sensorValueType)] = value
	}

	return mappedData
}

//FromMap converts a map[string]interface{} to WeatherData
func FromMap(values map[string]interface{}) (*WeatherData, error) {
	var data = new(WeatherData)
	data.Values = make(map[SensorValueType]float64)
	var err error

	for key, val := range values {
		switch value := val.(type) {
		case float64:
			data.Values[SensorValueType(key)] = float64(value)
		case string:
			if key == SensorId {
				data.SensorId, err = uuid.Parse(value)
				if err != nil {
					return nil, err
				}
			} else if key == TimeStamp {
				data.TimeStamp, err = time.Parse(time.RFC3339, value)
				if err != nil {
					return nil, err
				}
			}
		case uuid.UUID:
			if key == SensorId {
				data.SensorId = value
			}
		case time.Time:
			if key == TimeStamp {
				data.TimeStamp = value
			}
		default:
		}
	}

	return data, nil
}

//GetOnlyQueriedFields execute onlyQueriedValues on WeatherData slice an return this
func GetOnlyQueriedFields(dataPoints []*WeatherData, query *WeatherQuery) []*WeatherData {
	for _, data := range dataPoints {
		data.OnlyQueriedValues(query)
	}

	if query.MaxDataPoints >= 0 && len(dataPoints) > query.MaxDataPoints {
		var result = make([]*WeatherData, query.MaxDataPoints)
		division := float64(len(dataPoints)) / float64(query.MaxDataPoints)
		for i := 0; i < query.MaxDataPoints; i++ {
			result[i] = dataPoints[int(math.Round(float64(i)*division))]
		}
		return result
	}

	return dataPoints
}

//ToMap mapps all WeatherData of a slice ToMap
func ToMap(dataPoints []*WeatherData) []map[string]interface{} {
	var result = make([]map[string]interface{}, 0)
	for _, data := range dataPoints {
		result = append(result, data.ToMap())
	}
	return result
}
