package bench_request_body_validator

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hypnoglow/oas2"
	"github.com/hypnoglow/oas2/adapter/gorilla"
	"github.com/hypnoglow/oas2/e2e/bench_request_body_validator/testdata"
)

// go test -bench=. -benchmem -benchtime 5s -memprofile mem.out -run BenchmarkRequestBodyValidator ./e2e/bench_request_body_validator
func BenchmarkRequestBodyValidator(b *testing.B) {
	router := setupServer(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/recorder/records/batch", makeBody(b))
		req.Header.Set("Content-Type", "application/json")
		b.StartTimer()

		router.ServeHTTP(w, req)
		require.Equal(b, w.Code, http.StatusOK)
	}
}

func setupServer(t testing.TB) http.Handler {
	t.Helper()

	recorderServer := testdata.NewRecorderServer()

	doc := testdata.RecorderSpec(t)
	basis := oas.NewResolvingBasis(doc, oas_gorilla.NewResolver(doc))

	router := mux.NewRouter()
	opRouter := oas_gorilla.NewOperationRouter(router)
	opRouter = opRouter.WithDocument(doc)
	opRouter = opRouter.WithMiddleware(basis.OperationContext())
	opRouter = opRouter.WithMiddleware(basis.RequestBodyValidator())
	opRouter = opRouter.WithOperationHandlers(map[string]http.Handler{
		"recordMessagesBatch": recorderServer.RecordMessagesBatch(),
	})
	assert.NoError(t, opRouter.Route())

	return router
}

func makeBody(t testing.TB) io.Reader {
	t.Helper()

	var body testdata.Body

	body.Data.Messages = make([]testdata.Message, 1000)
	for i := range body.Data.Messages {
		body.Data.Messages[i] = testdata.Message{
			Text: randomText(16),
			Time: time.Now(),
			Kind: randomKind(i),
		}
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	return bytes.NewReader(b)
}

func randomKind(i int) testdata.MessageKind {
	ki := i % 3
	switch ki {
	case 0:
		return testdata.MessageKindImportant
	case 1:
		return testdata.MessageKindRegular
	default:
		return testdata.MessageKindSpam
	}
}

// ---

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randomText(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

/*

Without validator:

goos: darwin
goarch: amd64
pkg: github.com/hypnoglow/oas2/e2e/bench_request_body_validator
BenchmarkRequestBodyValidator-8   	     300	  23878377 ns/op	 4883404 B/op	   30085 allocs/op
PASS

Current:

goos: darwin
goarch: amd64
pkg: github.com/hypnoglow/oas2/e2e/bench_request_body_validator
BenchmarkRequestBodyValidator-8   	      50	 321915067 ns/op	248219781 B/op	 2010492 allocs/op
PASS
ok  	github.com/hypnoglow/oas2/

(foo) without the validation (return nil):

goos: darwin
goarch: amd64
pkg: github.com/hypnoglow/oas2/e2e/bench_request_body_validator
BenchmarkRequestBodyValidator-8   	     200	  41850622 ns/op	14008269 B/op	  140407 allocs/op
PASS

With cache, version 1

goos: darwin
goarch: amd64
pkg: github.com/hypnoglow/oas2/e2e/bench_request_body_validator
BenchmarkRequestBodyValidator-8   	      50	 301242776 ns/op	248231011 B/op	 2010724 allocs/op
PASS

// -------------------------------------

pkg: github.com/hypnoglow/oas2/e2e/bench_request_body_validator
BenchmarkRequestBodyValidator-8   	     200	  42339679 ns/op	14008290 B/op	  140408 allocs/op
PASS

with iter:

pkg: github.com/hypnoglow/oas2/e2e/bench_request_body_validator
BenchmarkRequestBodyValidator-8   	     300	  26612612 ns/op	16675587 B/op	  210407 allocs/op
PASS

with validator (cached):

pkg: github.com/hypnoglow/oas2/e2e/bench_request_body_validator
BenchmarkRequestBodyValidator-8   	      20	 290009005 ns/op	250900839 B/op	 2080729 allocs/op
PASS

// ------------------------------------- 1,000 -------------------------------------

without cache:

pkg: github.com/hypnoglow/oas2/e2e/bench_request_body_validator
BenchmarkRequestBodyValidator-8   	     300	  27179841 ns/op	24484944 B/op	  208353 allocs/op
PASS

with cache:

pkg: github.com/hypnoglow/oas2/e2e/bench_request_body_validator
BenchmarkRequestBodyValidator-8   	     300	  23600512 ns/op	24487185 B/op	  208525 allocs/op
PASS

*/
