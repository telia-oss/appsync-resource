package resource

import (
	"fmt"
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"
)

// Put will update the resource.
func Put(request PutRequest, manager Appsync, sourceDir string) (*PutResponse, error) {

	schema := request.Params.Schema

	if request.Params.SchemaPath != "" && len(schema) == 0 {
		data, err := ioutil.ReadFile(path.Join(sourceDir, request.Params.SchemaPath))

		if err != nil {
			return nil, err
		}

		schema = string(data)
	}

	if schema != "" {

		err := manager.UpdateSchema(
			request.Source.APIName,
			[]byte(schema),
		)

		if err != nil {
			return nil, err
		}
	}

	resolvers := request.Params.Resolvers
	if request.Params.ResolversPath != "" && len(resolvers) == 0 {
		data, err := ioutil.ReadFile(path.Join(sourceDir, request.Params.ResolversPath))

		if err != nil {
			return nil, err
		}

		resolvers = string(data)
	}

	if resolvers != "" {
		definitionFile, err := loadResolversFromYaml(resolvers)

		if err != nil {
			return nil, err
		}

		err = manager.UpdateResolvers(
			request.Source.APIName,
			definitionFile.Resolvers,
		)

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func loadResolversFromYaml(contents string) (*ResolverDefinitionYamlFile, error) {
	def := ResolverDefinitionYamlFile{}

	err := yaml.Unmarshal([]byte(contents), &def)

	if err != nil {
		return nil, fmt.Errorf("Could not parse resolver definition file (%s)", err)
	}

	return &def, nil
}
