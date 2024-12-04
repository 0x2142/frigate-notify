package apiv1

import (
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestPostReload(t *testing.T) {
	_, api := humatest.New(t)

	Registerv1Routes(api)

	resp := api.Post("/api/v1/reload")

	if resp.Code != http.StatusAccepted {
		t.Error("Expected HTTP 202, got ", resp.Code)
	}
}
