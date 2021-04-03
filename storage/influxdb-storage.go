package storage

import (
	"context"
	"fmt"
	"log"
	"time"
	"weather-data/config"

	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

//influxStorage is the Storage implementation for InfluxDB
type influxStorage struct {
	config      config.InfluxConfig
	measurement string
	client      influxdb2.Client
}

//NewInfluxStorage Factory
func NewInfluxStorage(cfg config.InfluxConfig) (*influxStorage, error) {
	influx := new(influxStorage)
	influx.config = cfg
	influx.client = influxdb2.NewClient(cfg.Host, cfg.Token)
	influx.measurement = "data"
	log.Print("Successfully created influx-client")
	return influx, nil
}

//Save WeatherData to InfluxDB
func (storage *influxStorage) Save(data WeatherData) error {
	tags := map[string]string{
		"sensorId": data.SensorId.String()}

	fields := make(map[string]interface{})

	for k, v := range data.Values {
		fields[string(k)] = v
	}

	datapoint := influxdb2.NewPoint(storage.measurement,
		tags,
		fields,
		data.TimeStamp)

	writeAPI := storage.client.WriteAPI(storage.config.Organization, storage.config.Bucket)
	writeAPI.WritePoint(datapoint)
	return nil
}

//GetData datapoints from InfluxDB
func (storage *influxStorage) GetData(query *WeatherQuery) ([]*WeatherData, error) {
	fluxQuery := storage.createFluxQuery(query)
	res, err := storage.executeFluxQuery(fluxQuery)
	return res, err
}

func (storage *influxStorage) createFluxQuery(query *WeatherQuery) string {
	fields := ""
	concat := ""

	for _, sensorValueType := range GetSensorValueTypes() {
		if query.Values[sensorValueType] {
			fields = fmt.Sprintf("%v %v r._field == \"%v\"", fields, concat, string(sensorValueType))
			concat = "or"
		}
	}

	fields = fmt.Sprintf(" and ( %v )", fields)

	fluxQuery := fmt.Sprintf("from(bucket:\"%v\")|> range(start: %v, stop: %v) |> filter(fn: (r) => r._measurement == \"%v\" and r.sensorId == \"%v\" %v)", storage.config.Bucket, query.Start.Format(time.RFC3339), query.End.Format(time.RFC3339), storage.measurement, query.SensorId, fields)
	return fluxQuery
}

func (storage *influxStorage) executeFluxQuery(query string) ([]*WeatherData, error) {

	queryAPI := storage.client.QueryAPI(storage.config.Organization)
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

		for _, sensorValueType := range GetSensorValueTypes() {
			if result.Record().Field() == string(sensorValueType) {
				data.Values[sensorValueType] = result.Record().Value().(float64)
			}
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
	newData.Values = make(map[SensorValueType]float64)
	return &newData, false
}

//Close InfluxDB connection
func (storage *influxStorage) Close() error {
	storage.client.Close()
	return nil
}
