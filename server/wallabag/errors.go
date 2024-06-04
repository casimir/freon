package wallabag

import (
	"fmt"
	"io"
	"net/http"
)

type WallabagApiError struct {
	response *http.Response
}

func (e WallabagApiError) Error() string {
	body, _ := io.ReadAll(e.response.Body)
	return fmt.Sprintf("wallabag API error: %s: %s", e.response.Status, string(body))
}

type WallabagNotAuthenticatedError struct{}

func (e WallabagNotAuthenticatedError) Error() string {
	return "not authenticated (a new authentication with password is needed)"
}

type InvalidOptionError struct {
	Field string
}

func (e *InvalidOptionError) Error() string {
	return "invalid option: " + e.Field
}
