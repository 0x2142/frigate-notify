package config

import "github.com/0x2142/frigate-notify/models"

var Internal models.InternalConfig

func init() {
	Internal.Status.Notifications.Enabled = true
	Internal.AppVersion = "v0.4.0-dev"
	Internal.Status.Health = "n/a"
	Internal.Status.API = "n/a"
	Internal.Status.Frigate.API = "n/a"
	Internal.Status.Frigate.MQTT = "n/a"
	Internal.Status.Monitor = "n/a"
}