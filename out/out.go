package out

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appsync"
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

	RequestMappingTemplate struct {
		Version   string
		Operation string
		Payload   string
	}
	Resolver struct {
		DataSourceName         string
		FieldName              string
		RequestMappingTemplate RequestMappingTemplate
		ResponseMapping        string
		TypeName               string
	}
)

func NewAwsConfig(
	accessKey string,
	secretKey string,
	sessionToken string,
	regionName string,
) *aws.Config {
	var creds *credentials.Credentials

	if accessKey == "" && secretKey == "" {
		creds = credentials.AnonymousCredentials
	} else {
		creds = credentials.NewStaticCredentials(accessKey, secretKey, sessionToken)
	}

	if len(regionName) == 0 {
		regionName = "eu-west-1"
	}

	awsConfig := &aws.Config{
		Region:      aws.String(regionName),
		Credentials: creds,
	}

	return awsConfig
}
func createOrUpdateResolvers(appsyncClient *appsync.AppSync, resolvers []Resolver, apiID string, logger *log.Logger) (string, string, string, string) {
	// number of resolvers successfully created
	var nResolversSuccessfullyCreated = 0
	// number of resolver successfully updated
	var nResolversSuccessfullyUpdated = 0
	// number of resolver fail to create
	var nResolversfailCreated = 0
	// number of resolver fail to update
	var nResolversfailUpdate = 0

	for _, resolver := range resolvers {
		resolverResp, err := getResolver(appsyncClient, &appsync.GetResolverInput{
			ApiId:     aws.String(apiID),
			FieldName: aws.String(resolver.FieldName),
			TypeName:  aws.String(resolver.TypeName),
		})

		if err != nil {
			logger.Println("failed to fetch a resolver with error", err)
		}
		if resolverResp != nil {
			var params = &appsync.UpdateResolverInput{
				ApiId:                   aws.String(apiID),
				DataSourceName:          aws.String(resolver.DataSourceName),
				FieldName:               aws.String(resolver.FieldName),
				RequestMappingTemplate:  aws.String(fmt.Sprintf("%s", resolver.RequestMappingTemplate)),
				ResponseMappingTemplate: aws.String(resolver.ResponseMapping),
				TypeName:                aws.String(resolver.TypeName),
			}
			_, err := updateResolver(appsyncClient, params)
			if err != nil {
				nResolversfailUpdate++
				logger.Println("failed to fetch a resolver with error", err)
			}
			nResolversSuccessfullyUpdated++
		} else {
			var params = &appsync.CreateResolverInput{
				ApiId:                   aws.String(apiID),
				DataSourceName:          aws.String(resolver.DataSourceName),
				FieldName:               aws.String(resolver.FieldName),
				RequestMappingTemplate:  aws.String(fmt.Sprintf("%s", resolver.RequestMappingTemplate)),
				ResponseMappingTemplate: aws.String(resolver.ResponseMapping),
				TypeName:                aws.String(resolver.TypeName),
			}
			_, err := createResolver(appsyncClient, params)
			if err != nil {
				nResolversfailCreated++
				logger.Println("failed to fetch a resolver with error", err)
			}

			nResolversSuccessfullyCreated++
		}
	}
	return strconv.Itoa(nResolversSuccessfullyCreated), strconv.Itoa(nResolversfailCreated), strconv.Itoa(nResolversSuccessfullyUpdated), strconv.Itoa(nResolversfailUpdate)
}
func getResolver(appsyncClient *appsync.AppSync, params *appsync.GetResolverInput) (*appsync.Resolver, error) {
	req, resp := appsyncClient.GetResolverRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Resolver, nil

}

func updateResolver(appsyncClient *appsync.AppSync, params *appsync.UpdateResolverInput) (*appsync.Resolver, error) {
	req, resp := appsyncClient.UpdateResolverRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Resolver, nil
}

func createResolver(appsyncClient *appsync.AppSync, params *appsync.CreateResolverInput) (*appsync.Resolver, error) {
	req, resp := appsyncClient.CreateResolverRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Resolver, nil
}

