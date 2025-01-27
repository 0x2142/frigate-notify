package apiv1

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestGetHealthz(t *testing.T) {
	_, api := humatest.New(t)

	Registerv1Routes(api)

	resp := api.Get("/api/v1/healthz")

	if resp.Code != http.StatusOK {
		t.Error("Expected HTTP 200, got ", resp.Code)
	}

	var healthzResponse map[string]interface{}
	json.Unmarshal([]byte(resp.Body.Bytes()), &healthzResponse)

	if healthzResponse["status"] != "ok" {
		t.Error("Response body did not match expected result")
	}
}
