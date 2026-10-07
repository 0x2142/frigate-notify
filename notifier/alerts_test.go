package notifier

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestSendAlertFiltersAfterSnapshot(t *testing.T) {
	originalConfig := config.ConfigData
	originalLast := config.Internal.Status.LastNotification
	originalWebhookStatus := config.Internal.Status.Notifications.Webhook
	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		config.ConfigData = originalConfig
		config.Internal.Status.LastNotification = originalLast
		config.Internal.Status.Notifications.Webhook = originalWebhookStatus
		http.DefaultTransport = originalTransport
		notificationCooldown = cooldownTracker{}
	})
	config.ConfigData = config.Config{
		Frigate: models.Frigate{Server: "http://frigate.test"},
		Alerts: models.Alerts{
			General: models.General{Cooldown: 60, MaxSnapRetry: 1},
			Webhook: []models.Webhook{{
				AlertCommon: models.AlertCommon{Enabled: true},
				Server:      "http://webhook.test", Method: "POST",
			}},
		},
	}
	config.Internal.Status.Notifications.Webhook = make([]models.NotifierStatus, 1)
	config.Internal.Status.LastNotification = time.Time{}
	notificationCooldown = cooldownTracker{}
	snapshotFetched := false
	sent := make(chan struct{}, 1)
	http.DefaultTransport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "/snapshot.jpg") {
			snapshotFetched = true
			// Activate provider quiet hours during the fetch, without waiting for a clock boundary.
			now := time.Now()
			config.ConfigData.Alerts.Webhook[0].Filters.Quiet = models.Quiet{
				Start: now.Add(-time.Minute).Format("15:04"),
				End:   now.Add(time.Minute).Format("15:04"),
			}
		} else {
			sent <- struct{}{}
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("snapshot")), Header: make(http.Header)}, nil
	})
	SendAlert([]models.Event{{Camera: "doorbell", ID: "quiet-test", HasSnapshot: true}})
	if !snapshotFetched {
		t.Fatal("snapshot was not fetched")
	}
	if !config.Internal.Status.LastNotification.IsZero() {
		// Wait for the unexpected sender to capture its configuration before cleanup.
		select {
		case <-sent:
		case <-time.After(5 * time.Second):
		}
		t.Fatal("notification dispatched despite quiet hours starting during snapshot retrieval")
	}
	if !notificationCooldown.lastGlobal.IsZero() {
		t.Fatal("provider-filtered alert must not consume a cooldown")
	}
}

func TestSendTestAlertCooldown(t *testing.T) {
	for _, active := range []bool{false, true} {
		name := "fresh cooldown"
		if active {
			name = "active cooldown"
		}
		t.Run(name, func(t *testing.T) {
			originalConfig := config.ConfigData
			originalLast := config.Internal.Status.LastNotification
			originalWebhookStatus := config.Internal.Status.Notifications.Webhook
			t.Cleanup(func() {
				config.ConfigData = originalConfig
				config.Internal.Status.LastNotification = originalLast
				config.Internal.Status.Notifications.Webhook = originalWebhookStatus
				notificationCooldown = cooldownTracker{}
			})

			requests := make(chan WebhookPayload, 1)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var payload WebhookPayload
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Errorf("invalid webhook payload: %v", err)
				}
				w.WriteHeader(http.StatusOK)
				requests <- payload
			}))
			defer server.Close()

			settings := models.General{Cooldown: 60, CameraCooldown: map[string]int{"doorbell": 120}}
			config.ConfigData = config.Config{Alerts: models.Alerts{
				General: settings,
				Webhook: []models.Webhook{{
					AlertCommon: models.AlertCommon{
						Enabled: true,
						Filters: models.AlertFilter{Cameras: []string{"doorbell"}},
					},
					Server: server.URL,
					Method: "POST",
				}},
			}}
			config.Internal.Status.Notifications.Webhook = make([]models.NotifierStatus, 1)
			notificationCooldown = cooldownTracker{}
			if active {
				notificationCooldown.allow("doorbell", settings)
			}
			lastGlobal := notificationCooldown.lastGlobal
			lastCamera := notificationCooldown.lastCamera["doorbell"]

			// Manual tests still obey provider filters and safely handle an empty response.
			sentinel := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			config.Internal.Status.LastNotification = sentinel
			SendTestAlert(nil)
			SendTestAlert([]models.Event{{Camera: "driveway", ID: "filtered"}})
			if config.Internal.Status.LastNotification != sentinel {
				t.Fatal("empty or provider-filtered tests must not dispatch")
			}

			SendTestAlert([]models.Event{{Camera: "doorbell", ID: "manual-test"}})
			select {
			case payload := <-requests:
				if payload.ID != "manual-test" {
					t.Fatalf("unexpected dispatched event: %s", payload.ID)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("manual test was not delivered")
			}
			if notificationCooldown.lastGlobal != lastGlobal || notificationCooldown.lastCamera["doorbell"] != lastCamera {
				t.Fatal("manual tests must not start or extend either cooldown")
			}
			if got := notificationCooldown.allow("doorbell", settings); got == active {
				t.Fatalf("next real alert allowed=%v, want %v", got, !active)
			}
		})
	}
}
