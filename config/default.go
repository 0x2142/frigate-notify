package config

import "github.com/0x2142/frigate-notify/models"

// Default config values
var DefaultConfig Config = Config{
	App: models.App{
		Mode: "reviews",
		API: models.API{
			Enabled: false,
			Port:    8000}},
	Frigate: models.Frigate{
		Server:    "",
		Insecure:  false,
		PublicURL: "",
		Headers:   nil,
		StartupCheck: models.StartupCheck{
			Attempts: 5,
			Interval: 30,
		},
		WebAPI: models.WebAPI{
			Enabled:  false,
			Interval: 30,
			TestMode: false,
		},
		MQTT: models.MQTT{
			Enabled:     false,
			Server:      "",
			Port:        1883,
			ClientID:    "frigate-notify",
			Username:    "",
			Password:    "",
			TopicPrefix: "frigate",
		},
		Cameras: models.Cameras{
			Exclude: nil,
		},
	},
	Alerts: models.Alerts{
		General: models.General{
			Title:            "Frigate Alert",
			TimeFormat:       "",
			NoSnap:           "allow",
			SnapBbox:         false,
			SnapTimestamp:    false,
			SnapCrop:         false,
			MaxSnapRetry:     10,
			NotifyOnce:       false,
			NotifyDetections: false,
			RecheckDelay:     0,
			AudioOnly:        "allow",
		},
		Quiet: models.Quiet{
			Start: "",
			End:   "",
		},
		Zones: models.Zones{
			Unzoned: "allow",
			Allow:   nil,
			Block:   nil,
		},
		Labels: models.Labels{
			MinScore: 0,
			Allow:    nil,
			Block:    nil,
		},
		SubLabels: models.Labels{
			MinScore: 0,
			Allow:    nil,
			Block:    nil,
		},
		LicensePlate: models.LicensePlate{
			Enabled: false,
			Allow:   nil,
			Block:   nil,
		},
	},
	Monitor: models.Monitor{
		Enabled:  false,
		URL:      "",
		Interval: 60,
		Insecure: false,
	},
}
