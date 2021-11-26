package storage

import (
	"context"
	"errors"
	"log"
	"weather-data/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongodbSensorRegistry struct {
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

	err = client.Connect(context.Background())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Print(err)
		return nil, err
	}

	weathersensorsDB := client.Database(mongoCfg.Database)
	sensorRegistry.sensorCollection = weathersensorsDB.Collection(mongoCfg.Collection)

	log.Print("successfully created mongodb connection")

	return sensorRegistry, nil
}

func (registry *mongodbSensorRegistry) RegisterSensor(sensor *WeatherSensor) (*WeatherSensor, error) {
	sensor.Id = uuid.New()
	_, err := registry.sensorCollection.InsertOne(context.Background(), sensor)

	return sensor, err
}

func (registry *mongodbSensorRegistry) ExistSensorName(name string) (bool, error) {
	cursor, err := registry.sensorCollection.Find(context.Background(), bson.M{"name": name})
	if err != nil {
		log.Print(err)
		return false, err
	}

	return cursor.Next(context.Background()), nil
}

func (registry *mongodbSensorRegistry) GetSensor(sensorId uuid.UUID) (*WeatherSensor, error) {
	cursor, err := registry.sensorCollection.Find(context.Background(), bson.M{"id": sensorId})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if !cursor.Next(context.Background()) {
		return nil, errors.New("sensor does not exist")
	}

	var sensor *WeatherSensor
	if err = cursor.Decode(&sensor); err != nil {
		log.Print(err)
		return nil, err
	}
	return sensor, nil
}

func (registry *mongodbSensorRegistry) ExistSensor(sensorId uuid.UUID) (bool, error) {
	cursor, err := registry.sensorCollection.Find(context.Background(), bson.M{"id": sensorId})
	if err != nil {
		log.Print(err)
		return false, err
	}

	return cursor.Next(context.Background()), nil
}

func (registry *mongodbSensorRegistry) GetSensors() ([]*WeatherSensor, error) {
	cursor, err := registry.sensorCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var readData []*WeatherSensor = make([]*WeatherSensor, 0)
	if err = cursor.All(context.Background(), &readData); err != nil {
		log.Print(err)
		return nil, err
	}

	return readData, nil
}

func (registry *mongodbSensorRegistry) GetSensorsOfUser(userId string) ([]*WeatherSensor, error) {
	cursor, err := registry.sensorCollection.Find(context.Background(), bson.M{"userid": userId})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var readData []*WeatherSensor = make([]*WeatherSensor, 0)
	if err = cursor.All(context.Background(), &readData); err != nil {
		log.Print(err)
		return nil, err
	}

	return readData, nil
}

func (registry *mongodbSensorRegistry) DeleteSensor(sensorId uuid.UUID) error {
	res, err := registry.sensorCollection.DeleteOne(context.Background(), bson.M{"id": sensorId})
	if err != nil {
		log.Print(err)
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no sensor could be deleted")
	}

	return nil
}

func (registry *mongodbSensorRegistry) UpdateSensor(sensor *WeatherSensor) error {
	res, err := registry.sensorCollection.ReplaceOne(
		context.Background(),
		bson.M{"id": sensor.Id},
		sensor)
	if err != nil {
		log.Print(err)
		return err
	}
	if res.MatchedCount == 0 || res.ModifiedCount == 0 {
		return errors.New("no sensor could be updated")
	}
	return nil
}

func (registry *mongodbSensorRegistry) Close() error {
	err := registry.client.Disconnect(context.Background())
	return err
}
