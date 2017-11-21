package oas

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/go-openapi/spec"
)

func TestDecodeForm(t *testing.T) {
	t.Run("ok for simple form", func(t *testing.T) {
		params := []spec.Parameter{
			*spec.QueryParam("name").Typed("string", "").WithDefault("Anonymous"),
			*spec.QueryParam("region_id").Typed("integer", "int32"),
			*spec.QueryParam("tags").CollectionOf(spec.NewItems().Typed("string", ""), ""),
		}

		values := map[string][]string{
			"region_id": {"44"},
			"tags":      {"apples", "bananas"},
		}
		req := makeFormRequest(values)

		type member struct {
			Name     string   `oas:"name"`
			RegionID int32    `oas:"region_id"`
			Tags     []string `oas:"tags"`
		}

		var m member
		if err := DecodeForm(params, req, &m); err != nil {
			panic(err)
		}

		expected := member{
			Name:     "Anonymous",
			RegionID: 44,
			Tags:     []string{"apples", "bananas"},
		}
		if !reflect.DeepEqual(expected, m) {
			t.Fatalf("Expected %v but got %v", expected, m)
		}
	})

	t.Run("ok for multipart form", func(t *testing.T) {
		params := []spec.Parameter{
			*spec.QueryParam("name").Typed("string", "").WithDefault("Anonymous"),
			*spec.QueryParam("region_id").Typed("integer", "int32"),
			*spec.QueryParam("tags").CollectionOf(spec.NewItems().Typed("string", ""), ""),
		}

		values := map[string][]string{
			"region_id": {"44"},
			"tags":      {"apples", "bananas"},
		}
		req := makeMultipartRequest(values, nil)

		type member struct {
			Name     string   `oas:"name"`
			RegionID int32    `oas:"region_id"`
			Tags     []string `oas:"tags"`
		}

		var m member
		if err := DecodeForm(params, req, &m); err != nil {
			panic(err)
		}

		expected := member{
			Name:     "Anonymous",
			RegionID: 44,
			Tags:     []string{"apples", "bananas"},
		}
		if !reflect.DeepEqual(expected, m) {
			t.Fatalf("Expected %v but got %v", expected, m)
		}
	})

	t.Run("ok for file", func(t *testing.T) {
		params := []spec.Parameter{
			*spec.QueryParam("file").Typed("file", ""),
		}

		files := []multipartFile{
			{
				FieldName: "file",
				FileName:  "test.txt",
				Content:   strings.NewReader("It works!"),
			},
		}
		req := makeMultipartRequest(nil, files)

		type somethingWithFile struct {
			File *multipart.FileHeader `oas:"file"`
		}

		var some somethingWithFile
		if err := DecodeForm(params, req, &some); err != nil {
			panic(err)
		}

		f, err := some.File.Open()
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		expected := "It works!"
		if string(b) != expected {
			t.Fatalf("Expected file content to be %v but got %v", expected, string(b))
		}
	})
}

func makeMultipartRequest(values map[string][]string, files []multipartFile) *http.Request {
	buf := &bytes.Buffer{}

	form := multipart.NewWriter(buf)
	var err error

	// add values
	for k, vals := range values {
		for _, val := range vals {
			err = form.WriteField(k, val)
			if err != nil {
				panic(err)
			}
		}
	}

	// add files
	for _, mf := range files {
		f, err := form.CreateFormFile(mf.FieldName, mf.FileName)
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(f, mf.Content); err != nil {
			panic(err)
		}
	}

	// finalize
	err = form.Close()
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/members", buf)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", form.FormDataContentType())

	return req
}

func makeFormRequest(values map[string][]string) *http.Request {
	form := url.Values(values)

	req, err := http.NewRequest(http.MethodPost, "/members", strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req
}

func ExampleDecodeForm() {
	// In real app parameters will be taken from spec document (yaml or json).
	params := []spec.Parameter{
		*spec.QueryParam("name").Typed("string", ""),
		*spec.QueryParam("region_id").Typed("integer", "int32"),
	}

	// TODO
	buf := &bytes.Buffer{}
	form := multipart.NewWriter(buf)
	form.WriteField("name", "John")
	form.WriteField("region_id", "44")
	form.Close()
	req, err := http.NewRequest(http.MethodPost, "/members", buf)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", form.FormDataContentType())

	type member struct {
		Name     string `oas:"name"`
		RegionID int32  `oas:"region_id"`
	}

	var m member
	if err := DecodeForm(params, req, &m); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)

	// Output:
	// oas.member{Name:"John", RegionID:44}
}

func makeMultipartBody() *http.Request {
	buf := &bytes.Buffer{}
	form := multipart.NewWriter(buf)
	form.WriteField("name", "John")
	form.WriteField("region_id", "44")
	form.Close()
	//f, err := form.CreateFormFile("file", "test.txt")
	//if err != nil {
	//	panic(err)
	//}
	//if _, err := f.Write([]byte("It works!")); err != nil {
	//	panic(err)
	//}
	//if err := form.Close(); err != nil {
	//	panic(err)
	//}

	req, err := http.NewRequest(http.MethodPost, "/members", buf)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", form.FormDataContentType())
	//req.ContentLength = emptyMultipartSize(fieldname, filename) + size

	return req
}

func emptyMultipartSize(fieldname, filename string) int64 {
	body := &bytes.Buffer{}
	form := multipart.NewWriter(body)
	form.CreateFormFile(fieldname, filename)
	form.Close()
	return int64(body.Len())
}

type multipartFile struct {
	FieldName string
	FileName  string
	Content   io.Reader
}
