package main

import (
	"errors"
	"log"
	"os"

	"github.com/casimir/freon/cmd"
	"github.com/joho/godotenv"
)

// experiments:
// - CORS
// - private token authentication
// - proxy
// - save endpoint
// - frigoligo preferences sync
// - remove all tags of an entry
// - real incremental sync (including deletions)
// - read progression
// - post-consume step
//   - tags suggestion
//   - summarization (?)
// - https://archive.org/help/wayback_api.php
// - notification (?)

func main() {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("could not load .env file: %v", err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
