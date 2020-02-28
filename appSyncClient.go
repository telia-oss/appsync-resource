package resource

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appsync"

	yaml "gopkg.in/yaml.v2"
)

var env = os.Getenv("ENV")

type (
	Resolvers struct {
		Resolvers []Resolver `yaml:"resolvers"`
	}

	Resolver struct {
		DataSourceName          string `yaml:"dataSource"`
		FieldName               string `yaml:"fieldName"`
		RequestMappingTemplate  string `yaml:"requestMappingTemplate"`
		ResponseMappingTemplate string `yaml:"responseMappingTemplate"`
		TypeName                string `yaml:"typeName"`
	}
)

// AppSync interface
type AppSync interface {
	CreateOrUpdateResolvers(apiID string, resolversFile []byte) (string, string, string, string, error)
	StartSchemaCreationOrUpdate(apiID string, schema []byte) error
	GetSchemaCreationStatus(apiID string) (string, string, error)
}

type appSyncClient struct {
	appSyncClient *appsync.AppSync
	session       *session.Session
}

func NewAppSyncClient(
	awsConfig *aws.Config,
) AppSync {
	session := session.New(awsConfig)
	client := appsync.New(session, awsConfig)

	return &appSyncClient{
		appSyncClient: client,
		session:       session,
	}
}

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

	var awsConfig *aws.Config
	if env == "development" {
		awsConfig = &aws.Config{
			Region: aws.String(regionName),
		}
	} else {
		awsConfig = &aws.Config{
			Region:      aws.String(regionName),
			Credentials: creds,
		}
	}

	return awsConfig
}

func (client *appSyncClient) CreateOrUpdateResolvers(apiID string, resolversFile []byte) (string, string, string, string, error) {
	// number of resolvers successfully created
	var nResolversSuccessfullyCreated = 0
	// number of resolver successfully updated
	var nResolversSuccessfullyUpdated = 0
	// number of resolver fail to create
	var nResolversfailCreated = 0
	// number of resolver fail to update
	var nResolversfailUpdate = 0

	var resolvers Resolvers
	err := yaml.Unmarshal(resolversFile, &resolvers)
	if err != nil {
		return strconv.Itoa(nResolversSuccessfullyCreated), strconv.Itoa(nResolversfailCreated), strconv.Itoa(nResolversSuccessfullyUpdated), strconv.Itoa(nResolversfailUpdate), err
	}

	for _, resolver := range resolvers.Resolvers {
		resolverFieldName := resolver.FieldName
		resolverTypeName := resolver.TypeName
		resolverResp, err := client.getResolver(&appsync.GetResolverInput{
			ApiId:     aws.String(apiID),
			FieldName: aws.String(resolverFieldName),
			TypeName:  aws.String(resolverTypeName),
		})

		if err != nil {
			resolver := fmt.Sprintf("Resolver, FieldName:%s, TypeName: %s, Error: %s", resolverFieldName, resolverTypeName, err)
			fmt.Println("faild to fetch", resolver)
		}
		if resolverResp != nil {
			var params = &appsync.UpdateResolverInput{
				ApiId:                   aws.String(apiID),
				DataSourceName:          aws.String(resolver.DataSourceName),
				FieldName:               aws.String(resolver.FieldName),
				RequestMappingTemplate:  aws.String(fmt.Sprintf("%s", resolver.RequestMappingTemplate)),
				ResponseMappingTemplate: aws.String(resolver.ResponseMappingTemplate),
				TypeName:                aws.String(resolver.TypeName),
			}
			_, err := client.updateResolver(params)
			if err != nil {
				nResolversfailUpdate++
			}
			nResolversSuccessfullyUpdated++
		} else {
			var params = &appsync.CreateResolverInput{
				ApiId:                   aws.String(apiID),
				DataSourceName:          aws.String(resolver.DataSourceName),
				FieldName:               aws.String(resolver.FieldName),
				RequestMappingTemplate:  aws.String(fmt.Sprintf("%s", resolver.RequestMappingTemplate)),
				ResponseMappingTemplate: aws.String(resolver.ResponseMappingTemplate),
				TypeName:                aws.String(resolver.TypeName),
			}
			_, err := client.createResolver(params)
			if err != nil {
				nResolversfailCreated++
			}

			nResolversSuccessfullyCreated++
		}
	}
	return strconv.Itoa(nResolversSuccessfullyCreated), strconv.Itoa(nResolversfailCreated), strconv.Itoa(nResolversSuccessfullyUpdated), strconv.Itoa(nResolversfailUpdate), nil
}

func (client *appSyncClient) getResolver(params *appsync.GetResolverInput) (*appsync.Resolver, error) {
	req, resp := client.appSyncClient.GetResolverRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Resolver, nil

}

func (client *appSyncClient) updateResolver(params *appsync.UpdateResolverInput) (*appsync.Resolver, error) {
	req, resp := client.appSyncClient.UpdateResolverRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Resolver, nil
}

func (client *appSyncClient) createResolver(params *appsync.CreateResolverInput) (*appsync.Resolver, error) {
	req, resp := client.appSyncClient.CreateResolverRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Resolver, nil
}

func (client *appSyncClient) StartSchemaCreationOrUpdate(apiID string, schema []byte) error {

	schemaCreateParams := &appsync.StartSchemaCreationInput{
		ApiId:      aws.String(apiID),
		Definition: schema,
	}

	schemaStatusParams := &appsync.GetSchemaCreationStatusInput{
		ApiId: aws.String(apiID),
	}

	req, resp := client.appSyncClient.StartSchemaCreationRequest(schemaCreateParams)
	err := req.Send()
	if err != nil {
		return err

	}
	status := *resp.Status
	for status == "PROCESSING" {
		time.Sleep(3 * time.Second)

		status, _, err = client.getSchemaCreationStatus(schemaStatusParams)

		if err != nil {
			return err
		}
	}
	return nil
}

func (client *appSyncClient) getSchemaCreationStatus(schemaStatusParams *appsync.GetSchemaCreationStatusInput) (string, string, error) {
	StatusOutput, err := client.appSyncClient.GetSchemaCreationStatus(schemaStatusParams)
	if err != nil {
		return "", "", err
	}
	creationStatus := *StatusOutput.Status
	creationDetails := *StatusOutput.Details

	return creationStatus, creationDetails, nil
}

func (client *appSyncClient) GetSchemaCreationStatus(apiID string) (string, string, error) {
	schemaStatusParams := &appsync.GetSchemaCreationStatusInput{
		ApiId: aws.String(apiID),
	}

	creationStatus, creationDetails, err := client.getSchemaCreationStatus(schemaStatusParams)
	return creationStatus, creationDetails, err
}
