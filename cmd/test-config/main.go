package main

import (
	"imt-atlantique.project.group.fr/meteo-airport/internal/config_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"time"
)

func main() {
	config_helper.SetDefaultConfigFileName("sensor_config.yaml")
	if config, err := config_helper.LoadDefaultSensorConfig(); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config.Broker.Client, "aClientId")
		if err := client.Connect(); err != nil {
			panic(err)
		}
		defer client.Disconnect()

		for {
			measurement := sensor.Measurement{
				SensorID:  config.Sensor.SensorID,
				AirportID: config.Sensor.AirportID,
				Type:      sensor.MeasurementType(config.Sensor.Topic),
				Value:     40.0,
				Unit:      config.Sensor.Unit,
				Timestamp: time.Now(),
			}

			if err := measurement.PublishOnMQTT(1, false, client); err != nil {
				panic(err)
			}

		}
	}

}