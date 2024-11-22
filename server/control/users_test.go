package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casimir/freon/auth"
	"github.com/casimir/freon/serialize"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func setupRouter(t *testing.T, adminSession ...bool) (*gin.Engine, *auth.User) {
	var isAdmin bool
	if len(adminSession) == 1 {
		isAdmin = adminSession[0]
	}
	user := auth.MustCreateUser(t.Name()+"_session", "", isAdmin)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router.Group("", auth.HardcodedAuth(user.ID)))

	return router, user
}

func assertUserQuery(t *testing.T, router *gin.Engine, expectedID, expectedUsername string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+expectedID, nil)
	router.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var payload []serialize.Field
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	gotUser := auth.User{}
	for _, field := range payload {
		switch field.Name {
		case "ID":
			gotUser.ID = uuid.MustParse(field.Value.(string))
		case "Username":
			gotUser.Username = field.Value.(string)
		}
	}

	if expectedID != gotUser.ID.String() {
		t.Fatalf("expected id %s, got %s", expectedID, gotUser.ID)
	}
	if expectedUsername != gotUser.Username {
		t.Fatalf("expected name %s, got %s", expectedUsername, gotUser.Username)
	}
}

func TestUserRetrieve(t *testing.T) {
	otherUser := auth.MustCreateUser(t.Name(), "", false)

	router, _ := setupRouter(t)
	assertUserQuery(t, router, otherUser.ID.String(), otherUser.Username)
}

func TestUserCreate(t *testing.T) {
	username := t.Name()
	payload := []byte(fmt.Sprintf(`{"username":%q,"password":"test"}`, username))
	req, _ := http.NewRequest("POST", "/users", bytes.NewReader(payload))

	t.Run("NotAdmin", func(t *testing.T) {
		router, _ := setupRouter(t)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if http.StatusForbidden != w.Code {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("Admin", func(t *testing.T) {
		router, _ := setupRouter(t, true)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if http.StatusCreated != w.Code {
			t.Fatalf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}

		assertUserQuery(t, router, response.ID, username)
	})
}

func TestUserDelete(t *testing.T) {
	user := auth.MustCreateUser(t.Name(), "", false)

	t.Run("NotAdmin", func(t *testing.T) {
		router, _ := setupRouter(t)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/"+user.ID.String(), nil)
		router.ServeHTTP(w, req)

		if http.StatusForbidden != w.Code {
			t.Fatalf("expected status %d, got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("Admin", func(t *testing.T) {
		router, _ := setupRouter(t, true)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/"+user.ID.String(), nil)
		router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		result, err := auth.FindUserByID(user.ID.String())
		if err != nil {
			t.Fatal(err)
		}
		if result != nil {
			t.Fatalf("expected user to be deleted")
		}
	})
}

func TestUserMe(t *testing.T) {
	router, user := setupRouter(t)
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
	if user.ID.String() != payload.ID {
		t.Fatalf("expected id %s, got %s", user.ID, payload.ID)
	}
	if user.Username != payload.Username {
		t.Fatalf("expected name %s, got %s", user.Username, payload.Username)
	}
}
