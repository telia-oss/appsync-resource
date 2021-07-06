package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/telia-oss/appsync-resource/out"
)

func createOutput(output interface{}, encoder *json.Encoder, logger *log.Logger) error {
	return encoder.Encode(output)
}

func main() {

	var (
		input   out.InputJSON
		decoder = json.NewDecoder(os.Stdin)
		encoder = json.NewEncoder(os.Stdout)
		logger  = log.New(os.Stderr, "resource:", log.Lshortfile)
	)

	whichCi := os.Getenv("INPUT_CI")

	if whichCi == "github" {
		sessionToken := os.Getenv("INPUT_SESSION_TOKEN")
		secretAccessKey := os.Getenv("INPUT_SECRET_ACCESS_KEY")
		accessKeyId := os.Getenv("INPUT_ACCESS_KEY_ID")
		regionName := os.Getenv("INPUT_REGION_NAME")
		apiID := os.Getenv("INPUT_API_ID")
		schemaFile := os.Getenv("INPUT_SCHEMA_FILE")
		resolversFile := os.Getenv("INPUT_RESOLVERS_FILE")

		input.Source = make(map[string]string)
		input.Params = make(map[string]string)
		input.Source["api_id"] = apiID
		input.Source["access_key_id"] = accessKeyId
		input.Source["secret_access_key"] = secretAccessKey
		input.Source["session_token"] = sessionToken
		input.Source["region_name"] = regionName
		input.Params["schema_file"] = schemaFile
		input.Params["resolvers_file"] = resolversFile
	} else if err := decoder.Decode(&input); err != nil {
		logger.Fatalf("Failed to decode to stdin: %s", err)
	}

	output, err := out.Command(input, logger)
	if err != nil {
		logger.Fatalf("Error execute out command: %s", err)
	}

	if err := createOutput(output, encoder, logger); err != nil {
		logger.Fatalf("Failed to encode to stdout: %s", err)
	}

}
