package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"weather-data/config"
	"weather-data/storage"
	"weather-data/weathersource"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type weatherRestApi struct {
	connection     string
	config         config.RestConfig
	weaterStorage  storage.WeatherStorage
	weatherSource  weathersource.WeatherSourceBase
	sensorRegistry storage.SensorRegistry
}

//SetupAPI sets the REST-API up
func NewRestAPI(connection string, weatherStorage storage.WeatherStorage, sensorRegistry storage.SensorRegistry, config config.RestConfig) *weatherRestApi {
	api := new(weatherRestApi)
	api.connection = connection
	api.weaterStorage = weatherStorage
	api.sensorRegistry = sensorRegistry
	api.config = config
	return api
}

//Start a new Rest-API instance
func (api *weatherRestApi) Start() error {
	router := api.handleRequests()

	originsOk := handlers.AllowedOrigins([]string{api.config.AccessControlAllowOriginHeader})

	return http.ListenAndServe(api.connection, handlers.CORS(originsOk)(router))
}

//Close the rest api
func (api *weatherRestApi) Close() {
}

func (api *weatherRestApi) handleRequests() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/", api.homePageHandler)

	//random weather data
	router.HandleFunc("/{_dummy:(?i)random}", api.randomWeatherHandler).Methods("GET")
	router.HandleFunc("/{_dummy:(?i)randomlist}", api.randomWeatherListHandler).Methods("GET")

	//sensor specific stuff
	sensorRouter := router.PathPrefix("/{_dummy:(?i)sensor}").Subrouter()
	sensorRouter.Use(api.IsAuthorized)

	sensorRouter.HandleFunc("/{id}/{_dummy:(?i)weather-data}", api.getWeatherDataHandler).Methods("GET")
	sensorRouter.HandleFunc("/{id}/{_dummy:(?i)weather-data}", api.addWeatherDataHandler).Methods("POST")

	sensorRouter.HandleFunc("", api.getAllWeatherSensorHandler).Methods("GET")
	sensorRouter.HandleFunc("/{id}", api.getWeatherSensorHandler).Methods("GET")
	sensorRouter.HandleFunc("/{id}", api.updateWeatherSensorHandler).Methods("PUT")
	sensorRouter.HandleFunc("/{id}", api.deleteWeatherSensorHandler).Methods("DELETE")

	//registration
	router.HandleFunc("/{_dummy:(?i)register/sensor}/{name}", api.registerWeatherSensorHandler).Methods("POST")

	return router
}

func (api *weatherRestApi) randomWeatherHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(storage.NewRandomWeatherData().ToMap())
}

func (api *weatherRestApi) randomWeatherListHandler(w http.ResponseWriter, r *http.Request) {
	var datapoints = make([]*storage.WeatherData, 0)

	for i := 0; i < 10; i++ {
		datapoints = append(datapoints, storage.NewRandomWeatherData())
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(storage.ToMap(datapoints))
}

func (api *weatherRestApi) getWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query, err := storage.ParseFromUrlQuery(r.URL.Query())
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	query.SensorId, err = uuid.Parse(id)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	data, err := api.weaterStorage.GetData(query)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	res := storage.ToMap(storage.GetOnlyQueriedFields(data, query))

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (api *weatherRestApi) addWeatherDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var data = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	delete(data, "sensorId")

	weatherData, err := storage.FromMap(data)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	weatherData.SensorId, err = uuid.Parse(id)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	api.addNewWeatherData(*weatherData)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(weatherData.ToMap())
}

func (api *weatherRestApi) registerWeatherSensorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	sensor, err := api.sensorRegistry.RegisterSensorByName(name)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sensor)
}

func (api *weatherRestApi) getAllWeatherSensorHandler(w http.ResponseWriter, r *http.Request) {
	var weatherSensors []*storage.WeatherSensor
	var err error

	userId := r.URL.Query().Get("userId")

	if len(userId) == 0 {
		weatherSensors, err = api.sensorRegistry.GetSensors()
	} else {
		weatherSensors, err = api.sensorRegistry.GetSensorsOfUser(userId)
	}

	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weatherSensors)
}

func (api *weatherRestApi) getWeatherSensorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sensorId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	weatherSensor, err := api.sensorRegistry.ResolveSensorById(sensorId)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weatherSensor)
}

func (api *weatherRestApi) updateWeatherSensorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sensorId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var sensor storage.WeatherSensor
	err = json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	sensor.Id = sensorId

	exist, err := api.sensorRegistry.ExistSensor(sensorId)
	if !exist || err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	err = api.sensorRegistry.UpdateSensor(&sensor)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensor)
}

func (api *weatherRestApi) deleteWeatherSensorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sensorId, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	err = api.sensorRegistry.DeleteSensor(sensorId)
	if err != nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (api *weatherRestApi) homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Weather API!")
}

func (api *weatherRestApi) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.config.UseTokenAuthorization {
			next.ServeHTTP(w, r)
			return
		}

		req, err := http.NewRequest(http.MethodGet, api.config.ValidateTokenUrl, &bytes.Buffer{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.Header = r.Header.Clone()
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if resp.StatusCode == http.StatusOK {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "", http.StatusUnauthorized)
	})
}

//AddNewWeatherDataCallback adds a new callbackMethod for incoming weather data
func (api *weatherRestApi) AddNewWeatherDataCallback(callback weathersource.NewWeatherDataCallbackFunc) {
	api.weatherSource.AddNewWeatherDataCallback(callback)
}

func (api *weatherRestApi) addNewWeatherData(weatherData storage.WeatherData) {
	api.weatherSource.NewWeatherData(weatherData)
}
