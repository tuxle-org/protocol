package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Letter struct {
	Kind   LetterKind
	Params map[string]string
	Body   string
}

func NewLetter(kind LetterKind) Letter {
	return Letter{
		Kind:   kind,
		Params: map[string]string{},
		Body:   "",
	}
}

func (letter Letter) ensureContainsParam(param string) error {
	_, ok := letter.Params[param]
	if !ok {
		return ErrMissingParam{
			Param: param,
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

func (letter Letter) ParseKind(in string) (Letter, error) {
	parts := strings.Split(in, ".")

	var kind LetterKind
	switch parts[0] {
	case "error":
		if len(parts) != 2 {
			return letter, errors.New("Missing letter kind variant")
		}
		kind = ErrorKind{
			Value: ErrKindVariant(parts[1]),
		}
	case "user":
		if len(parts) != 2 {
			return letter, errors.New("Missing letter kind variant")
		}
		kind = UserKind{
			Value: UserKindVariant(parts[1]),
		}
	case "reciept":
		kind = RecieptKind{}
	case "message":
		kind = MessageKind{}
	case "undefined":
		kind = UndefinedKind{}
	}

	letter.Kind = kind
	return letter, nil
}

// Ensures that all data is correct for the letter type
func (letter Letter) Validate() error {
	switch kind := letter.Kind.(type) {

	case ErrorKind:
		switch kind.Variant() {
		case ERR_INTERNAL:
		case ERR_DENIED:
		case ERR_REQUEST:
		case ERR_TIMEOUT:
		case ERR_CUSTOM:
		default:
			return ErrInvalidVariant{
				Kind:  kind.Kind(),
				Value: string(kind.Variant()),
			}
		}
		return letter.ensureBodyNotEmpty()

	case MessageKind:
		return letter.ensureBodyNotEmpty()

	case UserKind:
		switch kind.Variant() {
		case USER_CREATE:
		case USER_DELETE:
		case USER_LOGIN:
			return errors.Join(
				letter.ensureContainsParam("user_id"),
				letter.ensureContainsParam("password"),
			)
		case USER_LOGOUT:
		case USER_MODIFY:
		case USER_QUERY:
		case USER_BAN:
		case USER_UNBAN:
		default:
			return ErrInvalidVariant{
				Kind:  kind.Kind(),
				Value: string(kind.Variant()),
			}
		}

	case RecieptKind:

	default:
		return fmt.Errorf("unexpected protocol.LetterKind: %#v", kind)
	}

	return nil
}

// =error.internal
// param1=value1
//
// body

func (letter Letter) WriteHeader(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, "=%s\n", letter.Kind.String())
	return err
}

// Writes the params, all keys are in orbitrary order
//
// WARN: Do not use for tests
func (letter Letter) WriteParams(writer io.Writer) error {
	for key, value := range letter.Params {
		_, err := fmt.Fprintf(writer, "%s=%s\n", key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// Writes the params but with all the keys sorted
func (letter Letter) WriteParamsSorted(writer io.Writer) error {
	keys := make([]string, 0, len(letter.Params))
	for key := range letter.Params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		_, err := fmt.Fprintf(writer, "%s=%s\n", key, letter.Params[key])
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

// Writes the header, params (orbitrary order) and body
func (letter Letter) Write(writer io.Writer) error {
	return errors.Join(
		letter.WriteParams(writer),
		letter.WriteBody(writer),
	)
}

// Format:
//
// [param1]=[value1]\n
// [param2]=[value2]\n
// [paramN]=[valueN]\n
// \n
// [body]\x00
func ReadLetter(reader *bufio.Reader) (Letter, error) {
	letter := NewLetter(UndefinedKind{})

	byte, err := reader.ReadByte()
	if byte != '=' {
		return letter, ErrInvalidFormat{}
	}

	header, err := reader.ReadString('\n')
	if err != nil {
		return letter, ErrInvalidFormat{Err: err}
	}

	letter, err = letter.ParseKind(header[:len(header)-1])
	if err != nil {
		return letter, ErrInvalidFormat{Err: err}
	}

	index := 1
	for {
		// Check for end-of-params sequence "\n\n"
		peek, err := reader.Peek(1)
		if err == nil && peek[0] == '\n' {
			reader.Discard(1)
			break
		}

		// Read param name until ':'
		paramKey, err := reader.ReadString('=')
		if err != nil {
			return letter, ErrReadingParamKey{
				Err:   err,
				Index: index,
			}
		}

		paramValue, err := reader.ReadString('\n')
		if err != nil {
			return letter, ErrReadingParamValue{
				Err:   err,
				Index: index,
			}
		}

		if paramKey == "=" {
			return letter, ErrParamKeyIsEmpty{
				Index: index,
			}
		}

		letter.Params[paramKey[:len(paramKey)-1]] = paramValue[:len(paramValue)-1]
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
