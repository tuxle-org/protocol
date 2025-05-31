package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
)

type Letter struct {
	Header map[string]string
	Body   string
}

func NewLetter() Letter {
	return Letter{
		Header: map[string]string{},
		Body:   "",
	}
}

func (letter Letter) Type() string {
	value, ok := letter.Header["type"]
	if !ok {
		return M_UNSPECIFIED
	}
	return value
}

func (letter Letter) ensureContainsHeader(header string) error {
	_, ok := letter.Header[header]
	if !ok {
		return ErrMissingHeader{
			Header: header,
		}
	}
	return nil
}

func (letter Letter) ensureBodyNotEmpty() error {
	if len(letter.Body) == 0 {
		return ErrBodyIsEmpty{}
	}
	return nil
}

// Ensures that all data is correct for the letter type
func (letter Letter) Validate() error {
	switch letter.Type() {

	case M_ERR:
		return errors.Join(
			letter.ensureContainsHeader("subject"),
			letter.ensureBodyNotEmpty(),
		)

	case M_MESSAGE:
		return letter.ensureBodyNotEmpty()

	case M_AUTH:
		err := letter.ensureContainsHeader("operation")
		if err != nil {
			return err
		}
		switch letter.Header["operation"] {
		case AUTH_CREATE:
		case AUTH_DELETE:
		case AUTH_LOGIN:
			return errors.Join(
				letter.ensureContainsHeader("user_id"),
				letter.ensureContainsHeader("password"),
			)
		case AUTH_LOGOUT:
		case AUTH_MODIFY:
		default:
			return ErrInvalidHeader{
				Header: "operation",
				Value:  letter.Header["operation"],
			}
		}
	}

	return nil
}

// Writes the header, all keys are in orbitrary order
//
// WARN: Do not use for tests
func (letter Letter) WriteHeader(writer io.Writer) error {
	for key, value := range letter.Header {
		_, err := fmt.Fprintf(writer, "%s=%s\n", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// Writes the header but with all the keys sorted
func (letter Letter) WriteHeaderSorted(writer io.Writer) error {
	keys := make([]string, 0, len(letter.Header))
	for key := range letter.Header {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		_, err := fmt.Fprintf(writer, "%s=%s\n", key, letter.Header[key])
		if err != nil {
			return err
		}
	}

	return nil
}

func (letter Letter) WriteBody(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "\n%s\x00", letter.Body)
	return err
}

// Format:
//
// [header1]=[value1]\n
// [header2]=[value2]\n
// [headerN]=[valueN]\n
// \n
// [body]\x00
func ReadLetter(reader *bufio.Reader) (Letter, error) {
	letter := NewLetter()
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
			return letter, ErrReadingHeaderKey{
				Err:   err,
				Index: index,
			}
		}

		headerValue, err := reader.ReadString('\n')
		if err != nil {
			return letter, ErrReadingHeaderValue{
				Err:   err,
				Index: index,
			}
		}

		if headerKey == "=" {
			return letter, ErrHeaderKeyIsEmpty{
				Index: index,
			}
		}

		letter.Header[headerKey[:len(headerKey)-1]] = headerValue[:len(headerValue)-1]
		index++
	}

	body, err := reader.ReadString('\x00')
	if err != nil {
		return letter, ErrReadingBody{
			Err: err,
		}
	}
	letter.Body = body[:len(body)-1]

	return letter, nil
}
