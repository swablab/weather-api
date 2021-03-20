#build the application
go build main.go

#set environment variables for weather-api configuration
Set-Item -Path "Env:WEATHER-API-INFLUX_URL" -Value "https://influx.default-address.com"
Set-Item -Path "Env:WEATHER-API-INFLUX_TOKEN" -Value "default-token"
Set-Item -Path "Env:WEATHER-API-INFLUX_ORG" -Value "default-org"
Set-Item -Path "Env:WEATHER-API-INFLUX_BUCKET" -Value  "default-bucket"

Set-Item -Path "Env:WEATHER-API-MQTT_URL" -Value "tcp://default-address.com:1883"
Set-Item -Path "Env:WEATHER-API-MQTT_TOPIC" -Value "sensor/#"

Set-Item -Path "Env:WEATHER-API-ALLOW_UNREGISTERED_SENSORS" -Value "true"

#start application
Start-Process "main.exe" -Wait -NoNewWindow
