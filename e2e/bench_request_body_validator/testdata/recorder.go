package testdata

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"

	"github.com/hypnoglow/oas2"
)

// RecorderSpec returns an OpenAPI spec for recorder server.
func RecorderSpec(t testing.TB) *oas.Document {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	doc, err := oas.LoadFile(path.Join(path.Dir(filename), "recorder.yaml"))
	require.NoError(t, err)
	return doc
}

// ---

func NewRecorderServer() *RecorderServer {
	return &RecorderServer{
		recordMessagesBatchHandler: recordMessagesBatchHandler{},
	}
}

type RecorderServer struct {
	recordMessagesBatchHandler http.Handler
}

func (s RecorderServer) RecordMessagesBatch() http.Handler {
	return s.recordMessagesBatchHandler
}

// RecorderHandler is a simple handler that greets using a name.
type recordMessagesBatchHandler struct{}

// ServeHTTP implements http.Handler.
func (recordMessagesBatchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body Body
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if err = json.Unmarshal(b, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(w, `{"data":{"info":"recorded %d messages"}`, len(body.Data.Messages))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Body struct {
	Data Data `json:"data"`
}

type Data struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Text string      `json:"text"`
	Time time.Time   `json:"time"`
	Kind MessageKind `json:"kind"`
}

type MessageKind string

const (
	MessageKindImportant MessageKind = "important"
	MessageKindRegular   MessageKind = "regular"
	MessageKindSpam      MessageKind = "spam"
)

// ---

// TestGreeter tests greeter server.
func TestGreeter(t *testing.T, srv *httptest.Server) {
	t.Helper()

	resp, err := srv.Client().Get(srv.URL + "/api/greeting?name=Foo")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected response, status=%s, body=%q", resp.Status, string(b))
	}

	expected := `{"greeting":"Hello, Foo!"}`
	if string(b) != expected {
		t.Fatalf("Expected %q but got %q", expected, string(b))
	}
}
