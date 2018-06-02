package oas

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDynamicSpecHandler(t *testing.T) {
	doc := loadDocFile(t, "testdata/petstore_1.yml")

	h := DynamicSpecHandler(doc.Spec())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	req.Header.Set("X-Forwarded-Host", "foo.bar.com")
	req.Header.Set("X-Forwarded-Proto", "https")

	h.ServeHTTP(rr, req)

	writtenDoc := loadDocBytes(rr.Body.Bytes())

	if writtenDoc.Spec().Host != "foo.bar.com" {
		t.Errorf("Expected host to be foo.bar.com but got %q", writtenDoc.Spec().Host)
	}

	if !reflect.DeepEqual(writtenDoc.Spec().Schemes, []string{"https"}) {
		t.Errorf("Expected schemes to be [https] but got %v", writtenDoc.Spec().Schemes)
	}

	// check that original spec fields remain same

	if doc.Spec().Host != "petstore.swagger.io" {
		t.Errorf("Expected original spec host hasn't changed but got %v", doc.Spec().Host)
	}

	if !reflect.DeepEqual(doc.Spec().Schemes, []string{"http"}) {
		t.Errorf("Expected original spec schemes hasn't changed but got %v", doc.Spec().Schemes)
	}
}

func TestStaticSpecHandler(t *testing.T) {
	doc := loadDocFile(t, "testdata/petstore_1.yml")

	doc.Spec().Info.Version = "1.2.3"
	doc.Spec().Host = "foo.bar.com"
	doc.Spec().Schemes = []string{"https"}

	h := StaticSpecHandler(doc.Spec())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/foo", nil)

	h.ServeHTTP(rr, req)

	writtenDoc := loadDocBytes(rr.Body.Bytes())

	if writtenDoc.Spec().Info.Version != "1.2.3" {
		t.Errorf("Expected version to be 1.2.3 but got %q", writtenDoc.Spec().Info.Version)
	}

	if writtenDoc.Spec().Host != "foo.bar.com" {
		t.Errorf("Expected host to be foo.bar.com but got %q", writtenDoc.Spec().Host)
	}

	if !reflect.DeepEqual(writtenDoc.Spec().Schemes, []string{"https"}) {
		t.Errorf("Expected schemes to be [https] but got %v", writtenDoc.Spec().Schemes)
	}
}
