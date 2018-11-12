package out

import (
	"log"
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
