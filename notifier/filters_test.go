package notifier

import (
	"testing"
)

func TestCheckZones_Default(t *testing.T) {
	zones := []string{"b", "c"}

	if checkZones([]string{"a"}, zones, false, false) {
		t.Error("Expected: false, Got: true")
	}

	if !checkZones([]string{"b"}, zones, false, false) {
		t.Error("Expected: true, Got: false")
	}

	if !checkZones([]string{"a", "b"}, zones, false, false) {
		t.Error("Expected: true, Got: false")
	}
}

func TestCheckZones_MultiZone(t *testing.T) {
	zones := []string{"b", "c"}

	if checkZones([]string{"a"}, zones, true, false) {
		t.Error("Expected: false, Got: true")
	}

	if checkZones([]string{"b"}, zones, true, false) {
		t.Error("Expected: false, Got: true")
	}

	if !checkZones([]string{"c", "b", "a"}, zones, true, false) {
		t.Error("Expected: true, Got: false")
	}

	if !checkZones([]string{"a", "b", "c"}, zones, true, false) {
		t.Error("Expected: true, Got: false")
	}
}

func TestCheckZones_ZoneOrderEnforced(t *testing.T) {
	zones := []string{"b", "c"}

	if checkZones([]string{"a"}, zones, true, true) {
		t.Error("Expected: false, Got: true")
	}

	if checkZones([]string{"b"}, zones, true, true) {
		t.Error("Expected: false, Got: true")
	}

	if checkZones([]string{"c", "b", "a"}, zones, true, true) {
		t.Error("Expected: false, Got: true")
	}

	if !checkZones([]string{"a", "b", "c"}, zones, true, true) {
		t.Error("Expected: true, Got: false")
	}
}
