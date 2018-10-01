package oas

import (
	"strings"
)

// MultiError describes an error that wraps multiple errors that share
// the common message.
type MultiError interface {
	Message() string
	Errors() []error
}

func newMultiError(msg string, errs ...error) multiError {
	return multiError{
		msg:  msg,
		errs: errs,
	}
}

type multiError struct {
	msg  string
	errs []error
}

// Error implements error.
func (me multiError) Error() string {
	var ss []string
	for _, err := range me.errs {
		ss = append(ss, err.Error())
	}
	s := strings.Join(ss, ", ")
	if me.msg != "" {
		s = me.msg + ": " + s
	}
	return s
}

func (me multiError) Message() string {
	return me.msg
}

// Errors implements MultiError.
func (me multiError) Errors() []error {
	return me.errs
}
