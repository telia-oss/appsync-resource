package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/telia-oss/appsync-resource"
)

func main() {
	var request resource.PutRequest

	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("missing arguments")
	}

	if err := request.Source.Validate(); err != nil {
		log.Fatalf("invalid source configuration: %s", err)
	}
	appsync, err := resource.NewAppsyncClient(&request.Source)
	if err != nil {
		log.Fatalf("failed to create appsync manager: %s", err)
	}

	response, err := resource.Put(request, appsync, os.Args[1])
	if err != nil {
		log.Fatalf("put failed: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %s", err)
	}
}
