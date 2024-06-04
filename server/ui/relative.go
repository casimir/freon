//go:build !embed

package ui

import (
	"embed"
)

var IsEmbedded = false

var uiFS embed.FS
