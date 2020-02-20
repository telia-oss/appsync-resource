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

Given a schema specified by `schema_file`, to update/create AppSync  schema Or/And Given a resolvers JSON specified by `resolvers_file`, to update AppSync existing schema resolvers.

#### Parameters

* `schema_file`: *Optional.* .grapqh schema File provided by an output of a task, if you didn't specify `resolvers_file` this field is *Required.*.

* `resolvers_file`: *Optional.* .yml resolver String provided by an output of a task, if you didn't specify `schema_file` this field is *Required.*.
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
      api_id: ((api-id))
```

### Plan

``` yaml
- put: appsync-resource
  params: 
    schema_file: "path/to/schema.graphql"
    resolvers_file: "path/to/resolvers.yml"
```

### Params file example

[schema.graphql](https://github.com/telia-oss/appsync-resource/blob/master/schema.graphql)

[resolvers.yml](https://github.com/telia-oss/appsync-resource/blob/master/resolvers.yml)

