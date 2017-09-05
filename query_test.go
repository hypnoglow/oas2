package oas2

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
)

func TestDecodeQuery(t *testing.T) {
	type (
		user struct {
			Name           string `oas:"name"`
			Sex            string `oas:"sex"`
			fieldWithNoTag string
		}

		member struct {
			Nickname    string  `oas:"nickname"`
			Age         int32   `oas:"age"`
			LovesApples bool    `oas:"loves_apples"`
			Height      float32 `oas:"height"`
		}
	)

	number := 1

	cases := []struct {
		ps            []spec.Parameter
		q             url.Values
		dst           interface{}
		expectedData  interface{}
		expectedError error
	}{
		{
			// Simple value
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name: "name",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
			},
			q: url.Values{
				"name": []string{"John"},
			},
			dst: &user{},
			expectedData: &user{
				Name: "John",
			},
		},
		{
			// query parameter that is not defined in struct
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name: "name",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
				{
					ParamProps: spec.ParamProps{
						Name: "birthdate",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
			},
			q: url.Values{
				"name":      []string{"John"},
				"birthdate": []string{"1970-01-01"},
			},
			dst: &user{},
			expectedData: &user{
				Name: "John",
			},
		},
		{
			// With default value
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name: "name",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
				{
					ParamProps: spec.ParamProps{
						Name: "sex",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type:    "string",
						Default: "Male",
					},
				},
			},
			q: url.Values{
				"name": []string{"John"},
			},
			dst: &user{},
			expectedData: &user{
				Name: "John",
				Sex:  "Male",
			},
		},
		{
			// With default value of wrong type
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name: "name",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
				{
					ParamProps: spec.ParamProps{
						Name: "sex",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type:    "string",
						Default: 123,
					},
				},
			},
			q: url.Values{
				"name": []string{"John"},
			},
			dst: &user{},
			expectedData: &user{
				Name: "John",
			},
			expectedError: fmt.Errorf("field Sex is not assignable to int"),
		},
		{
			// Different types of query parameters
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name: "nickname",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
				},
				{
					ParamProps: spec.ParamProps{
						Name: "age",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type:   "integer",
						Format: "int32",
					},
				},
				{
					ParamProps: spec.ParamProps{
						Name: "loves_apples",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type: "boolean",
					},
				},
				{
					ParamProps: spec.ParamProps{
						Name: "height",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type:   "number",
						Format: "float",
					},
				},
			},
			q: url.Values{
				"nickname":     []string{"Princess"},
				"age":          []string{"40"},
				"loves_apples": []string{"yes"},
				"height":       []string{"185.5"},
			},
			dst: &member{},
			expectedData: &member{
				Nickname:    "Princess",
				Age:         40,
				LovesApples: true,
				Height:      185.5,
			},
		},
		{
			// dst passed by value
			dst:           member{},
			expectedData:  member{},
			expectedError: fmt.Errorf("dst is not a pointer to struct (cannot modify)"),
		},
		{
			// dst is not a pointer to struct
			dst:           &number,
			expectedData:  &number,
			expectedError: fmt.Errorf("dst is not a pointer to struct (cannot modify)"),
		},
		{
			// value is not convertible
			ps: []spec.Parameter{
				{
					ParamProps: spec.ParamProps{
						Name: "age",
						In:   "query",
					},
					SimpleSchema: spec.SimpleSchema{
						Type:   "integer",
						Format: "int32",
					},
				},
			},
			q: url.Values{
				"age": []string{"Twenty Two"},
			},
			dst:          &member{},
			expectedData: &member{},
			expectedError: fmt.Errorf(
				"cannot use values %v as parameter %s with type %s and format %s",
				[]string{"Twenty Two"},
				"age",
				"integer",
				"int32",
			),
		},
	}

	for _, c := range cases {
		err := DecodeQuery(c.ps, c.q, c.dst)
		if !reflect.DeepEqual(c.expectedError, err) {
			t.Errorf("Expected error to be %v but got %v", c.expectedError, err)
		}

		if !reflect.DeepEqual(c.expectedData, c.dst) {
			t.Errorf("Expected dst to be %v but got %v", c.expectedData, c.dst)
		}
	}
}
