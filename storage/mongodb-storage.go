package storage

import (
	"context"
	"fmt"
	"log"
	"time"
	"weather-data/config"

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

func NewMongodbSensorRegistry(mongoCfg config.MongoConfig) (*mongodbSensorRegistry, error) {
	sensorRegistry := new(mongodbSensorRegistry)

	options := options.Client().ApplyURI(mongoCfg.Host).SetAuth(options.Credential{Username: mongoCfg.Username, Password: mongoCfg.Password})

	client, err := mongo.NewClient(options)
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

	weathersensorsDB := client.Database(mongoCfg.Database)
	sensorRegistry.sensorCollection = weathersensorsDB.Collection(mongoCfg.Collection)

	log.Print("successfully created mongodb connection")

	return sensorRegistry, nil
}

func (registry *mongodbSensorRegistry) RegisterSensorByName(name string) (*WeatherSensor, error) {
	exist, err := registry.ExistSensorName(name)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("Sensorname already exists")
	}
	sensor := new(WeatherSensor)
	sensor.Name = name
	sensor.Id = uuid.New()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = registry.sensorCollection.InsertOne(ctx, sensor)

	return sensor, err
}

func (registry *mongodbSensorRegistry) ExistSensorName(name string) (bool, error) {
	sensors, err := registry.GetSensors()
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	for _, s := range sensors {
		if s.Name == name {
			return true, nil
		}
	}
	return false, nil
}

func (registry *mongodbSensorRegistry) ResolveSensorById(sensorId uuid.UUID) (*WeatherSensor, error) {
	sensors, err := registry.GetSensors()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for _, s := range sensors {
		if s.Id == sensorId {
			return s, nil
		}
	}
	return nil, fmt.Errorf("sensor does not exist")
}

func (registry *mongodbSensorRegistry) ExistSensor(sensor *WeatherSensor) (bool, error) {
	sensors, err := registry.GetSensors()
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	for _, s := range sensors {
		if s.Id == sensor.Id {
			return true, nil
		}
	}
	return false, nil
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