func startSchemaCreationOrUpdate(appsyncClient *appsync.AppSync, schemaCreateParams *appsync.StartSchemaCreationInput) error {
	req, resp := appsyncClient.StartSchemaCreationRequest(schemaCreateParams)
	err := req.Send()
	if err != nil {
		return err

	}
	status := *resp.Status
	if status == "PROCESSING" {
		time.Sleep(time.Second * 3)
	}
	return nil
}

func getSchemaCreationStatus(appsyncClient *appsync.AppSync, schemaStatusParams *appsync.GetSchemaCreationStatusInput, logger *log.Logger) (string, string) {
	StatusOutput, err := appsyncClient.GetSchemaCreationStatus(schemaStatusParams)
	if err != nil {
		logger.Println("Failed to get Schema Creation status, However the Schema creation might be succeeded, check the AWS console and re-tigger the build if the schema not created/updated: %s", err)
	}
	creationStatus := *StatusOutput.Status
	creationDetails := *StatusOutput.Details

	return creationStatus, creationDetails
}

// Out will update the resource.
func Out(input InputJSON, logger *log.Logger) (outOutputJSON, error) {

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

	schemaContent, ok := input.Params["schemaContent"]
	if !ok {
		return outOutputJSON{}, errors.New("schemaContent not set")
	}

	resolversContent, ok := input.Params["resolversContent"]
	if !ok {
		return outOutputJSON{}, errors.New("resolversContent not set")
	}

	var ref = input.Version.Ref
	var output outOutputJSON

	// AWS creds
	awsConfig := NewAwsConfig(
		accessKey,
		secretKey,
		sessionToken,
		regionName,
	)
	//update schema
	session, err := session.NewSession(awsConfig)
	if err != nil {
		logger.Fatalf("failed to create a new session: %s", err)
	}

	// Create a AppSync client from just a session.
	appsyncClient := appsync.New(session)

	// Create or update schema
	// TODO refactor code
	if schemaContent != "" {
		schema := []byte(schemaContent)
		var schemaCreateParams = &appsync.StartSchemaCreationInput{
			ApiId:      aws.String(apiID),
			Definition: schema,
		}

		var schemaStatusParams = &appsync.GetSchemaCreationStatusInput{
			ApiId: aws.String(apiID),
		}

		// Start create or update schema
		error := startSchemaCreationOrUpdate(appsyncClient, schemaCreateParams)
		if error != nil {
			logger.Fatalf("failed to create/update the schema: %s", error)
		}

		// get schema creation status
		creationStatus, creationDetails := getSchemaCreationStatus(appsyncClient, schemaStatusParams, logger)

		// OUTPUT
		output = outOutputJSON{
			Version: version{Ref: ref},
			Metadata: []metadata{
				{Name: "creationStatus", Value: creationStatus},
				{Name: "creationDetails", Value: creationDetails},
			},
		}
	}
	// update Resolvers
	if resolversContent != "" {

		resolverJSONTpl := fmt.Sprintf("`%s`", resolversContent)
		var resolvers []Resolver
		val := []byte(resolverJSONTpl)
		s, _ := strconv.Unquote(string(val))
		json.Unmarshal([]byte(s), &resolvers)

		nResolversSuccessfullyCreated, nResolversfailCreated, nResolversSuccessfullyUpdated, nResolversfailUpdate := createOrUpdateResolvers(appsyncClient, resolvers, apiID, logger)
		// OUTPUT
		output = outOutputJSON{
			Version: version{Ref: ref},
			Metadata: []metadata{
				{Name: "number of resolvers successfully created", Value: nResolversSuccessfullyCreated},
				{Name: "number of resolver successfully updated", Value: nResolversSuccessfullyUpdated},
				{Name: "number of resolver fail to create", Value: nResolversfailCreated},
				{Name: "number of resolver fail to update", Value: nResolversfailUpdate},
			},
		}
	}
	return output, nil

}
