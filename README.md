# weather-api

## Applikation ausführen

Eine lokal ausgeführte Test-Instanz der Wetter-API muss mit URLs, Tokens und ähnlichem über Umgebungsvariablen konfiguriert werden.  
Das PowerShell-Skript `run_default.ps1` ist eine Vorlage für den start einer eigenen Instanz, lediglich die Umgebungsvariablen müsssen hierzu angepasst werden. Am besten wird der Inhalt dieses Skriptes in ein weiteres Skript (z.B. `run.ps1`) kopiert. Dieses wird von Git ignoriert, geheime Zugangsdaten (z.B. zu MQTT Broker, InfluxDB) werden so nicht ins Git-Repository eingefügt.