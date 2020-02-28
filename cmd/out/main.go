package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/telia-oss/appsync-resource/out"
)

func createOutput(output interface{}, encoder *json.Encoder, logger *log.Logger) error {
	_, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	return encoder.Encode(output)

}

func main() {

	var (
		input   out.InputJSON
		decoder = json.NewDecoder(os.Stdin)
		encoder = json.NewEncoder(os.Stdout)
		logger  = log.New(os.Stderr, "resource:", log.Lshortfile)
	)

	if err := decoder.Decode(&input); err != nil {
		logger.Fatalf("Failed to decode to stdin: %s", err)
	}

	output, err := out.Command(input, logger)
	if err != nil {
		logger.Fatalf("Input missing a field: %s", err)
	}

	if err := createOutput(output, encoder, logger); err != nil {
		logger.Fatalf("Failed to encode to stdout: %s", err)
	}

}
