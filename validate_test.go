package oas2

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
)

func TestValidateQuery(t *testing.T) {
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
	}

	for _, c := range cases {
		errs := ValidateQuery(c.ps, c.q)
		if !reflect.DeepEqual(c.expectedErrors, errs) {
			t.Errorf("Expected errors to be\n%#v\n but got\n%#v", c.expectedErrors, errs)
		}
	}
}
