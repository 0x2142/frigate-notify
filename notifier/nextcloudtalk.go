package notifier

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/0x2142/frigate-notify/config"
	"github.com/0x2142/frigate-notify/models"
	"github.com/0x2142/frigate-notify/util"
)

type ncTalkClient struct {
	server     string
	username   string
	password   string
	roomToken  string
	uploadPath string
	insecure   bool
	http       *http.Client
}

type ocsMeta struct {
	Status     string `json:"status"`
	Statuscode int    `json:"statuscode"`
	Message    string `json:"message"`
}

type ocsEnvelope struct {
	OCS struct {
		Meta ocsMeta `json:"meta"`
	} `json:"ocs"`
}

func newNcTalkClient(profile models.NextcloudTalk) *ncTalkClient {
	server := strings.TrimRight(profile.Server, "/")
	uploadPath := profile.UploadPath
	if uploadPath == "" {
		uploadPath = "/frigate-notify"
	}
	if !strings.HasPrefix(uploadPath, "/") {
		uploadPath = "/" + uploadPath
	}
	uploadPath = strings.TrimRight(uploadPath, "/")

	client := &http.Client{Timeout: time.Duration(util.HTTPTimeout) * time.Second}
	if profile.Insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &ncTalkClient{
		server:     server,
		username:   profile.Username,
		password:   profile.Password,
		roomToken:  profile.RoomToken,
		uploadPath: uploadPath,
		insecure:   profile.Insecure,
		http:       client,
	}
}

func (c *ncTalkClient) authHeader() string {
	creds := c.username + ":" + c.password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
}

func (c *ncTalkClient) webdavURL(filePath string) string {
	escapedUser := url.PathEscape(c.username)
	return c.server + "/remote.php/dav/files/" + escapedUser + filePath
}

func (c *ncTalkClient) doRequest(method, reqURL string, body io.Reader, contentType string) ([]byte, int, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = body
	} else {
		bodyReader = http.NoBody
	}

	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", c.authHeader())
	req.Header.Set("OCS-APIRequest", "true")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", util.AppUserAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return respBody, resp.StatusCode, nil
}

func parseOCS(body []byte) error {
	var envelope ocsEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("invalid OCS response: %w", err)
	}
	if envelope.OCS.Meta.Statuscode >= 200 && envelope.OCS.Meta.Statuscode < 300 {
		return nil
	}
	if envelope.OCS.Meta.Message != "" {
		return fmt.Errorf("OCS error %d: %s", envelope.OCS.Meta.Statuscode, envelope.OCS.Meta.Message)
	}
	return fmt.Errorf("OCS error %d", envelope.OCS.Meta.Statuscode)
}

func (c *ncTalkClient) ensureUploadDir() error {
	dirURL := c.webdavURL(c.uploadPath)
	_, status, err := c.doRequest("MKCOL", dirURL, nil, "")
	if err != nil {
		return err
	}
	if status == http.StatusCreated || status == http.StatusOK || status == http.StatusNoContent {
		return nil
	}
	if status == http.StatusMethodNotAllowed || status == http.StatusConflict || status == http.StatusForbidden {
		// Directory likely already exists
		return nil
	}
	return fmt.Errorf("failed to create upload directory, status %d", status)
}

