package main

import (
	"encoding/json"
	"log"
	"os"
)

type (
	version struct {
		Ref string `json:"ref"`
	}
	InputJSON struct {
		Params  map[string]string `json:"params"`
		Source  map[string]string `json:"source"`
		Version version           `json:"version"`
	}
)

func main() {

	var (
		input   InputJSON
		decoder = json.NewDecoder(os.Stdin)
		logger  = log.New(os.Stderr, "resource:", log.Lshortfile)
	)

	if err := decoder.Decode(&input); err != nil {
		logger.Fatalf("Failed to decode to stdin: %s", err)
	}

}
