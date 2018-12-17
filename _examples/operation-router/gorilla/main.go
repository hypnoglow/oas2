package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/_examples/app"
	_ "github.com/hypnoglow/oas2/adapter/gorilla/init"
)

func main() {
	var specPath string
	flag.StringVar(&specPath, "spec", "", "Path to an OpenAPI spec file")
	flag.Parse()

	doc, err := oas.LoadFile(specPath)
	if err != nil {
		log.Fatalln(err)
	}

	err = http.ListenAndServe(":8080", api(doc))
	log.Fatal(err)
}

func api(doc *oas.Document) http.Handler {
	// Create basis that provides middlewares.
	basis := oas.NewResolvingBasis("gorilla", doc)

	srv := app.NewServer()

	// Prepare problem handler options for validation issues.
	reqProblemHandler := oas.WithProblemHandlerFunc(srv.HandleRequestProblem)
	respProblemHandler := oas.WithProblemHandlerFunc(srv.HandleResponseProblem)

	// Create the root router.
	router := mux.NewRouter()

	// Build routing for the API using operation router.
	err := basis.OperationRouter(router).
		WithOperationHandlers(map[string]http.Handler{
			"getSum":  http.HandlerFunc(srv.GetSum),
			"postSum": http.HandlerFunc(srv.PostSum),
		}).
		WithMiddleware(
			// Add content-type validators.
			basis.RequestContentTypeValidator(reqProblemHandler),
			basis.ResponseContentTypeValidator(respProblemHandler),
			// Add query & body validators.
			basis.QueryValidator(reqProblemHandler),
			basis.RequestBodyValidator(reqProblemHandler),
			basis.ResponseBodyValidator(respProblemHandler),
		).
		WithMissingOperationHandlerFunc(missingOperationHandler).
		Build()
	if err != nil {
		panic(err)
	}

	// Serve the spec itself so users can observe the API.
	router.Path("/openapi/v1").
		Methods(http.MethodGet).
		Handler(oas.NewStaticSpecHandler(doc))

	// Add healthcheck.
	router.Path("/healthz").
		Methods(http.MethodGet).
		HandlerFunc(srv.Health)

	return router
}

func missingOperationHandler(operationID string) {
	log.Printf("[WARN] missing handler for operation %s", operationID)
}
