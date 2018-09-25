// Package app provides common application logic for all examples.
package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/hypnoglow/oas2"
)

func NewServer() *Server {
	return &Server{}
}

type Server struct {
	sumAccumulator int64
	sumCount       int64
	mx             sync.Mutex
}

func (srv *Server) PostSum(w http.ResponseWriter, req *http.Request) {
	var params postSumRequest
	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srv.mx.Lock()

	srv.sumAccumulator += params.Number
	srv.sumCount += 1

	resp := postSumResponse{Sum: srv.sumAccumulator}

	srv.mx.Unlock()

	srv.respondJSON(w, resp)
}

type postSumRequest struct {
	Number int64 `json:"number"`
}

type postSumResponse struct {
	Sum int64 `json:"sum"`
}

func (srv *Server) GetSum(w http.ResponseWriter, req *http.Request) {
	var params getSumRequest
	if err := oas.DecodeQuery(req, &params); err != nil {
		log.Printf("[WARN] %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	srv.mx.Lock()

	resp := getSumResponse{Sum: srv.sumAccumulator}
	if params.Count {
		cnt := srv.sumCount
		resp.Count = &cnt
	}

	srv.mx.Unlock()

	srv.respondJSON(w, resp)
}

type getSumRequest struct {
	Count bool `oas:"count"`
}

type getSumResponse struct {
	Sum   int64  `json:"sum"`
	Count *int64 `json:"count,omitempty"`
}

func (srv *Server) respondJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("[ERROR] %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (srv *Server) HandleRequestProblem(p oas.Problem) {
	resp := errorResponse{}

	switch te := p.Cause().(type) {
	case oas.MultiError:
		for _, e := range te.Errors() {
			resp.Errors = append(resp.Errors,
				fmt.Sprintf("%s: %v", te.Message(), e),
			)
		}
	default:
		resp.Errors = append(resp.Errors, te.Error())
	}

	p.ResponseWriter().Header().Set("Content-Type", "application/json; charset=utf-8")
	p.ResponseWriter().WriteHeader(http.StatusBadRequest)
	srv.respondJSON(p.ResponseWriter(), resp)
}

func (srv *Server) HandleResponseProblem(p oas.Problem) {
	log.Printf("[WARN] oas problem on request %s %s: %v", p.Request().Method, p.Request().URL.String(), p.Cause())
}

type errorResponse struct {
	Errors []string `json:"errors"`
}

func (srv *Server) Health(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "I am alive!")
}
