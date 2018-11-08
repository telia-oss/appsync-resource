[![Build Status](https://travis-ci.org/telia-oss/appsync-resource.svg?branch=master)](https://travis-ci.org/telia-oss/appsync-resource)
[![Go Report Card](https://goreportcard.com/badge/github.com/telia-oss/appsync-resource)](https://goreportcard.com/report/github.com/telia-oss/appsync-resource)
![](https://img.shields.io/maintenance/yes/2018.svg)


# AWS AppSync resource

A Concourse resource to update AppSync schema. Written in Go.

## Source Configuration

| Parameter          | Required      | Example                  | Description
| ------------------ | ------------- | -------------            | ------------------------------ |
| api_name           | Yes           | some-api-name            |                                |
| access_key_id      | Yes           | {YOUR_ACCESS_KEY_ID}     |                                |
| access_key_secret  | Yes           | {YOUR_SECRET_ACCESS_KEY} |                                |
| session_token      | Yes           | {YOUR_SESSION_TOKEN}     |                                |
| region_name        | No            | eu-west-1                | AWS region DEFAULT: eu-west-1  |

### `out`: Update or Create schema/resolvers.

Provide a schema specified by `schema` or `schemaPath`, to update AppSync existing schema or to create new AppSync schema.

Provide resolvers definition specified by `resolvers` or `resolversPath` to update or create new AppSync resolvers.

#### Parameters

* `schema`: GraphQL schema String.
* `schemaPath`: path to a GraphQL schema file provided by an output of a task for example.
* `resolvers`: resolver defintions.
* `resolversPath`: path to a yaml file provided by an output of a task for example.

## Example Configuration

### Resource type

``` yaml
resource_types:
- name: appsync-resource
    type: docker-image
    source:
      repository: teliaoss/appsync-resource
```

### Resource

``` yaml
resource:
- name: appsync-resource
    type: resource
    source:
      access_key_id: ((access-key))
      secret_access_key: ((secret-key))
      session_token: ((session-token))
      region_name: "eu-west-1"
      api_id: "znvjdp3n25epx"
```

### Plan

For updating only schema
``` yaml
- put: appsync-resource
  params: 
    schema: "schema {query:Query} type Query { getTodos: [Todo]} type Todo { id: ID! name: String description: Int priority: Int}"
```

For updating only resolvers
``` yaml
- put: appsync-resource
  params: 
    resolvers: |
      resolvers:
        dataSource: data_source_name
        typeName: SomeType
        fieldName: SomeField
        requestMapping: >
          \#set( $event = {
            "stuff": $context.arguments.input.stuff
          } )

          {
              "version" : "2017-02-28",
              "operation": "Invoke",
              "payload": $util.toJson($event)
          }
        responseMapping: $util.toJson($context.result)
```

Or both from files
``` yaml
- put: appsync-resource
  params:
      schemaPath:  output/graphql/schema.graphql
      resolversPath: output/graphql/resolvers.yml
```

With `resolvers.yml` file looking like this:

``` yaml
resolvers:
  dataSource: data_source_name
  typeName: SomeType
  fieldName: SomeField
  requestMapping: >
    \#set( $event = {
      "stuff": $context.arguments.input.stuff
    } )

    {
        "version" : "2017-02-28",
        "operation": "Invoke",
        "payload": $util.toJson($event)
    }
  responseMapping: $util.toJson($context.result)
```