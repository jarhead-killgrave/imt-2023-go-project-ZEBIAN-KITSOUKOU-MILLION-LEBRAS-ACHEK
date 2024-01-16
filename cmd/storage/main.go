package main

import (
	"os"

	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
)

func createManager(storageConfig *config.Storage) *storage.Manager {
	client := mqtt.NewClient(&storageConfig.MQTT)

	if connexionErr := client.Connect(); connexionErr != nil {
		log.Error("Error connecting to MQTT broker: %v", connexionErr)
		os.Exit(1)
	}

	manager := storage.NewManager(client)

	for measurement, setting := range storageConfig.Settings {
		if setting.InfluxDB != (config.InfluxDBSettings{}) {
			log.Info("Registering InfluxDB recorder for measurement %s", measurement)
			log.Info("InfluxDB settings: %v", setting.InfluxDB)
			influxDBRecorder, _ := storage.NewInfluxDBRecorder(setting.InfluxDB)
			manager.AddRecorder(sensor.MeasurementType(measurement), influxDBRecorder, setting.Qos)
		}

		if setting.CSV != (config.CSVSettings{}) {
			log.Info("Registering CSV recorder for measurement %s", measurement)

			csvRecorder, _ := storage.NewCSVRecorder(setting.CSV)
			manager.AddRecorder(sensor.MeasurementType(measurement), csvRecorder, 0)

			if _, err := os.Stat(setting.CSV.PathDirectory); os.IsNotExist(err) {
				log.Info("Creating the directory for saving the CSV files...")

				err := os.Mkdir(setting.CSV.PathDirectory, 0755)

				if err != nil {
					log.Error("Error creating the directory for saving the CSV files: %v", err)
					os.Exit(1)
				}
			}
		}
	}

	return manager
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Warn("No config file specified, using default path: config/config.yaml")
	} else {
		config.SetDefaultConfigFileName(args[0])
	}

	log.Info("Loading Configurations of the storage manager...")

	defaultStorageConfig, configErr := config.LoadDefaultStorageConfig()

	if configErr != nil {
		log.Error("Error loading defaultStorageConfig: %v", configErr)
		os.Exit(1)
	}

	log.Info("Starting storage manager...")

	manager := createManager(defaultStorageConfig)

	if err := manager.Start(); err != nil {
		log.Error("Error starting storage manager: %v", err)
		os.Exit(1)
	}

	defer func(manager *storage.Manager) {
		err := manager.Close()
		if err != nil {
			log.Error("Error closing storage manager: %v", err)
		}
	}(manager)

	select {}
}
