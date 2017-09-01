package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/go-openapi/spec"

	"github.com/hypnoglow/oas2"
)

type member struct {
	Name        string `oas:"name"`
	Age         int32  `oas:"age"`
	LovesApples bool   `oas:"loves_apples"`
}

func main() {
	// In real app query will be taken from *http.Request.
	query := url.Values{"name": []string{"John"}, "age": []string{"27"}}

	// In real app parameters will be taken from spec document (yaml or json).
	paramSpec := []spec.Parameter{
		{
			ParamProps:   spec.ParamProps{Name: "name", In: "query"},
			SimpleSchema: spec.SimpleSchema{Type: "string"},
		},
		{
			ParamProps:   spec.ParamProps{Name: "age", In: "query"},
			SimpleSchema: spec.SimpleSchema{Type: "integer", Format: "int32"},
		},
		{
			ParamProps:   spec.ParamProps{Name: "loves_apples", In: "query"},
			SimpleSchema: spec.SimpleSchema{Type: "boolean", Default: true},
		},
	}

	var m member
	if err := oas2.DecodeQuery(paramSpec, query, &m); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", m)
}
