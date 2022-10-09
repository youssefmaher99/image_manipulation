package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestUpload(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/upload", nil)
	w := httptest.NewRecorder()
	CheckStatus(w, req)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Errorf("expected %d got %d", http.StatusOK, res.StatusCode)
	}
}

func generateUUID() uuid.UUID {
	return uuid.New()
}
