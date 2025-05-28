package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
)

type Message struct {
	Header map[string]string
	Body   string
}

func NewMessage() Message {
	return Message{
		Header: map[string]string{},
		Body:   "",
	}
}

func (message Message) Type() string {
	value, ok := message.Header["type"]
	if !ok {
		return M_UNSPECIFIED
	}
	return value
}

func (message Message) ensureContainsHeader(header string) error {
	_, ok := message.Header[header]
	if !ok {
		return ErrMissingHeader{
			Header: header,
		}
	}
	return nil
}

func (message Message) ensureBodyNotEmpty() error {
	if len(message.Body) == 0 {
		return ErrBodyIsEmpty{}
	}
	return nil
}

// Ensures that all data is correct for the message type
func (message Message) Validate() error {
	switch message.Type() {

	case M_ERROR:
		return errors.Join(
			message.ensureContainsHeader("subject"),
			message.ensureBodyNotEmpty(),
		)

	case M_MESSAGE:
		return message.ensureBodyNotEmpty()
	}

	return nil
}

// Writes the header, all keys are in orbitrary order
//
// WARN: Do not use for tests
func (message Message) WriteHeader(writer io.Writer) error {
	for key, value := range message.Header {
		_, err := fmt.Fprintf(writer, "%s=%s\n", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// Writes the header but with all the keys sorted
func (message Message) WriteHeaderSorted(writer io.Writer) error {
	keys := make([]string, 0, len(message.Header))
	for key := range message.Header {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		_, err := fmt.Fprintf(writer, "%s=%s\n", key, message.Header[key])
		if err != nil {
			return err
		}
	}

	return nil
}

func (message Message) WriteBody(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "\n%s\x00", message.Body)
	return err
}

// Format:
//
// [header1]=[value1]\n
// [header2]=[value2]\n
// [headerN]=[valueN]\n
// \n
// [body]\x00
func ReadMessage(reader *bufio.Reader) (Message, error) {
	message := NewMessage()
	index := 1

	for {
		// Check for end-of-headers sequence "\n\n"
		peek, err := reader.Peek(1)
		if err == nil && peek[0] == '\n' {
			reader.Discard(1)
			break
		}

		// Read header name until ':'
		headerKey, err := reader.ReadString('=')
		if err != nil {
			return message, ErrReadingHeaderKey{
				Err:   err,
				Index: index,
			}
		}

		headerValue, err := reader.ReadString('\n')
		if err != nil {
			return message, ErrReadingHeaderValue{
				Err:   err,
				Index: index,
			}
		}

		if headerKey == "=" {
			return message, ErrHeaderKeyIsEmpty{
				Index: index,
			}
		}

		message.Header[headerKey[:len(headerKey)-1]] = headerValue[:len(headerValue)-1]
		index++
	}

	body, err := reader.ReadString('\x00')
	if err != nil {
		return message, ErrReadingBody{
			Err: err,
		}
	}
	message.Body = body[:len(body)-1]

	return message, nil
}
