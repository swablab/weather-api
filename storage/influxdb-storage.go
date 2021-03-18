package storage

import (
	"context"
	"fmt"
	"time"

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
		"location": data.Location}

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
func (storage *influxStorage) GetData() ([]*WeatherData, error) {

	query := fmt.Sprintf("from(bucket:\"%v\")|> range(start: -40m, stop: -20m) |> filter(fn: (r) => r._measurement == \"data\" and r.location == \"Hamburg\")", storage.bucket)

	res, err := storage.executeFluxQuery(query)
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
			return nil, err
		}
		location := result.Record().ValueByKey("location").(string)
		timestamp := result.Record().Time()

		data, contained := containsWeatherData(queryResults, location, timestamp)

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
			data.Location = location
			data.TimeStamp = timestamp
			queryResults = append(queryResults, data)
		}
	}

	return queryResults, nil
}

func containsWeatherData(weatherData []*WeatherData, location string, timestamp time.Time) (*WeatherData, bool) {
	for _, val := range weatherData {
		if val.Location == location && val.TimeStamp == timestamp {
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
