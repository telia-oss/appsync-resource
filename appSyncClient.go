package resource

import (
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
	CreateOrUpdateResolvers(resolvers Resolvers, apiID string, logger *log.Logger) (string, string, string, string)
	StartSchemaCreationOrUpdate(schemaCreateParams *appsync.StartSchemaCreationInput) error
	GetSchemaCreationStatus(schemaStatusParams *appsync.GetSchemaCreationStatusInput, logger *log.Logger) (string, string)
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

	awsConfig := &aws.Config{
		Region:      aws.String(regionName),
		Credentials: creds,
	}

	return awsConfig
}

func (client *appSyncClient) CreateOrUpdateResolvers(resolvers Resolvers, apiID string, logger *log.Logger) (string, string, string, string) {
	// number of resolvers successfully created
	var nResolversSuccessfullyCreated = 0
	// number of resolver successfully updated
	var nResolversSuccessfullyUpdated = 0
	// number of resolver fail to create
	var nResolversfailCreated = 0
	// number of resolver fail to update
	var nResolversfailUpdate = 0

	for _, resolver := range resolvers.Resolvers {
		resolverResp, err := client.getResolver(&appsync.GetResolverInput{
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
				ResponseMappingTemplate: aws.String(resolver.ResponseMappingTemplate),
				TypeName:                aws.String(resolver.TypeName),
			}
			_, err := client.updateResolver(params)
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
				ResponseMappingTemplate: aws.String(resolver.ResponseMappingTemplate),
				TypeName:                aws.String(resolver.TypeName),
			}
			_, err := client.createResolver(params)
			if err != nil {
				nResolversfailCreated++
				logger.Println("failed to fetch a resolver with error", err)
			}

			nResolversSuccessfullyCreated++
		}
	}
	return strconv.Itoa(nResolversSuccessfullyCreated), strconv.Itoa(nResolversfailCreated), strconv.Itoa(nResolversSuccessfullyUpdated), strconv.Itoa(nResolversfailUpdate)
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

func (client *appSyncClient) StartSchemaCreationOrUpdate(schemaCreateParams *appsync.StartSchemaCreationInput) error {
	req, resp := client.appSyncClient.StartSchemaCreationRequest(schemaCreateParams)
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

func (client *appSyncClient) GetSchemaCreationStatus(schemaStatusParams *appsync.GetSchemaCreationStatusInput, logger *log.Logger) (string, string) {
	StatusOutput, err := client.appSyncClient.GetSchemaCreationStatus(schemaStatusParams)
	if err != nil {
		logger.Println("Failed to get Schema Creation status, However the Schema creation might be succeeded, check the AWS console and re-tigger the build if the schema not created/updated: %s", err)
	}
	creationStatus := *StatusOutput.Status
	creationDetails := *StatusOutput.Details

	return creationStatus, creationDetails
}
