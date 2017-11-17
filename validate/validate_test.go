package validate

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
)

func TestQuery(t *testing.T) {
	var maxAge float64 = 18

	cases := []struct {
		ps             []spec.Parameter
		q              url.Values
		expectedErrors []error
	}{
		// not an "in: query" parameter is skipped
		{
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "name",
						In:       "path",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
			},
		},
		// error on additional parameter
		{
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "name",
						In:       "query",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
			},
			q: url.Values{"name": {"johnhoe"}, "age": {"27"}},
			expectedErrors: []error{
				ValidationErrorf("age", "27", "parameter age is unknown"),
			},
		},
		// error on parameter conversion
		{
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "age",
						In:       "query",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{
						Type:   "integer",
						Format: "int32",
					},
				},
			},
			q: url.Values{"age": {"johndoe"}},
			expectedErrors: []error{
				ValidationErrorf("age", "johndoe", "param age: cannot convert johndoe to int32"),
			},
		},
		// error on parameter validation
		{
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "age",
						In:       "query",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{
						Type:   "integer",
						Format: "int32",
					},
					CommonValidations: spec.CommonValidations{
						Minimum: &maxAge,
					},
				},
			},
			q: url.Values{"age": {"17"}},
			expectedErrors: []error{
				ValidationErrorf("age", int32(17), "age in query should be greater than or equal to 18"),
			},
		},
		// required parameter is missing
		{
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "age",
						In:       "query",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{
						Type:   "integer",
						Format: "int32",
					},
				},
			},
			q: url.Values{},
			expectedErrors: []error{
				ValidationErrorf("age", nil, "param age is required"),
			},
		},
	}

	for _, c := range cases {
		errs := Query(c.ps, c.q)
		if !reflect.DeepEqual(c.expectedErrors, errs) {
			t.Errorf("Expected errors to be\n%#v\n but got\n%#v", c.expectedErrors, errs)
		}
	}
}

func TestBody(t *testing.T) {
	cases := []struct {
		ps             []spec.Parameter
		data           interface{}
		expectedErrors []error
	}{
		// not an "in: body" parameter is skipped
		{
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "name",
						In:       "query",
						Required: true,
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
			},
		},
		// not an "in: query" parameter is skipped
		{
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name:     "user",
						In:       "body",
						Required: true,
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"object"},
								Properties: map[string]spec.Schema{
									"name": {
										SchemaProps: spec.SchemaProps{
											Type: spec.StringOrArray{"string"},
										},
									},
								},
								Required: []string{"name"},
							},
						},
					},
				},
			},
			data:           testhelperMakeUserData("John Doe"),
			expectedErrors: nil,
		},
	}

	for _, c := range cases {
		errs := Body(c.ps, c.data)
		if !reflect.DeepEqual(c.expectedErrors, errs) {
			t.Errorf("Expected errors to be %v but got %v", c.expectedErrors, errs)
		}
	}
}

func TestBySchema(t *testing.T) {
	cases := []struct {
		sch            *spec.Schema
		data           interface{}
		expectedErrors []error
	}{
		// ok
		{
			sch: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
					Properties: map[string]spec.Schema{
						"name": {
							SchemaProps: spec.SchemaProps{
								Type: spec.StringOrArray{"string"},
							},
						},
					},
					Required: []string{"name"},
				},
			},
			data: testhelperMakeUserData("John Doe"),
		},
		// string length not satisfied
		{
			sch: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: spec.StringOrArray{"object"},
					Properties: map[string]spec.Schema{
						"name": {
							SchemaProps: spec.SchemaProps{
								Type:      spec.StringOrArray{"string"},
								MinLength: int64Ptr(4),
							},
						},
					},
					Required: []string{"name"},
				},
			},
			data:           testhelperMakeUserData("Max"),
			expectedErrors: []error{ValidationErrorf("name", nil, "name in body should be at least 4 chars long")},
		},
	}

	for _, c := range cases {
		errs := BySchema(c.sch, c.data)
		if !reflect.DeepEqual(c.expectedErrors, errs) {
			t.Errorf("Expected errors to be %#v but got %#v", c.expectedErrors, errs)
		}
	}
}

func TestValidationError(t *testing.T) {
	ve := ValidationErrorf("name", nil, "name cannot be empty")

	if ve.Error() != "name cannot be empty" {
		t.Errorf("Unexpected error message")
	}

	if ve.Field() != "name" {
		t.Errorf("Unexpected error field")
	}

	if ve.Value() != nil {
		t.Errorf("Unexpected error value")
	}
}

func testhelperMakeUserData(name string) interface{} {
	var v interface{}
	js := fmt.Sprintf(`{"name": "%s"}`, name)
	if err := json.Unmarshal([]byte(js), &v); err != nil {
		panic(err)
	}
	return v
}

func int64Ptr(f int64) *int64 {
	return &f
}
