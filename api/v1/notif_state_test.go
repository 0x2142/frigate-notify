package apiv1

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestGetNotifState(t *testing.T) {
	_, api := humatest.New(t)

	Registerv1Routes(api)

	resp := api.Get("/api/v1/notif_state")

	if resp.Code != http.StatusOK {
		t.Error("Expected HTTP 200, got ", resp.Code)
	}
}

func TestPostNotifState(t *testing.T) {
	_, api := humatest.New(t)

	Registerv1Routes(api)

	resp := api.Post("/api/v1/notif_state", bytes.NewReader([]byte(`{"enabled": false}`)))

	if resp.Code != http.StatusAccepted {
		t.Error("Expected HTTP 202, got ", resp.Code)
	}
}
