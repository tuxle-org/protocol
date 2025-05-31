package protocol

import (
	"fmt"
)

type ErrReadingParamKey struct {
	Err   error
	Index int
}

func (err ErrReadingParamKey) Error() string {
	return fmt.Sprintf("reading param:%d key: %s", err.Index, err.Err.Error())
}

// . . . . . . . .

type ErrReadingParamValue struct {
	Err   error
	Index int
}

func (err ErrReadingParamValue) Error() string {
	return fmt.Sprintf("reading param:%d value: %s", err.Index, err.Err.Error())
}

// . . . . . . . .

type ErrParamKeyIsEmpty struct {
	Index int
}

func (err ErrParamKeyIsEmpty) Error() string {
	return fmt.Sprintf("param:%d key is empty", err.Index)
}

// . . . . . . . .

type ErrReadingBody struct {
	Err error
}

func (err ErrReadingBody) Error() string {
	return fmt.Sprintf("reading body: %s", err.Err.Error())
}

// . . . . . . . .

type ErrMissingParam struct {
	Param string
}

func (err ErrMissingParam) Error() string {
	return fmt.Sprintf("missing param: %s", err.Param)
}

// . . . . . . . .

type ErrBodyIsEmpty struct{}

func (err ErrBodyIsEmpty) Error() string {
	return fmt.Sprintf("letter body is empty")
}

// . . . . . . . .

type ErrInvalidVariant struct {
	Kind  string
	Value string
}

func (err ErrInvalidVariant) Error() string {
	return fmt.Sprintf(
		"letter %q subtype has invalid value of %q",
		err.Kind,
		err.Value,
	)
}

// . . . . . . . .

type ErrInvalidFormat struct {
	Err error
}

func (err ErrInvalidFormat) Error() string {
	return "provided input is not a valid letter: " + err.Err.Error()
}
