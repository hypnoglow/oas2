package oas2

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

func loadDoc() *loads.Document {
	yml, err := swag.BytesToYAMLDoc(sp)
	if err != nil {
		log.Fatalln(err)
	}
	jsn, err := swag.YAMLToJSON(yml)
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := loads.Analyzed(jsn, "2.0")
	if err != nil {
		log.Fatalln(err)
	}

	if err := spec.ExpandSpec(doc.Spec(), &spec.ExpandOptions{}); err != nil {
		log.Fatalln(err)
	}

	if err := validate.Spec(doc, strfmt.Default); err != nil {
		log.Fatalln(err)
	}

	return doc
}

var sp = []byte(`
swagger: "2.0"
info:
  description: "This is a sample server Petstore server."
  version: "1.0.0"
  title: "Swagger Petstore"
  termsOfService: "http://swagger.io/terms/"
  contact:
    email: "apiteam@swagger.io"
  license:
    name: "Apache 2.0"
    url: "http://www.apache.org/licenses/LICENSE-2.0.html"
host: "petstore.swagger.io"
basePath: "/v2"
tags:
- name: "pet"
  description: "Everything about your Pets"
  externalDocs:
    description: "Find out more"
    url: "http://swagger.io"
schemes:
- "http"
paths:
  /pet:
    post:
      tags:
      - "pet"
      summary: "Add a new pet to the store"
      operationId: "addPet"
      consumes:
      - "application/json"
      produces:
      - "application/json"
      parameters:
      - in: "body"
        name: "body"
        description: "Pet object that needs to be added to the store"
        required: true
        schema:
          $ref: "#/definitions/Pet"
      - in: query
        name: debug
        type: boolean
      responses:
        405:
          description: "Invalid input"
      security:
      - petstore_auth:
        - "write:pets"
        - "read:pets"
  /pet/{petId}:
    get:
      tags:
      - "pet"
      summary: "Find pet by ID"
      description: "Returns a single pet"
      operationId: "getPetById"
      produces:
      - "application/json"
      parameters:
      - name: "petId"
        in: "path"
        description: "ID of pet to return"
        required: true
        type: "integer"
        format: "int64"
      - in: query
        name: debug
        type: boolean
      responses:
        200:
          description: "successful operation"
          schema:
            $ref: "#/definitions/Pet"
        400:
          description: "Invalid ID supplied"
        404:
          description: "Pet not found"
      security:
      - api_key: []
  /user/login:
    get:
      tags:
      - "user"
      summary: "Logs user into the system"
      description: ""
      operationId: "loginUser"
      produces:
      - "application/json"
      parameters:
      - name: "username"
        in: "query"
        description: "The user name for login"
        required: true
        type: "string"
      - name: "password"
        in: "query"
        description: "The password for login in clear text"
        required: true
        type: "string"
      responses:
        200:
          description: "successful operation"
          schema:
            type: "string"
          headers:
            X-Rate-Limit:
              type: "integer"
              format: "int32"
              description: "calls per hour allowed by the user"
            X-Expires-After:
              type: "string"
              format: "date-time"
              description: "date in UTC when token expires"
        400:
          description: "Invalid username/password supplied"
securityDefinitions:
  petstore_auth:
    type: "oauth2"
    authorizationUrl: "http://petstore.swagger.io/oauth/dialog"
    flow: "implicit"
    scopes:
      write:pets: "modify pets in your account"
      read:pets: "read your pets"
  api_key:
    type: "apiKey"
    name: "api_key"
    in: "header"
definitions:
  Pet:
    type: "object"
    required:
    - "name"
    - "age"
    properties:
      id:
        type: "integer"
        format: "int64"
      name:
        type: "string"
        example: "doggie"
      age:
        type: "integer"
        format: "int32"
        example: 7
      tags:
        type: "array"
        items:
          type: "array"
      status:
        type: "string"
        description: "pet status in the store"
        enum:
        - "available"
        - "pending"
        - "sold"
  ApiResponse:
    type: "object"
    properties:
      code:
        type: "integer"
        format: "int32"
      type:
        type: "string"
      message:
        type: "string"
externalDocs:
  description: "Find out more about Swagger"
  url: "http://swagger.io"
`)
