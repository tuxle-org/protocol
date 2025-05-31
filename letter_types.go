package protocol

const (
	M_UNSPECIFIED = "unspecified"
	M_MESSAGE     = "message"
	M_RECEIPT     = "reciept"
	M_ERR         = "error"
	M_AUTH        = "auth"
)

const (
	ERR_INTERNAL = "internal"
	ERR_DENIED   = "permission_denied"
)

const (
	AUTH_LOGIN  = "log_in"
	AUTH_LOGOUT = "log_out"
	AUTH_CREATE = "create"
	AUTH_DELETE = "delete"
	AUTH_MODIFY = "modify"
)

func NewErrorLetter(subject string) Letter {
	return Letter{
		Header: map[string]string{
			"type":    M_ERR,
			"subject": subject,
		},
		Body: "",
	}
}
