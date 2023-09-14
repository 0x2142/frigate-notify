package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	frigate "github.com/0x2142/frigate-notify/events"
	"github.com/0x2142/frigate-notify/util"
)

func main() {
	// Parse config file flag
	var configFile string
	flag.StringVar(&configFile, "c", "", "Configuration file location (default \"./config.yml\")")
	flag.Parse()
	loadConfig(configFile)

	// Set up monitor
	if ConfigData.Monitor.Enabled {
		log.Println("App monitoring enabled.")
		go func() {
			for {
				_, err := util.HTTPGet(ConfigData.Monitor.URL, ConfigData.Monitor.Insecure)
				if err != nil {
					log.Printf("Error polling monitoring URL: %v", err)
				}
				log.Println("Completed monitoring check-in.")
				time.Sleep(time.Duration(ConfigData.Monitor.Interval) * time.Second)
			}
		}()
	}

	// Loop & watch for events
	if ConfigData.Frigate.WebAPI.Enabled {
		log.Println("App Started.")
		for {
			frigate.CheckForEvents()
			time.Sleep(time.Duration(ConfigData.Frigate.WebAPI.Interval) * time.Second)
		}
	}
	// Connect MQTT
	if ConfigData.Frigate.MQTT.Enabled {
		log.Println("Connecting to MQTT Server...")
		frigate.SubscribeMQTT()
		log.Println("App running. Press Ctrl-C to quit.")
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
	}

}
