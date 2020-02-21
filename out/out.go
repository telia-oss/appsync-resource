package out

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/telia-oss/appsync-resource"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appsync"
	yaml "gopkg.in/yaml.v2"
)

// Command will update the resource.
func Command(input InputJSON, logger *log.Logger) (outOutputJSON, error) {

	// PARSE THE JSON FILE input.json
	apiID, ok := input.Source["api_id"]
	if !ok {
		return outOutputJSON{}, errors.New("api_id not set")
	}

	accessKey, ok := input.Source["access_key_id"]
	if !ok {
		return outOutputJSON{}, errors.New("aws access_key_id not set")
	}

	secretKey, ok := input.Source["secret_access_key"]
	if !ok {
		return outOutputJSON{}, errors.New("aws secret_access_key not set")
	}

	sessionToken, ok := input.Source["session_token"]
	if !ok {
		return outOutputJSON{}, errors.New("aws session_token not set")
	}

	regionName, ok := input.Source["region_name"]
	if !ok {
		return outOutputJSON{}, errors.New("aws region_name not set")
	}

	schemaFile, _ := input.Params["schema_file"]

	resolversFile, _ := input.Params["resolvers_file"]

	var ref = input.Version.Ref
	var output outOutputJSON
	var resolverOutput []metadata
	var schemaOutput []metadata

	// AWS creds
	awsConfig := resource.NewAwsConfig(
		accessKey,
		secretKey,
		sessionToken,
		regionName,
	)

	client := resource.NewAppSyncClient(awsConfig)

	if schemaFile == "" && resolversFile == "" {
		return outOutputJSON{}, errors.New("resolversFile and schemaFile both are not set")
	}

	// Create or update schema
	if schemaFile != "" {
		schemaFilePath := fmt.Sprintf("%s/%s", os.Args[1], schemaFile)
		schema, _ := ioutil.ReadFile(schemaFilePath)

		var schemaCreateParams = &appsync.StartSchemaCreationInput{
			ApiId:      aws.String(apiID),
			Definition: schema,
		}

		var schemaStatusParams = &appsync.GetSchemaCreationStatusInput{
			ApiId: aws.String(apiID),
		}

		// Start create or update schema
		error := client.StartSchemaCreationOrUpdate(schemaCreateParams)
		if error != nil {
			logger.Fatalf("failed to create/update the schema: %s", error)
		}

		// get schema creation status
		creationStatus, creationDetails := client.GetSchemaCreationStatus(schemaStatusParams, logger)

		// OUTPUT
		schemaOutput = []metadata{
			{Name: "creationStatus", Value: creationStatus},
			{Name: "creationDetails", Value: creationDetails},
		}
	}
	// update Resolvers
	if resolversFile != "" {
		resolversFilePath := fmt.Sprintf("%s/%s", os.Args[1], resolversFile)
		resolversFile, _ := ioutil.ReadFile(resolversFilePath)
		var resolvers resource.Resolvers
		err := yaml.Unmarshal(resolversFile, &resolvers)
		if err != nil {
			panic(err)
		}

		nResolversSuccessfullyCreated, nResolversfailCreated, nResolversSuccessfullyUpdated, nResolversfailUpdate := client.CreateOrUpdateResolvers(resolvers, apiID, logger)
		// OUTPUT
		resolverOutput = []metadata{
			{Name: "number of resolvers successfully created", Value: nResolversSuccessfullyCreated},
			{Name: "number of resolver successfully updated", Value: nResolversSuccessfullyUpdated},
			{Name: "number of resolver fail to create", Value: nResolversfailCreated},
			{Name: "number of resolver fail to update", Value: nResolversfailUpdate},
		}
	}
	output = outOutputJSON{
		Version:  version{Ref: ref},
		Metadata: append(schemaOutput, resolverOutput...),
	}
	return output, nil

}
