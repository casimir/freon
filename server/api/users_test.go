package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casimir/freon/auth"
	"github.com/gin-gonic/gin"
)

func TestMe(t *testing.T) {
	testUser, err := auth.CreateUser("test", "test", false)
	if err != nil {
		t.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router.Group("", auth.HardcodedAuth(testUser.ID)))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/me", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	var payload struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.ID != testUser.ID.String() {
		t.Fatalf("expected id %s, got %s", testUser.ID, payload.ID)
	}
	if payload.Username != testUser.Username {
		t.Fatalf("expected name %s, got %s", testUser.Username, payload.Username)
	}
}
