package oas2

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
)

func TestBindParams(t *testing.T) {
	type (
		user struct {
			Name string `oas:"name"`
		}

		member struct {
			Nickname    string  `oas:"nickname"`
			Age         int32   `oas:"age"`
			LovesApples bool    `oas:"loves_apples"`
			Height      float32 `oas:"height"`
		}
	)

	cases := []struct {
		ps            []spec.Parameter
		q             url.Values
		data          interface{}
		expectedData  interface{}
		expectedError error
	}{
		{
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
			data: &user{},
			expectedData: &user{
				Name: "John",
			},
		},
		{
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
			data: &member{},
			expectedData: &member{
				Nickname:    "Princess",
				Age:         40,
				LovesApples: true,
				Height:      185.5,
			},
		},
	}

	for _, c := range cases {
		err := DecodeQuery(c.ps, c.q, c.data)
		if !reflect.DeepEqual(c.expectedError, err) {
			t.Errorf("Expected error to be %v but got %v", c.expectedError, err)
		}

		if !reflect.DeepEqual(c.expectedData, c.data) {
			t.Errorf("Expected data to be %v but got %v", c.expectedData, c.data)
		}
	}
}
