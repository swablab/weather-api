package storage

//WeatherStorage interface for different storage-implementations of weather data
type WeatherStorage interface {
	Save(WeatherData) error
	GetData(*WeatherQuery) ([]*WeatherData, error)
	Close() error
}
