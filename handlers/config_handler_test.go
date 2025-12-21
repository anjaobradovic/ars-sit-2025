package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateConfigHandler_InvalidJSON(t *testing.T) {
	body := []byte(`{ invalid json `)

	req := httptest.NewRequest(http.MethodPost, "/configs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := NewConfigHandler(nil)

	handler.CreateConfig(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}