func (c *ncTalkClient) uploadFile(filePath string, data []byte, contentType string) error {
	if err := c.ensureUploadDir(); err != nil {
		return err
	}

	fileURL := c.webdavURL(filePath)
	req, err := http.NewRequest(http.MethodPut, fileURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", c.authHeader())
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", util.AppUserAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("WebDAV upload failed with status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *ncTalkClient) deleteFile(filePath string) error {
	fileURL := c.webdavURL(filePath)
	_, status, err := c.doRequest(http.MethodDelete, fileURL, nil, "")
	if err != nil {
		return err
	}
	if status == http.StatusNotFound || status == http.StatusNoContent || (status >= 200 && status < 300) {
		return nil
	}
	return fmt.Errorf("WebDAV delete failed with status %d", status)
}

func (c *ncTalkClient) shareFileToTalk(filePath, caption string) error {
	talkMeta := map[string]string{}
	if caption != "" {
		talkMeta["caption"] = caption
	}
	metaJSON, err := json.Marshal(talkMeta)
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Set("shareType", "10")
	form.Set("shareWith", c.roomToken)
	form.Set("path", filePath)
	form.Set("talkMetaData", string(metaJSON))

	shareURL := c.server + "/ocs/v2.php/apps/files_sharing/api/v1/shares"
	body, status, err := c.doRequest(http.MethodPost, shareURL, strings.NewReader(form.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	if status < 200 || status > 299 {
		return fmt.Errorf("share request failed with status %d: %s", status, string(body))
	}
	return parseOCS(body)
}

func (c *ncTalkClient) sendTextMessage(message string) error {
	chatURL := c.server + "/ocs/v2.php/apps/spreed/api/v1/chat/" + url.PathEscape(c.roomToken)
	form := url.Values{}
	form.Set("message", message)

	body, status, err := c.doRequest(http.MethodPost, chatURL, strings.NewReader(form.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	if status < 200 || status > 299 {
		return fmt.Errorf("chat request failed with status %d: %s", status, string(body))
	}
	return parseOCS(body)
}

// SendNextcloudTalkMessage sends alert through Nextcloud Talk
func SendNextcloudTalkMessage(event models.Event, snapshot io.Reader, provider notifMeta) {
	profile := config.ConfigData.Alerts.NextcloudTalk[provider.index]
	status := &config.Internal.Status.Notifications.NextcloudTalk[provider.index]

	var message string
	if profile.Template != "" {
		message = renderMessage(profile.Template, event, "message", "Nextcloud Talk")
	} else {
		message = renderMessage("nextcloudtalk", event, "message", "Nextcloud Talk")
	}

	client := newNcTalkClient(profile)

	var media []byte
	var filePath string
	var contentType string

	if event.HasClip && profile.SendClip {
		clip := GetClip(event)
		if clip != nil {
			media, _ = io.ReadAll(clip)
		}
		if len(media) == 0 {
			event.HasClip = false
		} else {
			filePath = path.Join(client.uploadPath, event.ID+".mp4")
			contentType = "video/mp4"
		}
	}

	if len(media) == 0 && event.HasSnapshot && snapshot != nil {
		media, _ = io.ReadAll(snapshot)
		if len(media) > 0 {
			filePath = path.Join(client.uploadPath, event.ID+".jpg")
			contentType = http.DetectContentType(media)
		}
	}

	if len(media) > 0 {
		if err := client.uploadFile(filePath, media, contentType); err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Nextcloud Talk").
				Int("provider_id", provider.index).
				Err(err).
				Msg("Unable to upload media for alert")
			status.NotifFailure(err.Error())
			return
		}

		if err := client.shareFileToTalk(filePath, message); err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Nextcloud Talk").
				Int("provider_id", provider.index).
				Err(err).
				Msg("Unable to share media to Talk room")
			status.NotifFailure(err.Error())
			if !profile.KeepStagedFiles {
				_ = client.deleteFile(filePath)
			}
			return
		}

		if !profile.KeepStagedFiles {
			if err := client.deleteFile(filePath); err != nil {
				log.Debug().
					Str("event_id", event.ID).
					Str("provider", "Nextcloud Talk").
					Int("provider_id", provider.index).
					Err(err).
					Msg("Uploaded media shared but cleanup failed")
			}
		}
	} else {
		if err := client.sendTextMessage(message); err != nil {
			log.Warn().
				Str("event_id", event.ID).
				Str("provider", "Nextcloud Talk").
				Int("provider_id", provider.index).
				Err(err).
				Msg("Unable to send alert")
			status.NotifFailure(err.Error())
			return
		}
	}

	log.Info().
		Str("event_id", event.ID).
		Str("provider", "Nextcloud Talk").
		Int("provider_id", provider.index).
		Msg("Alert sent")
	status.NotifSuccess()
}
