package oas

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadSpec(t *testing.T) {

	t.Run("positive", func(t *testing.T) {
		fpath := "/tmp/spec.json"
		if err := ioutil.WriteFile(fpath, loadDoc().Raw(), 0755); err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		_, err := LoadSpec(fpath)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		_ = os.Remove(fpath)
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := LoadSpec("/tmp/non/existent/file.yaml")
		if err == nil {
			t.Fatal("Expected error, but got nil")
		}
	})

	t.Run("should fail on spec expansion", func(t *testing.T) {
		fpath := "/tmp/spec-that-fails-expansion.json"
		if err := ioutil.WriteFile(fpath, []byte(specThatFailsToExpand), 0755); err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		_, err := LoadSpec(fpath)
		if err == nil {
			t.Fatal("Expected error, but got nil")
		}

		_ = os.Remove(fpath)
	})

	t.Run("should fail on spec validation", func(t *testing.T) {
		fpath := "/tmp/spec-that-fails-validation.json"
		if err := ioutil.WriteFile(fpath, []byte(specThatFailsValidation), 0755); err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		_, err := LoadSpec(fpath)
		if err == nil {
			t.Fatal("Expected error, but got nil")
		}

		_ = os.Remove(fpath)
	})

}

const (
	specThatFailsToExpand = `
swagger: "2.0"
info:
  title: "Part of Swagger Petstore"
  version: "1.0.0"
host: "petstore.swagger.io"
basePath: "/v2"
tags:
- name: "pet"
paths:
  /pet/{petId}:
    get:
      tags:
      - "pet"
      summary: "Find pet by OperationID"
      description: "Returns a single pet"
      operationId: "getPetById"
      produces:
      - "application/json"
      parameters:
      - name: "petId"
        in: "path"
        description: "OperationID of pet to return"
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
        404:
          description: "Pet not found"
definitions:
  Dog:
    type: "object"
    required:
    - "id"
    - "name"
    properties:
      id:
        type: "integer"
        format: "int64"
      name:
        type: "string"
`

	specThatFailsValidation = `
swagger: "2.0"
info:
  title: "Part of Swagger Petstore"
  version: "1.0.0"
host: "petstore.swagger.io"
basePath: "/v2"
tags:
- name: "pet"
paths:
  /pet/{petId}:
    get:
      tags:
      - "pet"
      summary: "Find pet by OperationID"
      description: "Returns a single pet"
      operationId: "getPetById"
      parameters:
      - name: "petId"
        in: "path"
        description: "OperationID of pet to return"
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
        404:
          description: "Pet not found"
  /pet:
    post:
      tags:
      - "pet"
      summary: "Add a new pet to the store"
      operationId: "getPetById"
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
definitions:
  Pet:
    type: "object"
    required:
    - "id"
    - "name"
    properties:
      id:
        type: "integer"
        format: "int64"
      name:
        type: "string"
`
)
