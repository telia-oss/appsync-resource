[![Build Status](https://travis-ci.org/telia-oss/appsync-resource.svg?branch=master)](https://travis-ci.org/telia-oss/appsync-resource)
[![Go Report Card](https://goreportcard.com/badge/github.com/telia-oss/appsync-resource)](https://goreportcard.com/report/github.com/telia-oss/appsync-resource)
![](https://img.shields.io/maintenance/yes/2018.svg)


# AWS AppSync resource

A Concourse resource to update AppSync schema. Written in Go.

## Source Configuration

| Parameter          | Required      | Example                  | Description
| ------------------ | ------------- | -------------            | ------------------------------ |
| api_id             | Yes           | znvjdp3n25epx            |                                |
| access_key_id      | Yes           | {YOUR_ACCESS_KEY_ID}     |                                |
| secret_access_key  | Yes           | {YOUR_SECRET_ACCESS_KEY} |                                |
| session_token      | Yes           | {YOUR_SESSION_TOKEN}     |                                |
| region_name        | No            | eu-west-1                | AWS region DEFAULT: eu-west-1  |

### `out`: Update or Create schema.

Given a schema specified by `schemaContent`, to update/create AppSync  schema Or/And Given a resolvers JSON specified by `resolversContent`, to update AppSync existing schema resolvers.

#### Parameters

* `schemaContent`: *Optional.* .grapqh schema String provided by an output of a task, if you didn't specify `resolversContent` this field is *Required.*.

* `resolversContent`: *Optional.* .json resolver String provided by an output of a task, if you didn't specify `schemaContent` this field is *Required.*.
.

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

``` yaml
- put: appsync-resource
  params: 
    schemaContent: "schema {query:Query} type Query { getTodos: [Todo]} type Todo { id: ID! name: String description: Int priority: Int}"
    resolversContent: "{\"dataSourceName\": \"testd\", \"fieldName\": \"getTodos\", \"requestMappingTemplate\": {\"version\": \"2017-02-28\", \"operation\": \"Invoke\", \"payload\": \"$util.toJson($context.args)\"}, \"responseMapping\": \"$util.toJson($context.result)\", \"typeName\": \"Query\"}"
```


