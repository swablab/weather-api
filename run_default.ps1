#build the application
go build main.go

#set environment variables for weather-api configuration
Set-Item -Path "Env:INFLUX_HOST" -Value "localhost:8086"
Set-Item -Path "Env:INFLUX_TOKEN" -Value "token"
Set-Item -Path "Env:INFLUX_ORG" -Value "org-name"
Set-Item -Path "Env:INFLUX_BUCKET" -Value "bucket-name"

Set-Item -Path "Env:MQTT_HOST" -Value "localhost:1883"
Set-Item -Path "Env:MQTT_TOPIC" -Value "sensor/#"
Set-Item -Path "Env:MQTT_USER" -Value "mqtt"
Set-Item -Path "Env:MQTT_PASS" -Value "mqtt"
Set-Item -Path "Env:MQTT_PUBLISH_INTERVALL" -Value "2500"
Set-Item -Path "Env:MQTT_MIN_DIST_LAST_VALUE" -Value "250"
Set-Item -Path "Env:MQTT_ANONYMOUS" -Value "false"

Set-Item -Path "Env:MONGO_HOST" -Value "localhost:27017"
Set-Item -Path "Env:MONGO_DB" -Value "weathersensors"
Set-Item -Path "Env:MONGO_COLLECTION" -Value "sensors"
Set-Item -Path "Env:MONGO_USER" -Value "admin"
Set-Item -Path "Env:MONGO_PASS" -Value "admin"

Set-Item -Path "Env:ALLOW_UNREGISTERED_SENSORS" -Value "false"

#start application
Start-Process "main.exe" -Wait -NoNewWindow
