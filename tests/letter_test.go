package protocol_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/tuxle-org/protocol"
	"gotest.tools/assert"
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
			input:     fmt.Sprintf("%s%s%s\n%s\x00", header1, header2, header3, body),
			expectErr: nil,
			expectMsg: protocol.Letter{
				Header: map[string]string{
					"header_name_1": "some header value 1",
					"header_name_2": "some header value 2",
					"header_name_3": "",
				},
				Body: body,
			},
		},
		{
			input:     fmt.Sprintf("\n%s\x00", body),
			expectErr: nil,
			expectMsg: protocol.Letter{
				Header: map[string]string{},
				Body:   "# Title\nThis is example content of the letter here.",
			},
		},
		{
			input:     fmt.Sprintf("%s%s%sinvalid_header\n%s\x00", header1, header2, header3, body),
			expectErr: protocol.ErrReadingHeaderKey{},
			expectMsg: protocol.Letter{},
		},
		{
			input:     fmt.Sprintf("%s%s%s=empty key name\n%s\x00", header1, header2, header3, body),
			expectErr: protocol.ErrHeaderKeyIsEmpty{},
			expectMsg: protocol.Letter{},
		},
		{
			input:     fmt.Sprintf("%s%s%sinvalid_header=", header1, header2, header3),
			expectErr: protocol.ErrReadingHeaderValue{},
			expectMsg: protocol.Letter{},
		},
		{
			input:     fmt.Sprintf("%s%s%s\n%s", header1, header2, header3, body),
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
				Header: map[string]string{
					"header_name_1": "some header value 1",
					"header_name_2": "some header value 2",
					"header_name_3": "",
				},
				Body: body,
			},
			expectBuffer: header1 + header2 + header3 + "\n" + body + "\x00",
		},
		{
			input: protocol.Letter{
				Header: map[string]string{},
				Body:   body,
			},
			expectBuffer: "\n" + body + "\x00",
		},
		{
			input: protocol.Letter{
				Header: map[string]string{
					"header_name_1": "some header value 1",
					"header_name_2": "some header value 2",
					"header_name_3": "",
				},
				Body: "",
			},
			expectBuffer: header1 + header2 + header3 + "\n\x00",
		},
	}

	for _, testCase := range testCases {
		test.Run("", func(test *testing.T) {
			var buffer bytes.Buffer
			assert.NilError(test, testCase.input.WriteHeaderSorted(&buffer))
			assert.NilError(test, testCase.input.WriteBody(&buffer))
			assert.DeepEqual(test, buffer.String(), testCase.expectBuffer)
		})
	}
}

func TestLetterValidate(test *testing.T) {
	letter := protocol.Letter{
		Header: map[string]string{
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
		Header: map[string]string{
			"type": "error",
		},
		Body: "123",
	}
	assert.ErrorType(test, letter.Validate(), errors.Join(protocol.ErrMissingHeader{}))

	letter = protocol.NewErrorLetter(protocol.ERR_INTERNAL)
	assert.ErrorType(test, letter.Validate(), errors.Join(protocol.ErrBodyIsEmpty{}))

	letter = protocol.Letter{
		Header: map[string]string{
			"type": "error",
		},
		Body: "",
	}
	assert.ErrorType(test, letter.Validate(), errors.Join(protocol.ErrMissingHeader{}, protocol.ErrBodyIsEmpty{}))
}

func TestLetterValidateLetter(test *testing.T) {
	var letter protocol.Letter

	letter = protocol.Letter{
		Header: map[string]string{
			"type": "message",
		},
		Body: "",
	}
	assert.ErrorType(test, letter.Validate(), protocol.ErrBodyIsEmpty{})
}
