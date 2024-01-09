package main

import (
	"encoding/json"
	"fmt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"math/rand"
	"time"

	mqtt_helper "imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
)

func main() {
	if config, err := mqtt_helper.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml"); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config, "clientId")
		client.Connect()

		actualWind := 40.0
		min := 10.0
		max := 120.0

		for {
			actualWind = actualWind + (rand.Float64()-rand.Float64())*5
			if actualWind < min {
				actualWind = min
			}
			if actualWind > max {
				actualWind = max
			}
			data := sensor.Measurement{
				SensorID:  2,
				AirportID: "CDG",
				Type:      "Wind speed",
				Value:     actualWind,
				Timestamp: time.Now(),
			}
			jsonData, err := json.Marshal(data)

			if err != nil {
				fmt.Printf("could not marshal json: %s\n", err)
			} else {
				client.Publish("sensor/wind", 2, false, jsonData)
				fmt.Printf("%#v", data)
			}

			time.Sleep(2 * time.Second)
		}
	}
}
