
[![Go](https://github.com/swablab/weather-api/actions/workflows/go.yml/badge.svg)](https://github.com/swablab/weather-api/actions/workflows/go.yml)
# weather-api

Ziel des Projektes weather-api ist es, eine einfache API für die Erhebung von Wetterdaten zu entwickeln. Kleine, selbstgebaute Wetterstationen können ihre Messergebnisse an die weather-api senden. Dort werden diese Verarbeitet, persistiert und können später wieder abgefragt werden. 


## Benötigte Fremdapplikationen
Für den Betrieb der weather-api werden einige Fremdapplikationen benötigt.

### MongoDB
In der MongoDB Datenbank werden Sensordaten der Sensoren (z.B. Name, ID, Location, ...) gespeichert.

### InfluxDB
Anfallende Wetterdaten werden in einer InfluxDB (Timeseries DBMS) gespeichert. 

### MQTT-Broker (optional)
Die Wetter-API kann Wetterdaten von Sensoren unteranderem über MQTT entgegennehmen


## Umgebungsvariablen
Key | Default-Wert  | Auswirkung
-------- | ---------- | ----------
MONGO_HOST | localhost:27017 | Hostadresse mongodb
MONGO_DB   | weathersensors  | DB-Namen mongodb
MONGO_USER | admin | Username mongodb
MONGO_PASSWORD | admin | Passwort mongodb
MONGO_COLLECTION | sensors | mongodb-Collection, in der Wettersensoren gespeichert werden
INFLUX_HOST | localhost:8086 | Hostadresse influxdb
INFLUX_TOKEN | token | Token für influxDB
INFLUX_ORG | org_name | Organisationsnamen Influx
INFLUX_BUCKET | bucket_name | Bucket-Namen, in dem die Wetterdaten abgespeichert werden
MQTT_HOST | localhost:1883 | Hostadresse MQTT-Broker
MQTT_TOPIC | sensor/# | MQTT-Topic, in welchem nach Wetterdaten geschaut wird
MQTT_USER | mqtt | Username für MQTT
MQTT_PASSWORD | mqtt | Passwort für MQTT
MQTT_PUBLISH_DELAY | 1000 | Innerhalb dieser Zeitspanne wird ein Wetterdatensatz noch durch weiter eintreffende Werte ergänzt. Danach wird der Datensatz veröffentlicht (in Millisekunden)
MQTT_ANONYMOUS | false | Anonyme Anmeldung am MQTT-Broker verwenden (ohne Username und Passwort)
ACCESS_CONTROL_ALLOW_ORIGIN_HEADER | * | CORS-Header
USE_JWT_TOKEN_VALIDATION_URL | false | Tokenvalidierung an einer URL
JWT_TOKEN_VALIDATION_URL | localhost:5000 | URL für die JWT-Token Validierung
USE_JWT_TOKEN_VALIDATION_SECRET | true | Tokenvalidierung mit der Angabe eines Secrets
JWT_TOKEN_VALIDATION_SECRET | token_Secret_value | Secret um die Signatur des JWT-Tokens zu überprüfen
ALLOW_UNREGISTERED_SENSORS | false | Wetterdaten nicht registrierter Sensoren erlauben

