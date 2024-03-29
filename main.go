package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/0x2142/frigate-notify/config"
	frigate "github.com/0x2142/frigate-notify/events"
	"github.com/0x2142/frigate-notify/util"
)

var APP_VER = "v0.2.5"

func main() {
	log.Println("Frigate Notify -", APP_VER)
	log.Println("Starting...")
	// Parse config file flag
	var configFile string
	flag.StringVar(&configFile, "c", "", "Configuration file location (default \"./config.yml\")")
	flag.Parse()
	config.LoadConfig(configFile)

	// Set up monitor
	if config.ConfigData.Monitor.Enabled {
		log.Println("App monitoring enabled.")
		go func() {
			for {
				_, err := util.HTTPGet(config.ConfigData.Monitor.URL, config.ConfigData.Monitor.Insecure)
				if err != nil {
					log.Printf("Error polling monitoring URL: %v", err)
				}
				log.Println("Completed monitoring check-in.")
				time.Sleep(time.Duration(config.ConfigData.Monitor.Interval) * time.Second)
			}
		}()
	}

	// Loop & watch for events
	if config.ConfigData.Frigate.WebAPI.Enabled {
		log.Println("App running. Press Ctrl-C to quit.")
		for {
			frigate.CheckForEvents()
			time.Sleep(time.Duration(config.ConfigData.Frigate.WebAPI.Interval) * time.Second)
		}
	}
	// Connect MQTT
	if config.ConfigData.Frigate.MQTT.Enabled {
		log.Println("Connecting to MQTT Server...")
		frigate.SubscribeMQTT()
		log.Println("App running. Press Ctrl-C to quit.")
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
	}

}
