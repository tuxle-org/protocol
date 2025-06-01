package protocol_test

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/tuxle-org/protocol"
	"gotest.tools/assert"
)

var (
	TimeUnix = time.Now().UnixMilli()
	Time     = time.UnixMilli(TimeUnix)
)

func TestReadLetter(test *testing.T) {
	const header1 = "header_name_1=some header value 1\n"
	const header2 = "header_name_2=some header value 2\n"
	const header3 = "header_name_3=\n"
	const body = "# Title\nThis is example content of the letter here."
	testCases := []struct {
		input     string
		expectErr error
		expectMsg protocol.Letter
	}{
		{
			input: fmt.Sprintf(
				"=message @%d\n%s%s%s\n%s\x00",
				TimeUnix,
				header1,
				header2,
				header3,
				body,
			),
			expectErr: nil,
			expectMsg: protocol.Letter{
				Kind:      protocol.MessageKind{},
				Timestamp: Time,
				Params: map[string]string{
					"header_name_1": "some header value 1",
					"header_name_2": "some header value 2",
					"header_name_3": "",
				},
				Body: body,
			},
		},
		{
			input:     fmt.Sprintf("=message @%d\n\n%s\x00", TimeUnix, body),
			expectErr: nil,
			expectMsg: protocol.Letter{
				Kind:      protocol.MessageKind{},
				Timestamp: Time,
				Params:    map[string]string{},
				Body:      "# Title\nThis is example content of the letter here.",
			},
		},
		{
			input:     fmt.Sprintf("=message @%d\n%s%s%sinvalid_header\n%s\x00", TimeUnix, header1, header2, header3, body),
			expectErr: protocol.ErrReadingParamKey{},
			expectMsg: protocol.Letter{},
		},
		{
			input:     fmt.Sprintf("=message @%d\n%s%s%s=empty key name\n%s\x00", TimeUnix, header1, header2, header3, body),
			expectErr: protocol.ErrParamKeyIsEmpty{},
			expectMsg: protocol.Letter{},
		},
		{
			input:     fmt.Sprintf("=message @%d\n%s%s%sinvalid_header=", TimeUnix, header1, header2, header3),
			expectErr: protocol.ErrReadingParamValue{},
			expectMsg: protocol.Letter{},
		},
		{
			input:     fmt.Sprintf("=message @%d\n%s%s%s\n%s", TimeUnix, header1, header2, header3, body),
			expectErr: protocol.ErrReadingBody{},
			expectMsg: protocol.Letter{},
		},
	}

	for _, testCase := range testCases {
		test.Run(
			fmt.Sprintf("expect_err=%v", reflect.TypeOf(testCase.expectErr)),
			func(test *testing.T) {
				reader := bufio.NewReader(strings.NewReader(testCase.input))
				letter, err := protocol.ReadLetter(reader)
				if testCase.expectErr != nil {
					assert.ErrorType(test, err, testCase.expectErr)
				} else {
					assert.NilError(test, err)
					assert.DeepEqual(test, letter, testCase.expectMsg)
				}
			},
		)
	}
}

func TestLetterWrite(test *testing.T) {
	const header1 = "header_name_1=some header value 1\n"
	const header2 = "header_name_2=some header value 2\n"
	const header3 = "header_name_3=\n"
	const body = "# Title\nThis is example content of the letter here."
	testCases := []struct {
		input        protocol.Letter
		expectErr    bool
		expectBuffer string
	}{
		{
			input: protocol.Letter{
				Kind:      protocol.MessageKind{},
				Timestamp: Time,
				Params: map[string]string{
					"header_name_1": "some header value 1",
					"header_name_2": "some header value 2",
					"header_name_3": "",
				},
				Body: body,
			},
			expectBuffer: fmt.Sprintf("=message @%d\n%s%s%s\n%s\x00", TimeUnix, header1, header2, header3, body),
		},
		{
			input: protocol.Letter{
				Kind:      protocol.MessageKind{},
				Timestamp: Time,
				Params:    map[string]string{},
				Body:      body,
			},
			expectBuffer: fmt.Sprintf("=message @%d\n\n%s\x00", TimeUnix, body),
		},
		{
			input: protocol.Letter{
				Kind:      protocol.MessageKind{},
				Timestamp: Time,
				Params: map[string]string{
					"header_name_1": "some header value 1",
					"header_name_2": "some header value 2",
					"header_name_3": "",
				},
				Body: "",
			},
			expectBuffer: fmt.Sprintf("=message @%d\n%s%s%s\n\x00", TimeUnix, header1, header2, header3),
		},
	}

	for _, testCase := range testCases {
		test.Run("", func(test *testing.T) {
			var buffer bytes.Buffer
			assert.NilError(test, testCase.input.WriteHeader(&buffer))
			assert.NilError(test, testCase.input.WriteParamsSorted(&buffer))
			assert.NilError(test, testCase.input.WriteBody(&buffer))
			assert.DeepEqual(test, buffer.String(), testCase.expectBuffer)
		})
	}
}

func TestLetterValidate(test *testing.T) {
	letter := protocol.Letter{
		Kind:      protocol.MessageKind{},
		Timestamp: Time,
		Params: map[string]string{
			"header_name_1": "some header value 1",
			"header_name_2": "some header value 2",
			"header_name_3": "",
		},
		Body: "# Title\nThis is example content of the letter here.",
	}
	assert.NilError(test, letter.Validate())
}

func TestLetterValidateError(test *testing.T) {
	var letter protocol.Letter

	letter = protocol.Letter{
		Kind:      protocol.ErrorKind{},
		Timestamp: Time,
		Params:    map[string]string{},
		Body:      "123",
	}
	assert.ErrorType(test, letter.Validate(), protocol.ErrInvalidVariant{})

	letter = protocol.NewLetter(protocol.ErrorKind{Value: protocol.ERR_INTERNAL})
	assert.ErrorType(test, letter.Validate(), protocol.ErrBodyIsEmpty{})
}

func TestLetterValidateLetter(test *testing.T) {
	var letter protocol.Letter

	letter = protocol.Letter{
		Kind:      protocol.MessageKind{},
		Timestamp: Time,
		Params:    map[string]string{},
		Body:      "",
	}
	assert.ErrorType(test, letter.Validate(), protocol.ErrBodyIsEmpty{})
}
