package apiv1

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/0x2142/frigate-notify/config"
	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestGetVersion(t *testing.T) {
	_, api := humatest.New(t)

	Registerv1Routes(api)

	resp := api.Get("/api/v1/version")

	if resp.Code != http.StatusOK {
		t.Error("Expected HTTP 200, got ", resp.Code)
	}

	var versionResponse map[string]interface{}
	json.Unmarshal([]byte(resp.Body.Bytes()), &versionResponse)

	if versionResponse["version"] != config.Internal.AppVersion {
		t.Error("Response body did not match expected result")
	}
}
