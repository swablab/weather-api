package storage

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
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
	influx.measurement = "weather-data"
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

	for sensorValueType, value := range query.Values {
		if value {
			fields = fmt.Sprintf("%v %v r[\"_field\"] == \"%v\"", fields, concat, string(sensorValueType))
			concat = "or"
		}
	}

	fromTemplate := fmt.Sprintf("from(bucket:\"%v\")", storage.config.Bucket)
	rangeTemplate := fmt.Sprintf("|> range(start: %v, stop: %v)", query.Start.Format(time.RFC3339), query.End.Format(time.RFC3339))
	measurementTemplate := fmt.Sprintf("|> filter(fn: (r) => r[\"_measurement\"] == \"%v\")", storage.measurement)
	sensorIdsTemplate := fmt.Sprintf("|> filter(fn: (r) => r[\"sensorId\"] == \"%v\")", query.SensorId)
	fields = fmt.Sprintf("|> filter(fn: (r) => %v )", strings.Trim(fields, " "))

	fluxQuery := fmt.Sprintf("%v \n %v \n %v \n %v \n %v", fromTemplate, rangeTemplate, measurementTemplate, sensorIdsTemplate, fields)

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

		data.Values[SensorValueType(result.Record().Field())] = result.Record().Value().(float64)

		if !contained {
			data.SensorId = sensorId
			data.TimeStamp = timestamp
			queryResults = append(queryResults, data)
		}
	}

	//if some attibutes missed in a few datapoints they are not ordered by time
	//influx query e.g. first all humidity, than pressure and temperature at last. if there are some datapoints with only pressore and/or temperature they are the last inserted in the array
	sort.Slice(queryResults, func(p, q int) bool {
		return queryResults[p].TimeStamp.Before(queryResults[q].TimeStamp)
	})

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
