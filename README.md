[![Build Status](https://travis-ci.org/telia-oss/appsync-resource.svg?branch=master)](https://travis-ci.org/telia-oss/appsync-resource)
[![Go Report Card](https://goreportcard.com/badge/github.com/telia-oss/appsync-resource)](https://goreportcard.com/report/github.com/telia-oss/appsync-resource)

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

Given a schema specified by `schemaContent`, to update AppSync existing schema to create AppSync new schema.

#### Parameters

* `schemaContent`: *Required.* .grapqh schema String provided by an output of a task.


## Example Configuration

### Resource type

``` yaml
resource_types:
- name: appsync-resource
    type: docker-image
    source:
      repository: mhd999/resource
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
```
