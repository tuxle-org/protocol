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

func TestReadMessage(test *testing.T) {
	const header1 = "header_name_1=some header value 1\n"
	const header2 = "header_name_2=some header value 2\n"
	const header3 = "header_name_3=\n"
	const body = "# Title\nThis is example content of the message here."
	testCases := []struct {
		input     string
		expectErr error
		expectMsg protocol.Message
	}{
		{
			input:     fmt.Sprintf("%s%s%s\n%s\x00", header1, header2, header3, body),
			expectErr: nil,
			expectMsg: protocol.Message{
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
			expectMsg: protocol.Message{
				Header: map[string]string{},
				Body:   "# Title\nThis is example content of the message here.",
			},
		},
		{
			input:     fmt.Sprintf("%s%s%sinvalid_header\n%s\x00", header1, header2, header3, body),
			expectErr: protocol.ErrReadingHeaderKey{},
			expectMsg: protocol.Message{},
		},
		{
			input:     fmt.Sprintf("%s%s%s=empty key name\n%s\x00", header1, header2, header3, body),
			expectErr: protocol.ErrHeaderKeyIsEmpty{},
			expectMsg: protocol.Message{},
		},
		{
			input:     fmt.Sprintf("%s%s%sinvalid_header=", header1, header2, header3),
			expectErr: protocol.ErrReadingHeaderValue{},
			expectMsg: protocol.Message{},
		},
		{
			input:     fmt.Sprintf("%s%s%s\n%s", header1, header2, header3, body),
			expectErr: protocol.ErrReadingBody{},
			expectMsg: protocol.Message{},
		},
	}

	for _, testCase := range testCases {
		test.Run(
			fmt.Sprintf("expect_err=%v", reflect.TypeOf(testCase.expectErr)),
			func(test *testing.T) {
				reader := bufio.NewReader(strings.NewReader(testCase.input))
				message, err := protocol.ReadMessage(reader)
				if testCase.expectErr != nil {
					assert.ErrorType(test, err, testCase.expectErr)
				} else {
					assert.NilError(test, err)
					assert.DeepEqual(test, message, testCase.expectMsg)
				}
			},
		)
	}
}

func TestMessageWrite(test *testing.T) {
	const header1 = "header_name_1=some header value 1\n"
	const header2 = "header_name_2=some header value 2\n"
	const header3 = "header_name_3=\n"
	const body = "# Title\nThis is example content of the message here."
	testCases := []struct {
		input        protocol.Message
		expectErr    bool
		expectBuffer string
	}{
		{
			input: protocol.Message{
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
			input: protocol.Message{
				Header: map[string]string{},
				Body:   body,
			},
			expectBuffer: "\n" + body + "\x00",
		},
		{
			input: protocol.Message{
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

func TestMessageValidate(test *testing.T) {
	message := protocol.Message{
		Header: map[string]string{
			"header_name_1": "some header value 1",
			"header_name_2": "some header value 2",
			"header_name_3": "",
		},
		Body: "# Title\nThis is example content of the message here.",
	}
	assert.NilError(test, message.Validate())
}

func TestMessageValidateError(test *testing.T) {
	var message protocol.Message

	message = protocol.Message{
		Header: map[string]string{
			"type": "error",
		},
		Body: "123",
	}
	assert.ErrorType(test, message.Validate(), errors.Join(protocol.ErrMissingHeader{}))

	message = protocol.Message{
		Header: map[string]string{
			"type":    "error",
			"subject": "internal",
		},
		Body: "",
	}
	assert.ErrorType(test, message.Validate(), errors.Join(protocol.ErrBodyIsEmpty{}))

	message = protocol.Message{
		Header: map[string]string{
			"type": "error",
		},
		Body: "",
	}
	assert.ErrorType(test, message.Validate(), errors.Join(protocol.ErrMissingHeader{}, protocol.ErrBodyIsEmpty{}))
}

func TestMessageValidateMessage(test *testing.T) {
	var message protocol.Message

	message = protocol.Message{
		Header: map[string]string{
			"type": "message",
		},
		Body: "",
	}
	assert.ErrorType(test, message.Validate(), protocol.ErrBodyIsEmpty{})
}
