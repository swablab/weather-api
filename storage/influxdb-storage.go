package storage

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

//influxStorage is the Storage implementation for InfluxDB
type influxStorage struct {
	token        string
	bucket       string
	organization string
	url          string
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
	return influx, nil
}

//Save WeatherData to InfluxDB
func (storage *influxStorage) Save(data WeatherData) error {
	tags := map[string]string{
		"location": data.Location}

	fields := map[string]interface{}{
		"temperature": data.Temperature,
		"humidity":    data.Humidity,
		"preasure":    data.Preasure}

	datapoint := influxdb2.NewPoint("new2",
		tags,
		fields,
		data.TimeStamp)

	writeAPI := storage.client.WriteAPI(storage.organization, storage.bucket)
	writeAPI.WritePoint(datapoint)
	return nil
}

//Close InfluxDB connection
func (storage *influxStorage) Close() error {
	storage.client.Close()
	return nil
}
