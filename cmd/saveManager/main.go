package main

import (
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
	"time"
)

func main() {

	recorderInflux, err := storage.NewInfluxDBRecorder(
		storage.InfluxDBSettings{
			URL:          "http://localhost:8086",
			Token:        "c86LyNxa1gJSrIA0cZRq1M-chlzox9NC9WZ3jnK3-EQ372CIw0vR-rRN1dutlg7-upIJZb9dHXNvClPhYMzajg==",
			Bucket:       "meteo-airport",
			Organization: "meteo-airport",
			Measurement:  "measurements",
		})
	if err != nil {
		panic(err)
	}

	// Get the current Date(only day, month and year)
	date := time.Now().Format("2006-01-02")

	recorderCSV, err := storage.NewCSVRecorder(
		date+"-meteo-airport.csv",
		storage.CSVSettings{
			PathDirectory: "./data",
			Separator:     ';',
			TimeFormat:    time.RFC3339,
		})
	if err != nil {
		panic(err)
	}

	// List of recorders
	recorders := []storage.Recorder{
		recorderInflux,
		recorderCSV,
	}

	defer func(recorder *storage.InfluxDBRecorder) {
		err := recorder.Close()
		if err != nil {
			panic(err)
		}
	}(recorderInflux)

	measurement := sensor.Measurement{
		SensorID:  1,
		AirportID: "NTE",
		Value:     24.0,
		Unit:      "°C",
		Timestamp: time.Now(),
	}

	for _, recorder := range recorders {
		err := recorder.Record(&measurement)
		if err != nil {
			panic(err)
		}
	}
}
