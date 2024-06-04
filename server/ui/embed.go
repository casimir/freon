//go:build embed

package ui

import (
	"embed"
)

var IsEmbedded = true

//go:embed statics/*
var uiFS embed.FS
