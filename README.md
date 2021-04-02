# weather-api

## Benötigte Fremdapplikationen
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
MONGO_PASS | admin | Passwort mongodb
MONGO_COLLECTION | sensors | mongodb-Collection, in der Wettersensoren gespeichert werden
INFLUX_HOST | localhost:8086 | Hostadresse influxdb
INFLUX_TOKEN | token | Token für influxDB
INFLUX_ORG | org_name | Organisationsnamen Influx
INFLUX_BUCKET | bucket_name | Bucket-Namen, in dem die Wetterdaten abgespeichert werden
MQTT_HOST | localhost:1883 | Hostadresse MQTT-Broker
MQTT_TOPIC | sensor/# | MQTT-Topic, in welchem nach Wetterdaten geschaut wird
MQTT_USER | mqtt | Username für MQTT
MQTT_PASS | mqtt | Passwort für MQTT
MQTT_PUBLISH_INTERVALL | 2500 | Intervall, nachdem über MQTT empfangene Wetterdaten in die DB geschrieben werden (in Millisekunden)
MQTT_MIN_DIST_LAST_VALUE | 250 | Zeit, die Wetterdaten mindestens zurückgehalten werden, bevor diese in die DB geschrieben werden -> Innerhalb dieser Zeitspanne kann ein Wetterdatensatz noch durch andere Werte ergänzt werden(in Millisekunden)
MQTT_ANONYMOUS | false | Anonyme Anmeldung am MQTT-Broker verwenden (ohne Username und Passwort)
ALLOW_UNREGISTERED_SENSORS | false | Wetterdaten nicht registrierter Sensoren erlauben


## Applikation unter Windows ausführen

Eine lokal ausgeführte Test-Instanz der Wetter-API muss mit URLs, Tokens und ähnlichem über Umgebungsvariablen konfiguriert werden.  
Das PowerShell-Skript `run_default.ps1` ist eine Vorlage für den start einer eigenen Instanz, lediglich die Umgebungsvariablen müsssen hierzu angepasst werden. Am besten wird der Inhalt dieses Skriptes in ein weiteres Skript (z.B. `run.ps1`) kopiert. Dieses wird von Git ignoriert, geheime Zugangsdaten (z.B. zu MQTT Broker, InfluxDB) werden so nicht ins Git-Repository eingefügt.