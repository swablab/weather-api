package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weather-data/storage"
	"weather-data/weathersource"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type weatherRestApi struct {
	connection    string
	weaterStorage storage.WeatherStorage
	weatherSource weathersource.WeatherSourceBase
}

//SetupAPI sets the REST-API up
func NewRestAPI(connection string, weatherStorage storage.WeatherStorage) *weatherRestApi {
	api := new(weatherRestApi)
	api.connection = connection
	api.weaterStorage = weatherStorage
	return api
}

//Start a new Rest-API instance
func (api *weatherRestApi) Start() error {
	return http.ListenAndServe(api.connection, api.handleRequests())
}

func (api *weatherRestApi) Close() {
}

func (api *weatherRestApi) handleRequests() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", api.homePageHandler)
	router.HandleFunc("/random", api.randomWeatherHandler)
	router.HandleFunc("/randomlist", api.randomWeatherListHandler)
	router.HandleFunc("/addData", api.addDataHandler)
	router.HandleFunc("/getData", api.getData)
	router.HandleFunc("/registerWeatherSensor/{name}", api.registerWeatherSensor)
	return router
}

func (api *weatherRestApi) getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	data, err := api.weaterStorage.GetData()
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func (api *weatherRestApi) randomWeatherHandler(w http.ResponseWriter, r *http.Request) {
	datapoint := storage.NewRandomWeatherData("swablab")

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(datapoint)
}

func (api *weatherRestApi) randomWeatherListHandler(w http.ResponseWriter, r *http.Request) {
	var datapoints = make([]storage.WeatherData, 0)
	for i := 0; i < 10; i++ {
		datapoints = append(datapoints, storage.NewRandomWeatherData("swablab"))
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(datapoints)
}

func (api *weatherRestApi) addDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "only POST-Method allowed", http.StatusMethodNotAllowed)
		return
	}

	var data storage.WeatherData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api.addNewWeatherData(data)
}

func (api *weatherRestApi) homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Weather API!")
}

func (api *weatherRestApi) registerWeatherSensor(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "only POST-Method allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("content-type", "application/json")

	vars := mux.Vars(r)
	key := vars["name"]

	registration := SensorRegistration{
		Name: key,
		Id:   uuid.New(),
	}

	json.NewEncoder(w).Encode(registration)
}

//AddNewWeatherDataCallback adds a new callbackMethod for incoming weather data
func (api *weatherRestApi) AddNewWeatherDataCallback(callback weathersource.NewWeatherDataCallbackFunc) {
	api.weatherSource.AddNewWeatherDataCallback(callback)
}

func (api *weatherRestApi) addNewWeatherData(weatherData storage.WeatherData) {
	api.weatherSource.NewWeatherData(weatherData)
	api.weaterStorage.Save(weatherData)
}
