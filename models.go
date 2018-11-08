package resource

import "errors"

// Source represents the configuration for the resource.
type Source struct {
	APIName         string `json:"api_name"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SessionToken    string `json:"session_token"`
	RegionName      string `json:"region_name"`
}

// Validate validates the source object
func (s *Source) Validate() error {
	if s.APIName == "" {
		return errors.New("api_name must be set")
	}
	if s.AccessKeyID == "" {
		return errors.New("access_key_id must be set")
	}
	if s.AccessKeySecret == "" {
		return errors.New("access_key_secret must be set")
	}
	if s.SessionToken == "" {
		return errors.New("session_token must be set")
	}

	return nil
}

type ResolverDefinitionYamlFile struct {
	Resolvers []ResolverDefinition `yaml:"resolvers"`
}

type ResolverDefinition struct {
	DataSource      string `yaml:"dataSource"`
	TypeName        string `yaml:"typeName"`
	FieldName       string `yaml:"fieldName"`
	RequestMapping  string `yaml:"requestMapping"`
	ResponseMapping string `yaml:"responseMapping"`
}

type PutRequest struct {
	Source Source        `json:"source"`
	Params PutParameters `json:"params"`
}

type PutResponse struct{}

type PutParameters struct {
	Schema        string `json:"schema"`
	Resolvers     string `json:"resolvers"`
	SchemaPath    string `json:"schemaPath"`
	ResolversPath string `json:"resolversPath"`
}
