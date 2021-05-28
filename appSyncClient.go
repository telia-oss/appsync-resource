package resource

import (
	"fmt"
	"log"
	"os"
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
		Functions []Function `yaml:"functions"`
	}

	Resolver struct {
		DataSourceName          string   `yaml:"dataSource"`
		FieldName               string   `yaml:"fieldName"`
		RequestMappingTemplate  string   `yaml:"requestMappingTemplate"`
		ResponseMappingTemplate string   `yaml:"responseMappingTemplate"`
		Functions               []string `yaml:"functions"`
		TypeName                string   `yaml:"typeName"`
	}

	Function struct {
		DataSourceName          string `yaml:"dataSource"`
		Name                    string `yaml:"name"`
		RequestMappingTemplate  string `yaml:"requestMappingTemplate"`
		ResponseMappingTemplate string `yaml:"responseMappingTemplate"`
	}

	Statistics struct {
		Created        int
		Updated        int
		FailedToCreate int
		FailedToUpdate int
	}
)

// AppSync interface
type AppSync interface {
	CreateOrUpdateResolvers(apiID string, resolversFile []byte, logger *log.Logger) (Statistics, Statistics, error)
	StartSchemaCreationOrUpdate(apiID string, schema []byte) error
	GetSchemaCreationStatus(apiID string) (string, string, error)
}

type appSyncClient struct {
	appSyncClient *appsync.AppSync
	session       *session.Session
}

func NewAppSyncClient(
	awsConfig *aws.Config,
) (AppSync, error) {
	session, err := session.NewSession(awsConfig)

	if err != nil {
		return nil, err
	}

	client := appsync.New(session, awsConfig)

	return &appSyncClient{
		appSyncClient: client,
		session:       session,
	}, nil
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

func (client *appSyncClient) CreateOrUpdateResolvers(apiID string, resolversFile []byte, logger *log.Logger) (Statistics, Statistics, error) {
	resolverStatistics := Statistics{}
	functionStatistics := Statistics{}

	var resolvers Resolvers
	err := yaml.Unmarshal(resolversFile, &resolvers)
	if err != nil {
		return resolverStatistics, functionStatistics, err
	}

	functions, err := client.getFunctions(apiID)

	if err != nil {
		return resolverStatistics, functionStatistics, err
	}

	for _, function := range resolvers.Functions {
		functionName := function.Name

		existingFunction := getFunctionByName(functionName, functions)

		if existingFunction != nil {
			_, err := client.updateFunction(&appsync.UpdateFunctionInput{
				ApiId:                   aws.String(apiID),
				DataSourceName:          aws.String(function.DataSourceName),
				RequestMappingTemplate:  aws.String(function.RequestMappingTemplate),
				ResponseMappingTemplate: aws.String(function.ResponseMappingTemplate),
				FunctionId:              existingFunction.FunctionId,
				Description:             aws.String(function.Name),
				Name:                    aws.String(function.Name),
				FunctionVersion:         aws.String("2018-05-29"),
			})
			if err != nil {
				logger.Println(fmt.Sprintf("Function %s failed to update: %s", function.Name, err))
				functionStatistics.FailedToUpdate++
			} else {
				functionStatistics.Updated++
			}
		} else {
			function, err := client.createFunction(&appsync.CreateFunctionInput{
				ApiId:                   aws.String(apiID),
				DataSourceName:          aws.String(function.DataSourceName),
				RequestMappingTemplate:  aws.String(function.RequestMappingTemplate),
				ResponseMappingTemplate: aws.String(function.ResponseMappingTemplate),
				Description:             aws.String(function.Name),
				Name:                    aws.String(function.Name),
				FunctionVersion:         aws.String("2018-05-29"),
			})

			if err != nil {
				logger.Println(fmt.Sprintf("Function %s failed to create: %s", *function.Name, err))
				functionStatistics.FailedToCreate++
			} else {
				functions = append(functions, function)
				functionStatistics.Created++
			}
		}
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
			logger.Println("faild to fetch", resolver)
		}

		dataSourceName := aws.String(resolver.DataSourceName)
		var pipelineConfig *appsync.PipelineConfig
		resolverKind := appsync.ResolverKindUnit

		shouldContinue := true
		if len(resolver.Functions) > 0 {
			resolverKind = appsync.ResolverKindPipeline
			pipelineConfig = &appsync.PipelineConfig{}

			for i := range resolver.Functions {
				existingFunction := getFunctionByName(resolver.Functions[i], functions)

				if existingFunction == nil {
					logger.Printf("Failed to find function: %s, I will not continue with updating resolver type: %s, field: %s\n",
						resolver.Functions[i], resolver.TypeName, resolver.FieldName)
					shouldContinue = false
					break
				}

				pipelineConfig.Functions = append(pipelineConfig.Functions, existingFunction.FunctionId)
			}

			dataSourceName = nil
		}

		if !shouldContinue {
			continue
		}

		if resolverResp != nil {
			var params = &appsync.UpdateResolverInput{
				ApiId:                   aws.String(apiID),
				DataSourceName:          dataSourceName,
				FieldName:               aws.String(resolver.FieldName),
				RequestMappingTemplate:  aws.String(resolver.RequestMappingTemplate),
				ResponseMappingTemplate: aws.String(resolver.ResponseMappingTemplate),
				TypeName:                aws.String(resolver.TypeName),
				Kind:                    &resolverKind,
				PipelineConfig:          pipelineConfig,
			}
			_, err := client.updateResolver(params)
			if err != nil {
				logger.Println(fmt.Sprintf("Resolver on type %s and field %s failed to update: %s", resolver.TypeName, resolver.FieldName, err))
				resolverStatistics.FailedToUpdate++
			} else {
				resolverStatistics.Updated++
			}
		} else {
			var params = &appsync.CreateResolverInput{
				ApiId:                   aws.String(apiID),
				DataSourceName:          dataSourceName,
				FieldName:               aws.String(resolver.FieldName),
				RequestMappingTemplate:  aws.String(resolver.RequestMappingTemplate),
				ResponseMappingTemplate: aws.String(resolver.ResponseMappingTemplate),
				TypeName:                aws.String(resolver.TypeName),
				Kind:                    &resolverKind,
				PipelineConfig:          pipelineConfig,
			}
			_, err := client.createResolver(params)
			if err != nil {
				logger.Println(fmt.Sprintf("Resolver on type %s and field %s failed to create: %s", resolver.TypeName, resolver.FieldName, err))
				resolverStatistics.FailedToCreate++
			} else {
				resolverStatistics.Created++
			}
		}
	}

	return resolverStatistics, functionStatistics, nil
}

