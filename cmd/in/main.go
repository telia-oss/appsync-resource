package main

import (
	"encoding/json"
	"os"
)

type getResponse struct{}

func main() {
	json.NewEncoder(os.Stdout).Encode(getResponse{})
}
