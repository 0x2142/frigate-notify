package notifier

import (
	"sync"
	"time"

	"github.com/0x2142/frigate-notify/models"
)

// cooldownTracker reserves a notification slot atomically across concurrent alerts.
// History is kept in memory and resets when the application restarts.
type cooldownTracker struct {
	mu         sync.Mutex
	lastGlobal time.Time
	lastCamera map[string]time.Time
}

var notificationCooldown cooldownTracker

func (c *cooldownTracker) allow(camera string, settings models.General) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.allowAt(camera, settings, time.Now())
}

// allowAt requires the caller to hold mu; the explicit time also allows deterministic tests.
func (c *cooldownTracker) allowAt(camera string, settings models.General, now time.Time) bool {
	if settings.Cooldown > 0 && !c.lastGlobal.IsZero() && now.Sub(c.lastGlobal).Seconds() < float64(settings.Cooldown) {
		return false
	}
	last := c.lastCamera[camera]
	if cooldown := settings.CameraCooldown[camera]; cooldown > 0 && !last.IsZero() && now.Sub(last).Seconds() < float64(cooldown) {
		return false
	}
	if c.lastCamera == nil {
		c.lastCamera = make(map[string]time.Time)
	}
	c.lastGlobal = now
	c.lastCamera[camera] = now
	return true
}