func getFunctionByName(name string, functions []*appsync.FunctionConfiguration) *appsync.FunctionConfiguration {
	for _, function := range functions {
		if *function.Name == name {
			return function
		}
	}

	return nil
}

func (client *appSyncClient) getResolver(params *appsync.GetResolverInput) (*appsync.Resolver, error) {
	req, resp := client.appSyncClient.GetResolverRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Resolver, nil

}

func (client *appSyncClient) getFunctions(apiID string) ([]*appsync.FunctionConfiguration, error) {
	var out []*appsync.FunctionConfiguration

	params := appsync.ListFunctionsInput{
		ApiId:      aws.String(apiID),
		MaxResults: aws.Int64(25),
		NextToken:  nil,
	}

	for {
		req, resp := client.appSyncClient.ListFunctionsRequest(&params)

		err := req.Send()
		if err != nil {
			return nil, err
		}

		out = append(out, resp.Functions...)

		if resp.NextToken == nil {
			break
		}

		params.NextToken = resp.NextToken
	}

	return out, nil

}

func (client *appSyncClient) updateResolver(params *appsync.UpdateResolverInput) (*appsync.Resolver, error) {
	req, resp := client.appSyncClient.UpdateResolverRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Resolver, nil
}

func (client *appSyncClient) updateFunction(params *appsync.UpdateFunctionInput) (*appsync.FunctionConfiguration, error) {
	req, resp := client.appSyncClient.UpdateFunctionRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.FunctionConfiguration, nil
}

func (client *appSyncClient) createResolver(params *appsync.CreateResolverInput) (*appsync.Resolver, error) {
	req, resp := client.appSyncClient.CreateResolverRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.Resolver, nil
}

func (client *appSyncClient) createFunction(params *appsync.CreateFunctionInput) (*appsync.FunctionConfiguration, error) {
	req, resp := client.appSyncClient.CreateFunctionRequest(params)

	err := req.Send()
	if err != nil {
		return nil, err
	}

	return resp.FunctionConfiguration, nil
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
