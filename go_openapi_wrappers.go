package oas

import (
	"github.com/go-openapi/loads"
)

// This file contains wrappers around go-openapi/* packages exported types to
// not to force oas package users explicitly import go-openapi packages.

// Document represents a swagger spec document.
type Document struct {
	*loads.Document
}

func wrapDocument(doc *loads.Document) *Document {
	return &Document{Document: doc}
}
