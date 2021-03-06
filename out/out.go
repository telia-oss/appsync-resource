package out

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	resource "github.com/telia-oss/appsync-resource"
)

var env = os.Getenv("ENV")

// Command will update the resource.
func Command(input InputJSON, logger *log.Logger) (outOutputJSON, error) {
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

	schemaFile, ok := input.Params["schema_file"]

	resolversFile, ok := input.Params["resolvers_file"]

	var ref = input.Version.Ref
	var output outOutputJSON
	var resolverOutput []metadata
	var schemaOutput []metadata
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
	if schemaFile != "" {
		var schemaFilePath string
		if env == "development" {
			pwd, _ := os.Getwd()
			schemaFilePath = fmt.Sprintf("%s/%s", pwd, schemaFile)
		} else {
			schemaFilePath = fmt.Sprintf("%s/%s", os.Args[1], schemaFile)
		}
		schema, serr := ioutil.ReadFile(schemaFilePath)
		if serr != nil {
			logger.Fatalf("can't read the schema file: %s", serr)
		}

		error := client.StartSchemaCreationOrUpdate(apiID, schema)
		if error != nil {
			logger.Fatalf("failed to create/update the schema: %s", error)
		}

		creationStatus, creationDetails, err := client.GetSchemaCreationStatus(apiID)
		if err != nil {
			logger.Println("Failed to get Schema Creation status, However the Schema creation might be succeeded, check the AWS console and re-tigger the build if the schema not created/updated: %s", err)
			schemaOutput = []metadata{
				{Name: "creationStatus", Value: "unknown"},
				{Name: "creationDetails", Value: "unknown"},
			}
		} else {
			schemaOutput = []metadata{
				{Name: "creationStatus", Value: creationStatus},
				{Name: "creationDetails", Value: creationDetails},
			}
		}
	}
	if resolversFile != "" {
		var resolversFilePath string
		if env == "development" {
			pwd, _ := os.Getwd()
			resolversFilePath = fmt.Sprintf("%s/%s", pwd, resolversFile)
		} else {
			resolversFilePath = fmt.Sprintf("%s/%s", os.Args[1], resolversFile)
		}
		resolversFile, rerr := ioutil.ReadFile(resolversFilePath)
		if rerr != nil {
			logger.Fatalf("can't read the resolvers file: %s", rerr)
		}

		nResolversSuccessfullyCreated, nResolversfailCreated, nResolversSuccessfullyUpdated, nResolversfailUpdate, err := client.CreateOrUpdateResolvers(apiID, resolversFile, logger)
		if err != nil {
			logger.Println("failed to create/update", err)
		}
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
