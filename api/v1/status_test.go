package apiv1

import (
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestGetStatus(t *testing.T) {
	_, api := humatest.New(t)

	Registerv1Routes(api)

	resp := api.Get("/api/v1/status")

	if resp.Code != http.StatusOK {
		t.Error("Expected HTTP 200, got ", resp.Code)
	}
}
