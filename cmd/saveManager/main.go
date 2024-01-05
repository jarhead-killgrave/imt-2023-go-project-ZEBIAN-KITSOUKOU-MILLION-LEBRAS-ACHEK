package main

import (
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
	"time"
)

func main() {

	recorder, err := storage.NewCSVRecorder(
		"measurements.csv",
		storage.CSVSettings{
			PathDirectory: "./data-collector",
			Separator:     ';',
			TimeFormat:    time.RFC3339,
		})
	if err != nil {
		panic(err)
	}

	measurement := storage.Measurement{
		SensorID:  1,
		AirportID: "NTE",
		Type:      "temperature",
		Value:     20.0,
		Unit:      "°C",
		Timestamp: time.Now(),
	}

	recorder.Record(&measurement)

	recorder.Close()
}
