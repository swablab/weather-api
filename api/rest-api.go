package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"weather-data/storage"
	"weather-data/weathersource"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type weatherRestApi struct {
	connection     string
	weaterStorage  storage.WeatherStorage
	weatherSource  weathersource.WeatherSourceBase
	sensorRegistry storage.SensorRegistry
}

//SetupAPI sets the REST-API up
func NewRestAPI(connection string, weatherStorage storage.WeatherStorage, sensorRegistry storage.SensorRegistry) *weatherRestApi {
	api := new(weatherRestApi)
	api.connection = connection
	api.weaterStorage = weatherStorage
	api.sensorRegistry = sensorRegistry
	return api
}

//Start a new Rest-API instance
func (api *weatherRestApi) Start() error {
	handler := api.handleRequests()
	return http.ListenAndServe(api.connection, handler)
}

//Close the rest api
func (api *weatherRestApi) Close() {
}

func (api *weatherRestApi) handleRequests() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/", api.homePageHandler)
	router.HandleFunc("/random", api.randomWeatherHandler)
	router.HandleFunc("/randomlist", api.randomWeatherListHandler)
	router.HandleFunc("/addData", api.addDataHandler)
	router.HandleFunc("/getData/{id}", api.getData)
	router.HandleFunc("/registerWeatherSensor/{name}", api.registerWeatherSensor)
	return router
}

func caselessMatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (api *weatherRestApi) getData(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	query, err := storage.ParseFromUrlQuery(r.URL.Query())
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse query: %v", err.Error()), http.StatusBadRequest)
		return
	}

	query.SensorId, err = uuid.Parse(id)
	if err != nil {
		http.Error(w, "could not parse uuid", http.StatusBadRequest)
		return
	}

	data, err := api.weaterStorage.GetData(query)
	if err != nil {
		http.Error(w, "error executing query", http.StatusBadRequest)
		return
	}

	res := storage.ToMap(storage.GetOnlyQueriedFields(data, query))
	json.NewEncoder(w).Encode(res)
}

func (api *weatherRestApi) randomWeatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(storage.NewRandomWeatherData())
}

func (api *weatherRestApi) randomWeatherListHandler(w http.ResponseWriter, r *http.Request) {
	var datapoints = make([]*storage.WeatherData, 0)

	for i := 0; i < 10; i++ {
		datapoints = append(datapoints, storage.NewRandomWeatherData())
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(datapoints)
}

func (api *weatherRestApi) addDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "only POST-Method allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("content-type", "application/json")

	var data = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	weatherData, err := storage.FromMap(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = api.addNewWeatherData(*weatherData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(weatherData.ToMap())
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
	name := vars["name"]

	sensor, err := api.sensorRegistry.RegisterSensorByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(sensor)
}

//AddNewWeatherDataCallback adds a new callbackMethod for incoming weather data
func (api *weatherRestApi) AddNewWeatherDataCallback(callback weathersource.NewWeatherDataCallbackFunc) {
	api.weatherSource.AddNewWeatherDataCallback(callback)
}

func (api *weatherRestApi) addNewWeatherData(weatherData storage.WeatherData) error {
	return api.weatherSource.NewWeatherData(weatherData)
}
