package storage

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type WeatherQuery struct {
	Start         time.Time
	End           time.Time
	SensorId      uuid.UUID
	MaxDataPoints int
	Values        map[SensorValueType]bool
}

//NewWeatherQuery creates a new empty WeatherQuery
func NewWeatherQuery() *WeatherQuery {
	query := new(WeatherQuery)
	query.MaxDataPoints = -1
	query.Values = make(map[SensorValueType]bool)
	return query
}

func (query *WeatherQuery) Init() {
	query.Start = time.Now().Add(-1 * time.Hour * 24 * 14)
	query.End = time.Now()
	query.SensorId = uuid.Nil
	for _, sensorValueType := range GetSensorValueTypes() {
		query.Values[sensorValueType] = true
	}
}

func ParseFromUrlQuery(query url.Values) (*WeatherQuery, error) {
	result := NewWeatherQuery()
	result.Init()

	start := query.Get("start")
	end := query.Get("end")
	max := query.Get("maxDataPoints")

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

	if len(max) != 0 {
		if tval, err := strconv.Atoi(max); err == nil {
			result.MaxDataPoints = tval
		} else if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	for k, v := range query {
		if k == "start" || k == "end" {
			continue
		}
		if bval, err := strconv.ParseBool(v[0]); err == nil {
			result.Values[SensorValueType(k)] = bval
		}
	}

	return result, nil
}
