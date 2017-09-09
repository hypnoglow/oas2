package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/sirupsen/logrus"

	"github.com/hypnoglow/oas2"
)

func main() {
	var specPath string
	flag.StringVar(&specPath, "spec", "", "Path to spec.yaml")
	flag.Parse()

	doc, err := loadSpecDoc(specPath)
	if err != nil {
		log.Fatalln(err)
	}

	handlers := oas2.OperationHandlers{
		"addPet":       http.HandlerFunc(postPet),
		"loginUser":    http.HandlerFunc(getUserLogin),
		"getInventory": http.HandlerFunc(getStoreInventory),
	}

	lg := logrus.New()
	lg.SetLevel(logrus.DebugLevel)

	opts := []oas2.Option{
		oas2.LoggerOpt(lg),
		oas2.MiddlewareOpt(oas2.NewQueryValidator(doc.Spec(), errHandler).Apply),
	}

	router, err := oas2.NewRouter(doc.Spec(), handlers, opts...)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(http.ListenAndServe(":3000", router))
}

func postPet(w http.ResponseWriter, req *http.Request) {
	if _, err := io.WriteString(w, "postPet"); err != nil {
		log.Fatal(err)
	}
}

func getUserLogin(w http.ResponseWriter, req *http.Request) {
	if _, err := io.WriteString(w, "getUserLogin"); err != nil {
		log.Fatal(err)
	}
}

func getStoreInventory(w http.ResponseWriter, req *http.Request) {
	if _, err := io.WriteString(w, "getStoreInventory"); err != nil {
		log.Fatal(err)
	}
}

func errHandler(w http.ResponseWriter, errs []error) {
	// This is an example of composing an error for the response
	// from validation errors.

	type (
		errorItem struct {
			Message string      `json:"message"`
			Field   string      `json:"field"`
			Value   interface{} `json:"value"`
		}
		payload struct {
			Errors []errorItem `json:"errors"`
		}
	)

	type fielder interface {
		Field() string
	}

	type valuer interface {
		Value() interface{}
	}

	p := payload{Errors: make([]errorItem, 0)}
	for _, e := range errs {
		item := errorItem{Message: e.Error()}
		if fe, ok := e.(fielder); ok {
			item.Field = fe.Field()
		}
		if ve, ok := e.(valuer); ok {
			item.Value = ve.Value()
		}
		p.Errors = append(p.Errors, item)
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Fatalln(err)
	}
}

// loadSpecDoc loads a OpenAPI 2.0 Specification document.
func loadSpecDoc(path string) (*loads.Document, error) {
	document, err := loads.Spec(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to load spec: %s", err)
	}

	if err := spec.ExpandSpec(document.Spec(), &spec.ExpandOptions{RelativeBase: path}); err != nil {
		return nil, fmt.Errorf("Failed to expand spec: %s", err)
	}

	if err := validate.Spec(document, strfmt.Default); err != nil {
		return nil, fmt.Errorf("Spec is invalid: %s", err)
	}

	return document, nil
}
