package resource_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/telia-oss/appsync-resource"
	"github.com/telia-oss/appsync-resource/mocks"
)

func getSource() resource.Source {
	return resource.Source {
		AccessKeyID: "test",
		AccessKeySecret: "secret_test",
		APIName: "some_api",
		RegionName: "eu-west-1",
		SessionToken: "sesssion",
	}
}

func TestNoSchemaNoResolvers_NoCalls(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appsync := mock_resource.NewMockAppsync(ctrl)

	appsync.EXPECT().UpdateResolvers(gomock.Any(), gomock.Any()).Times(0)
	appsync.EXPECT().UpdateSchema(gomock.Any(), gomock.Any()).Times(0)

	_, err := resource.Put(resource.PutRequest{
		Source: getSource(),
	}, appsync, "")

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestWithSchema_UpdatesSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appsync := mock_resource.NewMockAppsync(ctrl)

	appsync.EXPECT().UpdateResolvers(gomock.Any(), gomock.Any()).Times(0)
	appsync.EXPECT().UpdateSchema("some_api", []byte("some schema")).Times(1)

	_, err := resource.Put(resource.PutRequest{
		Source: getSource(),
		Params: resource.PutParameters{
			Schema: "some schema",
		},
	}, appsync, "")

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestWithSchemaFromFile_UpdatesSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appsync := mock_resource.NewMockAppsync(ctrl)

	appsync.EXPECT().UpdateResolvers(gomock.Any(), gomock.Any()).Times(0)
	appsync.EXPECT().UpdateSchema("some_api", []byte("this is not a real schema")).Times(1)

	_, err := resource.Put(resource.PutRequest{
		Source: getSource(),
		Params: resource.PutParameters{
			SchemaPath: "mocks/data/schema.txt",
		},
	}, appsync, "")

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestWithResolversFromFile_UpdatesSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	appsync := mock_resource.NewMockAppsync(ctrl)

	appsync.EXPECT().UpdateSchema(gomock.Any(), gomock.Any()).Times(0)
	appsync.EXPECT().UpdateResolvers("some_api", gomock.Any()).Times(1)

	_, err := resource.Put(resource.PutRequest{
		Source: getSource(),
		Params: resource.PutParameters{
			ResolversPath: "mocks/data/resolvers.yml",
		},
	}, appsync, "")

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}