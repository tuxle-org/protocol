package protocol

type LetterKind interface {
	Kind() string
	String() string
}

// . . . . . . . .

type ErrKindVariant string

const (
	ERR_INTERNAL ErrKindVariant = "internal"
	ERR_DENIED   ErrKindVariant = "access_denied"
	ERR_REQUEST  ErrKindVariant = "bad_request"
	ERR_TIMEOUT  ErrKindVariant = "timed_out"
	ERR_CUSTOM   ErrKindVariant = "custom"
)

type ErrorKind struct {
	Value ErrKindVariant
}

func (kind ErrorKind) Kind() string {
	return "error"
}

func (kind ErrorKind) Variant() ErrKindVariant {
	return kind.Value
}

func (kind ErrorKind) String() string {
	return kind.Kind() + "/" + string(kind.Variant())
}

// . . . . . . . .

type RecieptKind struct{}

func (kind RecieptKind) Kind() string {
	return "reciept"
}

func (kind RecieptKind) String() string {
	return kind.Kind()
}

// . . . . . . . .

type MessageKind struct{}

func (kind MessageKind) Kind() string {
	return "message"
}

func (kind MessageKind) String() string {
	return kind.Kind()
}

// . . . . . . . .

type UserKindVariant string

const (
	USER_LOGIN  UserKindVariant = "log_in"
	USER_LOGOUT UserKindVariant = "log_out"
	USER_CREATE UserKindVariant = "create"
	USER_MODIFY UserKindVariant = "modify"
	USER_DELETE UserKindVariant = "delete"
	USER_QUERY  UserKindVariant = "query"
	USER_BAN    UserKindVariant = "ban"
	USER_UNBAN  UserKindVariant = "unban"
)

type UserKind struct {
	Value UserKindVariant
}

func (kind UserKind) Kind() string {
	return "user"
}

func (kind UserKind) Variant() UserKindVariant {
	return kind.Value
}

func (kind UserKind) String() string {
	return kind.Kind() + "/" + string(kind.Variant())
}

// . . . . . . . .

type UndefinedKind struct{}

func (kind UndefinedKind) Kind() string {
	return "undefined"
}

func (kind UndefinedKind) String() string {
	return kind.Kind()
}
