package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/_examples/app"
	"github.com/hypnoglow/oas2/adapter/gorilla"
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
	basis := oas.NewResolvingBasis(doc, oas_gorilla.NewResolver(doc))

	srv := app.NewServer()

	// Prepare problem handler options for validation issues.
	reqProblemHandler := oas.WithProblemHandlerFunc(srv.HandleRequestProblem)
	respProblemHandler := oas.WithProblemHandlerFunc(srv.HandleResponseProblem)

	// Create the root router.
	router := mux.NewRouter()

	// Build routing for the API using subrouter.
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(
		// First of all, use operation context middleware so other oas middlewares
		// can function properly.
		mux.MiddlewareFunc(basis.OperationContext()),
		// Add content-type validators.
		mux.MiddlewareFunc(basis.RequestContentTypeValidator(reqProblemHandler)),
		mux.MiddlewareFunc(basis.ResponseContentTypeValidator(respProblemHandler)),
		// Add query & body validators.
		mux.MiddlewareFunc(basis.QueryValidator(reqProblemHandler)),
		mux.MiddlewareFunc(basis.RequestBodyValidator(reqProblemHandler)),
		mux.MiddlewareFunc(basis.ResponseBodyValidator(respProblemHandler)),
	)
	// Handle routes.
	apiRouter.Path("/sum").Methods(http.MethodGet).HandlerFunc(srv.GetSum)
	apiRouter.Path("/sum").Methods(http.MethodPost).HandlerFunc(srv.PostSum)

	// Serve the spec itself so users can observe the API.
	router.Path("/openapi/v1").Methods(http.MethodGet).Handler(oas.NewStaticSpecHandler(doc))

	// Add healthcheck.
	router.Path("/healthz").Methods(http.MethodGet).HandlerFunc(srv.Health)

	return router
}
