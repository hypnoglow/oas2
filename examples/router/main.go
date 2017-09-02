package main

import (
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
		"getInventory": http.HandlerFunc(getStoreInventory),
	}

	lg := logrus.New()
	lg.SetLevel(logrus.DebugLevel)

	router, err := oas2.NewRouter(doc.Spec(), handlers, oas2.Logger(lg))
	if err != nil {
		log.Fatalln(err)
	}

	http.ListenAndServe(":3000", router)
}

func postPet(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "postPet")
}

func getStoreInventory(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "getStoreInventory")
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
