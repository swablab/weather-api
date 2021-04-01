package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongodbSensorRegistry struct {
	weatherSensors   []*WeatherSensor
	sensorCollection *mongo.Collection
	client           *mongo.Client
}

func NewMongodbSensorRegistry(connection, database, collection string) (*mongodbSensorRegistry, error) {
	sensorRegistry := new(mongodbSensorRegistry)

	client, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		return nil, err
	}

	sensorRegistry.client = client

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	weathersensorsDB := client.Database(database)
	sensorRegistry.sensorCollection = weathersensorsDB.Collection(collection)

	return sensorRegistry, nil
}

func (registry *mongodbSensorRegistry) RegisterSensorByName(name string) (*WeatherSensor, error) {
	if registry.ExistSensorName(name) {
		return nil, fmt.Errorf("Sensorname already exists")
	}
	sensor := new(WeatherSensor)
	sensor.Name = name
	sensor.Id = uuid.New()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := registry.sensorCollection.InsertOne(ctx, sensor)

	return sensor, err
}

func (registry *mongodbSensorRegistry) ExistSensorName(name string) bool {
	sensors, err := registry.GetSensors()
	if err != nil {
		log.Fatal(err)
		return false
	}
	for _, s := range sensors {
		if s.Name == name {
			return true
		}
	}
	return false
}

func (registry *mongodbSensorRegistry) ResolveSensorById(sensorId uuid.UUID) (*WeatherSensor, bool) {
	sensors, err := registry.GetSensors()
	if err != nil {
		log.Fatal(err)
		return nil, false
	}
	for _, s := range sensors {
		if s.Id == sensorId {
			return s, true
		}
	}
	return nil, false
}

func (registry *mongodbSensorRegistry) ExistSensor(sensor *WeatherSensor) bool {
	sensors, err := registry.GetSensors()
	if err != nil {
		log.Fatal(err)
		return false
	}
	for _, s := range sensors {
		if s.Id == sensor.Id {
			return true
		}
	}
	return false
}

func (registry *mongodbSensorRegistry) GetSensors() ([]*WeatherSensor, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := registry.sensorCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var readData []*WeatherSensor
	if err = cursor.All(ctx, &readData); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return readData, nil
}

func (registry *mongodbSensorRegistry) Close() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := registry.client.Disconnect(ctx)
	return err
}
