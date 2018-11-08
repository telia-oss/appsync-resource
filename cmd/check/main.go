package main

import (
	"encoding/json"
	"os"
)

type checkItem struct{}
type checkResponse []checkItem

func main() {
	json.NewEncoder(os.Stdout).Encode(checkResponse{})
}
