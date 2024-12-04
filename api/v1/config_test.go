package apiv1

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestGetConfig(t *testing.T) {
	_, api := humatest.New(t)

	Registerv1Routes(api)

	resp := api.Get("/api/v1/config")

	if resp.Code != http.StatusOK {
		t.Error("Expected HTTP 200, got ", resp.Code)
	}
}

func TestPutConfig(t *testing.T) {
	_, api := humatest.New(t)

	Registerv1Routes(api)

	newconfig := `{
   "config":{
      "frigate":{
         "server":"http://192.0.2.10:5000",
         "mqtt":{
            "enabled": true
         }
      },
      "alerts":{
      }
   },
	  "skipvalidate": true,
	  "skipbackup": true,
	  "skipsave": true,
	  "skipreload": true
}`

	resp := api.Put("/api/v1/config", bytes.NewReader([]byte(newconfig)))

	if resp.Code != http.StatusAccepted {
		t.Error("Expected HTTP 202, got ", resp.Code)
	}
}
