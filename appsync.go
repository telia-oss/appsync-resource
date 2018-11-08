package resource

import (
	"log"
	"time"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appsync"
)

// Appsync for testing purposes
type Appsync interface {
	GetGraphQLApi(apiName string) (*appsync.GraphqlApi, error)
	UpdateSchema(apiName string, schema []byte) error
	UpdateResolvers(apiName string, definitions []ResolverDefinition) error
}

// AppsyncClient is a wrapper around appsync functionality
type AppsyncClient struct {
	client *appsync.AppSync
}

// NewAppsyncClient Creates new AppSync client
func NewAppsyncClient(s *Source) (*AppsyncClient, error) {

	awsConfig := newAwsConfig(
		s.AccessKeyID,
		s.AccessKeySecret,
		s.SessionToken,
		s.RegionName,
	)

	session, err := session.NewSession(awsConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to create a new session: %s", err)
	}

	client := appsync.New(session, &aws.Config{Region: aws.String(s.RegionName)})
	
	return &AppsyncClient{
		client: client,
	}, nil
}

// GetGraphQLApi returns an API based on a name
func (client *AppsyncClient) GetGraphQLApi(apiName string) (*appsync.GraphqlApi, error) {
	response, err := client.client.ListGraphqlApis(&appsync.ListGraphqlApisInput{})

	if err != nil {
		return nil, fmt.Errorf("failed to list GraphQL apis: %s", err)
	}

 	for _, api := range response.GraphqlApis {
		if *api.Name == apiName {
			return api, nil
		}
	}

 	return nil, fmt.Errorf("No GraphQL api called %s found", apiName)
}

// UpdateSchema updates schema for the target GraphQL api specified by name
func (client *AppsyncClient) UpdateSchema(apiName string, schema []byte) error {
	api, err := client.GetGraphQLApi(apiName)

	if err != nil {
		return err
	}

	done := make(chan error)
	 
 	go func() {
		output, err := client.client.StartSchemaCreation(&appsync.StartSchemaCreationInput{
			ApiId:      api.ApiId,
			Definition: schema,
		})

		if err != nil {
			done <- fmt.Errorf("Failed to start creation creation: %s", err)
			return
		}
		
 		status := *output.Status
		for status == "PROCESSING" {
			time.Sleep(1000 * time.Millisecond)

 			creationStatusOutput, err := client.client.GetSchemaCreationStatus(&appsync.GetSchemaCreationStatusInput{
				ApiId: api.ApiId,
			})

			if err != nil {
				done <- fmt.Errorf("Failed to get schema creation status: %s", err)
				return
			}
	
			status = *creationStatusOutput.Status
			log.Println("Creating schema..." + status + ": " + *creationStatusOutput.Details)
		}
		
 		if status == "FAILED" {
			done <- fmt.Errorf("Failed to apply the schema, exitting")
		} else {
			done <- nil
		}
	}()

	return <-done
}

// UpdateResolvers Updates resolvers based on the definition file
func (client *AppsyncClient) UpdateResolvers(
	apiName string,
	definitions []ResolverDefinition) error {

	api, err := client.GetGraphQLApi(apiName)

	if err != nil {
		return err
	}

 	for _, resolver := range definitions {

 		output, err := client.client.GetResolver(&appsync.GetResolverInput{
			ApiId:     api.ApiId,
			FieldName: &resolver.FieldName,
			TypeName:  &resolver.TypeName,
		})

 		if output.Resolver != nil {
			_, err = client.client.DeleteResolver(&appsync.DeleteResolverInput{
				ApiId:     api.ApiId,
				FieldName: &resolver.FieldName,
				TypeName:  &resolver.TypeName,
			})
			 
			if err != nil {
				return fmt.Errorf("Failed to delete resolver %s,%s due to %s",
					resolver.TypeName,
					resolver.FieldName,
					err)
			}
		}
 		_, err = client.client.CreateResolver(&appsync.CreateResolverInput{
			ApiId:                   api.ApiId,
			DataSourceName:          &resolver.DataSource,
			FieldName:               &resolver.FieldName,
			TypeName:                &resolver.TypeName,
			RequestMappingTemplate:  &resolver.RequestMapping,
			ResponseMappingTemplate: &resolver.ResponseMapping,
		})

		if err != nil {
			return fmt.Errorf("Failed to create resolver %s,%s due to %s",
				resolver.TypeName,
				resolver.FieldName,
				err)
		}

		log.Printf("Successfully added/updated resolver: %s:%s = %s",
			resolver.TypeName,
			resolver.FieldName,
			resolver.DataSource)
	}

	return nil
}

func newAwsConfig(
	accessKey string,
	secretKey string,
	sessionToken string,
	regionName string,
) *aws.Config {
	creds := credentials.NewStaticCredentials(accessKey, secretKey, sessionToken)

	if len(regionName) == 0 {
		regionName = "eu-west-1"
	}

	awsConfig := &aws.Config{
		Region:      aws.String(regionName),
		Credentials: creds,
	}

	return awsConfig
}