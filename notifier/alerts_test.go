package notifier

import (
	"os"
	"testing"

	"github.com/0x2142/frigate-notify/models"
)

func TestMain(m *testing.M) {
	TemplateFiles = os.DirFS("..")
	os.Exit(m.Run())
}

func TestRenderMessage(t *testing.T) {
	emptyEvent := models.Event{}
	fullEvent := models.Event{
		Camera:      "front_door",
		HasSnapshot: true,
		HasClip:     true,
		Zones:       []string{"driveway", "entryway"},
		Extra: models.ExtraFields{
			FormattedTime:       "2026-06-20 15:04:05 +0000",
			CameraName:          "Front Door",
			LabelList:           "person, dog",
			SubLabelList:        "visitor",
			LicensePlateList:    "ABC123",
			Audio:               "speech",
			ZoneList:            "driveway, entryway",
			PublicURL:           "https://frigate.example",
			FrigateMajorVersion: 17,
			ReviewLink:          "https://frigate.example/review/evt-1",
			EventLink:           "https://frigate.example/api/events/evt-1/clip.mp4",
		},
	}

	tests := []struct {
		name     string
		template string
		event    models.Event
		expected string
	}{
		{
			name:     "html empty event",
			template: "html",
			event:    emptyEvent,
			expected: "Detection at <br />\n" +
				"Camera: <br />\n" +
				"Links: <a href=\"/cameras/\">Camera</a>\n" +
				"<br /><br />No snapshot available.",
		},
		{
			name:     "html full event",
			template: "html",
			event:    fullEvent,
			expected: "Detection at 2026-06-20 15:04:05 +0000<br />\n" +
				"Camera: Front Door<br />\n" +
				"Label(s): person, dog<br />\n" +
				"Sublabel(s): visitor<br />\n" +
				"License Plate(s): ABC123<br />\n" +
				"Audio: speech<br />\n" +
				"Zone(s): driveway, entryway<br />\n" +
				"Links: <a href=\"https://frigate.example/#front_door\">Camera</a> | <a href=\"https://frigate.example/review/evt-1\">Review Event</a>",
		},
		{
			name:     "markdown empty event",
			template: "markdown",
			event:    emptyEvent,
			expected: "Detection at   \n" +
				"Camera: \n" +
				"Links: [Camera](/cameras/)\n" +
				"\nNo snapshot available.",
		},
		{
			name:     "markdown full event",
			template: "markdown",
			event:    fullEvent,
			expected: "Detection at 2026-06-20 15:04:05 +0000  \n" +
				"Camera: Front Door\n" +
				"Label(s): person, dog\n" +
				"Sublabel(s): visitor\n" +
				"License Plate(s): ABC123\n" +
				"Audio: speech\n" +
				"Zone(s): driveway, entryway\n" +
				"Links: [Camera](https://frigate.example/#front_door) | [Review Event](https://frigate.example/review/evt-1)",
		},
		{
			name:     "plaintext empty event",
			template: "plaintext",
			event:    emptyEvent,
			expected: "Detection at \n" +
				"Camera: \n" +
				"Links: - Camera: /cameras/\n" +
				"\nNo snapshot available.",
		},
		{
			name:     "plaintext full event",
			template: "plaintext",
			event:    fullEvent,
			expected: "Detection at 2026-06-20 15:04:05 +0000\n" +
				"Camera: Front Door\n" +
				"Label(s): person, dog\n" +
				"Sublabel(s): visitor\n" +
				"License Plate(s): ABC123\n" +
				"Audio: speech\n" +
				"Zone(s): driveway, entryway\n" +
				"Links:\n" +
				" - Camera: https://frigate.example/#front_door\n" +
				" - Review Event: https://frigate.example/review/evt-1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := renderMessage(test.template, test.event, "message", "test")
			if result != test.expected {
				t.Errorf("Expected: %q, Got: %q", test.expected, result)
			}
		})
	}
}
