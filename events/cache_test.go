package frigate

import (
	"testing"

	"github.com/0x2142/frigate-notify/models"
)

func TestSetZoneAlerted(t *testing.T) {
	// Setup
	InitZoneCache()
	defer CloseZoneCache()
	event := models.Event{ID: "test-event-id", CurrentZones: []string{"test_zone", "test_zone"}}

	setZoneAlerted(event)

	expected := []string{"test_zone"}
	result, ok := zoneCache.Get(event.ID)
	if !ok {
		t.Error("Could not find event ID")
	}

	// Check if zone added
	if result[0] != expected[0] {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}

	// Check if duplicates removed
	if len(result) != 1 {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
}

func TestZoneAlreadyAlerted(t *testing.T) {
	// Setup
	InitZoneCache()
	defer CloseZoneCache()
	event := models.Event{ID: "test-event-id", CurrentZones: []string{"test_zone", "test_zone"}}

	// Test new event
	result := zoneAlreadyAlerted(event)
	if result != false {
		t.Errorf("Expected: false, Got: %v", result)
	}

	// Test adding new zone to existing event
	event.CurrentZones = append(event.CurrentZones, "another_zone")
	result = zoneAlreadyAlerted(event)
	if result != false {
		t.Errorf("Expected: false, Got: %v", result)
	}

	// Test event that has already generated alert
	result = zoneAlreadyAlerted(event)
	if result != true {
		t.Errorf("Expected: true, Got: %v", result)
	}
}

func TestDelZoneAlerted(t *testing.T) {
	// Setup
	InitZoneCache()
	defer CloseZoneCache()
	event := models.Event{ID: "test-event-id", CurrentZones: []string{"test_zone", "test_zone"}}

	// Create new cache entry
	setZoneAlerted(event)

	// Test delete
	delZoneAlerted(event)
	_, ok := zoneCache.Get(event.ID)
	if ok {
		t.Errorf("Cache entry not deleted")
	}
}
