package oas

import "github.com/go-openapi/spec"

const getSpecificationOperationPath = "/openapi.yaml"

func getSpecificationOperation() spec.PathItem {
	responses := &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			StatusCodeResponses: map[int]spec.Response{
				200: {
					ResponseProps: spec.ResponseProps{
						Description: "OpenAPI specification in YAML format",
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: []string{"file"},
							},
						},
					},
				},
				500: {
					ResponseProps: spec.ResponseProps{
						Description: "Internal Server Error",
					},
				},
			},
		},
	}

	operation := &spec.Operation{
		OperationProps: spec.OperationProps{
			// https://github.com/OAI/OpenAPI-Specification/issues/110
			Produces:     []string{"application/vnd.oai.openapi;version=2.0"},
			Schemes:      nil,
			Tags:         nil,
			Summary:      "Returns this OpenAPI specification in YAML format.",
			ExternalDocs: nil,
			ID:           "X-Get-Specification",
			Responses:    responses,
		},
	}

	return spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Get: operation,
		},
	}
}
