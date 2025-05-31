package protocol

import (
	"fmt"
)

type ErrReadingHeaderKey struct {
	Err   error
	Index int
}

func (err ErrReadingHeaderKey) Error() string {
	return fmt.Sprintf("reading header:%d key: %s", err.Index, err.Err.Error())
}

// . . . . . . . .

type ErrReadingHeaderValue struct {
	Err   error
	Index int
}

func (err ErrReadingHeaderValue) Error() string {
	return fmt.Sprintf("reading header:%d value: %s", err.Index, err.Err.Error())
}

// . . . . . . . .

type ErrHeaderKeyIsEmpty struct {
	Index int
}

func (err ErrHeaderKeyIsEmpty) Error() string {
	return fmt.Sprintf("header:%d key is empty", err.Index)
}

// . . . . . . . .

type ErrReadingBody struct {
	Err error
}

func (err ErrReadingBody) Error() string {
	return fmt.Sprintf("reading body: %s", err.Err.Error())
}

// . . . . . . . .

type ErrMissingHeader struct {
	Header string
}

func (err ErrMissingHeader) Error() string {
	return fmt.Sprintf("missing header: %s", err.Header)
}

// . . . . . . . .

type ErrBodyIsEmpty struct{}

func (err ErrBodyIsEmpty) Error() string {
	return fmt.Sprintf("letter body is empty")
}

// . . . . . . . .

type ErrInvalidHeader struct {
	Header string
	Value  string
}

func (err ErrInvalidHeader) Error() string {
	return fmt.Sprintf(
		"letter header %q has invalid value of %q",
		err.Header,
		err.Value,
	)
}
