package notifier

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	evt "maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
)

// SendMatrix pushes alert message to Matrix chat via webhook
func SendMatrix(event models.Event, snapshot io.Reader, provider notifMeta) {
	profile := config.ConfigData.Alerts.Matrix[provider.index]
	status := &config.Internal.Status.Notifications.Matrix[provider.index]

	var err error
	var message string
	// Build notification
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Matrix")
	} else {
		message = renderMessage("html", event, "message", "Matrix")
	}

	// New matrix client
	m, err := mautrix.NewClient(profile.Server, "", "")
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Matrix").
			Err(err).
			Int("provider_id", provider.index).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}

	// Ignore self-signed certs if set
	if profile.Insecure {
		httpClient := &http.Client{Timeout: 10 * time.Second}
		// Ignore SSL verification if set
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		m.Client = httpClient
	}

	// Handle login
	m.StateStore = mautrix.NewMemoryStateStore()
	ch, err := cryptohelper.NewCryptoHelper(m, []byte("asdf"), "./matrix.db")
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Matrix").
			Err(err).
			Int("provider_id", provider.index).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}
	ch.LoginAs = &mautrix.ReqLogin{
		Type:       mautrix.AuthTypePassword,
		Identifier: mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: profile.Username},
		Password:   profile.Password,
	}
	ch.Init(context.Background())
	m.Crypto = ch

	// Join room if needed
	_, err = m.JoinRoomByID(context.Background(), id.RoomID(profile.RoomID))
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Matrix").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}

	// Send snapshot image if available
	if event.HasSnapshot {
		var media *mautrix.RespMediaUpload
		img, _ := io.ReadAll(snapshot)
		media, _ = m.UploadMedia(context.Background(), mautrix.ReqUploadMedia{ContentBytes: img, ContentType: "image/jpeg", FileName: "snapshot.jpg"})
		_, err = m.SendMessageEvent(context.Background(), id.RoomID(profile.RoomID), evt.EventMessage, &evt.MessageEventContent{
			MsgType: evt.MsgImage,
			Body:    "snapshot.jpg",
			URL:     media.ContentURI.CUString(),
			Info:    &evt.FileInfo{MimeType: "image/jpeg"},
		})
		if err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Matrix").
				Int("provider_id", provider.index).
				Err(err).
				Msg("Unable to send alert")
			status.NotifFailure(err.Error())
			return
		}
	}

	// Send event details
	_, err = m.SendMessageEvent(context.Background(), id.RoomID(profile.RoomID), evt.EventMessage, &evt.MessageEventContent{
		MsgType:       evt.MsgText,
		Format:        "org.matrix.custom.html",
		FormattedBody: message,
	})
	if err != nil {
		log.Warn().
			Str("event_id", event.ID).
			Str("provider", "Matrix").
			Int("provider_id", provider.index).
			Err(err).
			Msg("Unable to send alert")
		status.NotifFailure(err.Error())
		return
	}

	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Matrix").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	status.NotifSuccess()
}
