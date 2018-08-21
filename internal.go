package oas

import (
	"net/http"
	"strings"

	"github.com/go-openapi/spec"
)

func newMultiMiddleware(fn PathTemplateFunc, doc *Document) Middleware {
	m := make(map[methodPath]operationWithParams)
	for method, pathOps := range doc.Analyzer.Operations() {
		for path, operation := range pathOps {
			key := methodPath{
				method: method, pathTemplate: path,
			}
			value := operationWithParams{
				operation: wrapOperation(operation),
				params:    doc.Analyzer.ParametersFor(operation.ID),
			}
			m[key] = value
		}
	}

	return func(next http.Handler) http.Handler {
		return &multiMiddleware{
			next:  next,
			fn:    fn,
			doc:   doc,
			opMap: m,
		}
	}
}

// multiMiddleware is an internal middleware that adds path templates, operations
// and parameters to requests.
type multiMiddleware struct {
	next  http.Handler
	fn    PathTemplateFunc
	doc   *Document
	opMap map[methodPath]operationWithParams
}

// ServeHTTP implements http.Handler.
func (m *multiMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pt := strings.TrimPrefix(m.fn(req), m.doc.BasePath())
	if pt == "" {
		m.next.ServeHTTP(w, req)
		return
	}

	req = WithPathTemplate(req, pt)

	key := methodPath{method: strings.ToUpper(req.Method), pathTemplate: pt}
	op, ok := m.opMap[key]
	if !ok {
		m.next.ServeHTTP(w, req)
		return
	}

	req = WithOperation(req, op.operation)
	req = WithParams(req, op.params)

	m.next.ServeHTTP(w, req)
}

type methodPath struct {
	method       string
	pathTemplate string
}

type operationWithParams struct {
	operation *Operation
	params    []spec.Parameter
}
