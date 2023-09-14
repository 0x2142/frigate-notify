package frigate

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/0x2142/frigate-notify/notifier"
	"github.com/0x2142/frigate-notify/util"
	"golang.org/x/exp/slices"
)

const eventsURI = "/api/events"
const snapshotURI = "/snapshot.jpg"

// LastEventTime tracks the timestamp of the last event seen
var LastEventTime float64 = float64(time.Now().Unix())

// CheckForEvents queries for all detection events since last alert time
func CheckForEvents() {
	params := "?include_thumbnails=0&after=" + strconv.FormatFloat(LastEventTime, 'f', 6, 64)
	// For testing, pull 1 event immediately
	//params := "?include_thumbnails=0&limit=1"
	url := FrigateServerURL + eventsURI + params
	log.Println("Checking for new events...")

	// Query events
	response, err := util.HTTPGet(url, FrigateInsecure)
	if err != nil {
		log.Printf("Cannot get events from %s", url)
	}

	var events []Event

	json.Unmarshal([]byte(response), &events)

	log.Printf("Found %v new events.", len(events))

	for _, event := range events {
		// Convert to human-readable timestamp
		eventTime := time.Unix(int64(event.StartTime), 0)

		// Update last event check time with most recent timestamp
		if event.StartTime > LastEventTime {
			LastEventTime = event.StartTime
		}

		// Skip excluded cameras
		if slices.Contains(ExcludeCameras, event.Camera) {
			log.Printf("Skipping event from excluded camera: %v", event.Camera)
			continue
		}

		log.Printf("Event ID %v detected %v in zone(s): %v", event.ID, event.Label, event.Zones)
		log.Println("Event Start time: ", eventTime)

		// If snapshot was collected, pull down image to send with alert
		var snapshot io.Reader
		var snapshotURL string
		if event.HasSnapshot {
			snapshotURL = FrigateServerURL + eventsURI + "/" + event.ID + snapshotURI
			snapshot = GetSnapshot(snapshotURL, event.ID)
		}

		message := buildMessage(eventTime, event)

		// Send alert with snapshot
		notifier.SendAlert(message, snapshotURL, snapshot)
	}

}

// GetSnapshot downloads a snapshot from Frigate
func GetSnapshot(snapshotURL, eventID string) io.Reader {
	response, err := util.HTTPGet(snapshotURL, FrigateInsecure)
	if err != nil {
		log.Println("Could not access snaphot. Error: ", err)
	}

	return bytes.NewReader(response)
}
