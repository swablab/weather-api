package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

//influxStorage is the Storage implementation for InfluxDB
type influxStorage struct {
	token        string
	bucket       string
	organization string
	url          string
	measurement  string
	client       influxdb2.Client
}

//NewInfluxStorage Factory
func NewInfluxStorage(token, bucket, organization, url string) (*influxStorage, error) {
	influx := new(influxStorage)
	influx.bucket = bucket
	influx.token = token
	influx.organization = organization
	influx.url = url
	influx.client = influxdb2.NewClient(url, token)
	influx.measurement = "data"
	return influx, nil
}

//Save WeatherData to InfluxDB
func (storage *influxStorage) Save(data WeatherData) error {
	tags := map[string]string{
		"sensorId": data.SensorId.String()}

	fields := map[string]interface{}{
		"temperature": data.Temperature,
		"humidity":    data.Humidity,
		"pressure":    data.Pressure,
		"co2level":    data.CO2Level}

	datapoint := influxdb2.NewPoint(storage.measurement,
		tags,
		fields,
		data.TimeStamp)

	writeAPI := storage.client.WriteAPI(storage.organization, storage.bucket)
	writeAPI.WritePoint(datapoint)
	return nil
}

//GetData datapoints from InfluxDB
func (storage *influxStorage) GetData(query *WeatherQuery) ([]*WeatherData, error) {

	fields := ""
	concat := ""

	if query.Temperature {
		fields = fmt.Sprintf("%v %v r._field == \"temperature\"", fields, concat)
		concat = "or"
	}

	if query.Humidity {
		fields = fmt.Sprintf("%v %v r._field == \"humidity\"", fields, concat)
		concat = "or"
	}

	if query.Pressure {
		fields = fmt.Sprintf("%v %v r._field == \"pressure\"", fields, concat)
		concat = "or"
	}

	if query.Co2Level {
		fields = fmt.Sprintf("%v %v r._field == \"co2level\"", fields, concat)
		concat = "or"
	}

	fields = fmt.Sprintf(" and ( %v )", fields)

	fluxQuery := fmt.Sprintf("from(bucket:\"%v\")|> range(start: %v, stop: %v) |> filter(fn: (r) => r._measurement == \"%v\" and r.sensorId == \"%v\" %v)", storage.bucket, query.Start.Format(time.RFC3339), query.End.Format(time.RFC3339), storage.measurement, query.SensorId, fields)

	res, err := storage.executeFluxQuery(fluxQuery)
	return res, err
}

func (storage *influxStorage) executeFluxQuery(query string) ([]*WeatherData, error) {

	queryAPI := storage.client.QueryAPI(storage.organization)
	result, err := queryAPI.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	var queryResults []*WeatherData

	for result.Next() {
		if result.Err() != nil {
			return nil, result.Err()
		}

		timestamp := result.Record().Time()
		sensorId, err := uuid.Parse(result.Record().ValueByKey("sensorId").(string))

		if err != nil {
			return nil, err
		}

		data, contained := containsWeatherData(queryResults, sensorId, timestamp)

		if result.Record().Field() == "temperature" {
			data.Temperature = result.Record().Value().(float64)
		}
		if result.Record().Field() == "pressure" {
			data.Pressure = result.Record().Value().(float64)
		}
		if result.Record().Field() == "humidity" {
			data.Humidity = result.Record().Value().(float64)
		}
		if result.Record().Field() == "co2level" {
			data.CO2Level = result.Record().Value().(float64)
		}

		if !contained {
			data.SensorId = sensorId
			data.TimeStamp = timestamp
			queryResults = append(queryResults, data)
		}
	}

	return queryResults, nil
}

func containsWeatherData(weatherData []*WeatherData, sensorId uuid.UUID, timestamp time.Time) (*WeatherData, bool) {
	for _, val := range weatherData {
		if val.SensorId == sensorId && val.TimeStamp == timestamp {
			return val, true
		}
	}
	var newData WeatherData
	return &newData, false
}

//Close InfluxDB connection
func (storage *influxStorage) Close() error {
	storage.client.Close()
	return nil
}
