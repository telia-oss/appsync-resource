package cmd

import (
	"log"
)

// In ...
func In(input InputJSON, logger *log.Logger) (outOutputJSON, error) {
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
