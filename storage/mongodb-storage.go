package storage

import (
	"context"
	"errors"
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
		return nil, errors.New("sensorname already exists")
	}
	sensor := new(WeatherSensor)
	sensor.Name = name
	sensor.Id = uuid.New()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = registry.sensorCollection.InsertOne(ctx, sensor)

	return sensor, err
}

func (registry *mongodbSensorRegistry) ExistSensorName(name string) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := registry.sensorCollection.Find(ctx, bson.M{"name": name})
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return cursor.Next(ctx), nil
}

func (registry *mongodbSensorRegistry) ResolveSensorById(sensorId uuid.UUID) (*WeatherSensor, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := registry.sensorCollection.Find(ctx, bson.M{"id": sensorId})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if !cursor.Next(ctx) {
		return nil, errors.New("sensor does not exist")
	}

	var sensor *WeatherSensor
	if err = cursor.Decode(&sensor); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return sensor, nil
}

func (registry *mongodbSensorRegistry) ExistSensor(sensorId uuid.UUID) (bool, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := registry.sensorCollection.Find(ctx, bson.M{"id": sensorId})
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return cursor.Next(ctx), nil
}

func (registry *mongodbSensorRegistry) GetSensors() ([]*WeatherSensor, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := registry.sensorCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var readData []*WeatherSensor = make([]*WeatherSensor, 0)
	if err = cursor.All(ctx, &readData); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return readData, nil
}

func (registry *mongodbSensorRegistry) DeleteSensor(sensorId uuid.UUID) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := registry.sensorCollection.DeleteOne(ctx, bson.M{"id": sensorId})
	if err != nil {
		log.Fatal(err)
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no sensor could be deleted")
	}

	return nil
}

func (registry *mongodbSensorRegistry) UpdateSensor(sensor *WeatherSensor) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := registry.sensorCollection.ReplaceOne(ctx,
		bson.M{"id": sensor.Id},
		sensor)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if res.MatchedCount == 0 || res.ModifiedCount == 0 {
		return errors.New("no sensor could be updated")
	}
	return nil
}

func (registry *mongodbSensorRegistry) Close() error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := registry.client.Disconnect(ctx)
	return err
}
