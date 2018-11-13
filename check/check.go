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
	// InputJSON ...
	InputJSON struct {
		Params  map[string]string `json:"params"`
		Source  map[string]string `json:"source"`
		Version version           `json:"version"`
	}
	metadata struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	checkOutputJSON []version
	inOutputJSON    struct {
		Version  version    `json:"version"`
		Metadata []metadata `json:"metadata"`
	}
	outOutputJSON inOutputJSON
)

// createOutput ...
func createOutput(output interface{}, encoder *json.Encoder, logger *log.Logger) error {
	b, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	logger.Print(string(b))

	// encode output to stdout
	return encoder.Encode(output)

}

// Check ...
func Check(input InputJSON, logger *log.Logger) (outOutputJSON, error) {
	var ref = input.Version.Ref

	// OUTPUT
	output := outOutputJSON{
		Version: version{Ref: ref},
		Metadata: []metadata{
			{Name: "", Value: ""},
			{Name: "", Value: ""},
		},
	}

	return output, nil

}

func main() {

	var (
		input   InputJSON
		decoder = json.NewDecoder(os.Stdin)
		encoder = json.NewEncoder(os.Stdout)
		logger  = log.New(os.Stderr, "resource:", log.Lshortfile)
	)

	if err := decoder.Decode(&input); err != nil {
		logger.Fatalf("Failed to decode to stdin: %s", err)
	}

	output, err := Check(input, logger)
	if err != nil {
		logger.Fatalf("Input missing a field: %s", err)
	}

	if err := createOutput(output, encoder, logger); err != nil {
		logger.Fatalf("Failed to encode to stdout: %s", err)
	}

}
