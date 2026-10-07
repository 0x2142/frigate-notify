package notifier

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
)

func TestCooldownTiming(t *testing.T) {
	settings := models.General{Cooldown: 10, CameraCooldown: map[string]int{"doorbell": 30}}
	var tracker cooldownTracker
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for _, step := range []struct {
		camera string
		delay  time.Duration
		want   bool
	}{
		{"doorbell", 0, true},
		{"driveway", 10*time.Second - time.Nanosecond, false},
		{"driveway", 10 * time.Second, true},
		{"doorbell", 20 * time.Second, false},
		// Rejected alerts must not extend either cooldown.
		{"doorbell", 30 * time.Second, true},
		{"driveway", 39 * time.Second, false},
		{"driveway", 40 * time.Second, true},
	} {
		if got := tracker.allowAt(step.camera, settings, start.Add(step.delay)); got != step.want {
			t.Errorf("camera %s at %s: got %v, want %v", step.camera, step.delay, got, step.want)
		}
	}
}

func TestCameraCooldownIndependent(t *testing.T) {
	settings := models.General{CameraCooldown: map[string]int{"doorbell": 30, "garage": 0}}
	var tracker cooldownTracker
	now := time.Now()
	if !tracker.allowAt("doorbell", settings, now) || tracker.allowAt("doorbell", settings, now) {
		t.Fatal("doorbell should notify only once within its cooldown")
	}
	for _, camera := range []string{"driveway", "garage"} {
		for range 2 {
			if !tracker.allowAt(camera, settings, now) {
				t.Errorf("camera %s without a cooldown should be allowed", camera)
			}
		}
	}
}

func TestCooldownDisabled(t *testing.T) {
	var tracker cooldownTracker
	for range 2 {
		if !tracker.allowAt("doorbell", models.General{}, time.Now()) {
			t.Fatal("default configuration should not suppress alerts")
		}
	}
}

func TestCooldownConcurrent(t *testing.T) {
	for _, settings := range []models.General{
		{Cooldown: 60},
		{CameraCooldown: map[string]int{"doorbell": 60}},
	} {
		var tracker cooldownTracker
		var allowed atomic.Int32
		var workers sync.WaitGroup
		start := make(chan struct{})
		for range 50 {
			workers.Go(func() {
				<-start
				if tracker.allow("doorbell", settings) {
					allowed.Add(1)
				}
			})
		}
		close(start)
		workers.Wait()
		if got := allowed.Load(); got != 1 {
			t.Errorf("concurrent alerts: got %d accepted, want 1", got)
		}
	}
}

func TestSendAlertCooldown(t *testing.T) {
	originalConfig := config.ConfigData
	originalLast := config.Internal.Status.LastNotification
	t.Cleanup(func() {
		config.ConfigData = originalConfig
		config.Internal.Status.LastNotification = originalLast
		notificationCooldown = cooldownTracker{}
	})
	notificationCooldown = cooldownTracker{}
	config.ConfigData = config.Config{Alerts: models.Alerts{
		General: models.General{Cooldown: 60},
		Webhook: []models.Webhook{{AlertCommon: models.AlertCommon{
			Enabled: true, Filters: models.AlertFilter{Cameras: []string{"driveway"}},
		}}},
	}}
	SendAlert(nil)
	SendAlert([]models.Event{{Camera: "doorbell"}})
	if !notificationCooldown.lastGlobal.IsZero() {
		t.Fatal("empty or provider-filtered alerts must not consume a cooldown")
	}
	// Multiple providers remain eligible together for a review containing multiple detections.
	config.ConfigData.Alerts.Webhook = append(config.ConfigData.Alerts.Webhook, config.ConfigData.Alerts.Webhook[0])
	events := []models.Event{{Camera: "driveway"}, {Camera: "driveway"}}
	if got := len(eligibleAlertSenders(events)); got != 2 {
		t.Fatalf("got %d eligible providers, want 2", got)
	}
	if !notificationCooldown.allow("driveway", config.ConfigData.Alerts.General) {
		t.Fatal("filtered alerts should leave the first notification slot available")
	}
	// With an active cooldown, SendAlert must return before dispatch or status updates.
	last := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	config.Internal.Status.LastNotification = last
	SendAlert(events)
	// Audio-only reviews also use SendAlert, with a camera and audio metadata.
	SendAlert([]models.Event{{Camera: "driveway", Extra: models.ExtraFields{Audio: "bark"}}})
	if config.Internal.Status.LastNotification != last {
		t.Fatal("suppressed notifications must not update LastNotification")
	}
}
