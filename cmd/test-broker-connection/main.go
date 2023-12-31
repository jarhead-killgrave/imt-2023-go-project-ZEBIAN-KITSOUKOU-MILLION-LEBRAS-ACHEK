package main

import (
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
)

func main() {
	if config, err := mqtt_helper.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml"); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config, "aClientId")
		if err := client.Connect(); err != nil {
			panic(err)
		}
		defer client.Disconnect()

		// Print the message when it is received
		if err := client.Subscribe("test", 1, func(client mqtt.Client, message mqtt.Message) {
			println(string(message.Payload()))
		}); err != nil {
			panic(err)
		}

		if err := client.Publish("test", 0, false, "Hello World"); err != nil {
			panic(err)
		}
		for true {

		}
	}
}
