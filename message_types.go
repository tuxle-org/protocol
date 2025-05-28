package protocol

const (
	M_UNSPECIFIED = "unspecified"
	M_MESSAGE     = "message"
	M_RECEIPT     = "reciept"
	M_ERROR       = "error"
)

func NewErrorMessage(subject string) Message {
	return Message{
		Header: map[string]string{
			"type":    M_ERROR,
			"subject": subject,
		},
		Body: "",
	}
}
