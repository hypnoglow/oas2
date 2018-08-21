package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/hypnoglow/oas2"
)

func main() {
	var specPath string
	flag.StringVar(&specPath, "spec", "", "Path to spec.yaml")
	flag.Parse()

	doc, err := oas.LoadFile(specPath)
	if err != nil {
		log.Fatalln(err)
	}

	handlers := oas.OperationHandlers{
		"addPet":       http.HandlerFunc(postPet),
		"loginUser":    http.HandlerFunc(getUserLogin),
		"getInventory": http.HandlerFunc(getStoreInventory),
	}

	// We are using logrus as a debug logger for router.
	lg := logrus.New()
	lg.SetLevel(logrus.DebugLevel)

	// Prepare error handler.
	errHandler := middlewareErrorHandler(lg)

	// Create the router
	router := oas.NewRouter(
		oas.DebugLog(lg.Debugf),
		oas.RouterMiddleware(oas.QueryValidator(errHandler)),
	)

	// Setup routing by spec.
	if err = router.AddSpec(doc, handlers); err != nil {
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

func middlewareErrorHandler(log logrus.FieldLogger) oas.RequestErrorHandler {
	return func(w http.ResponseWriter, req *http.Request, err error) (resume bool) {

		switch err.(type) {
		case oas.ValidationError:
			e := err.(oas.ValidationError)
			respondClientErrors(w, e.Errors())
			return false // do not continue

		default:
			log.Error("oas middleware: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
	}
}

func respondClientErrors(w http.ResponseWriter, errs []error) {
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
