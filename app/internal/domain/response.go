package domain

// SendMsg - sending message with parameters
type SendMsg[Kb any] struct {
	Text     string
	Keyboard Kb
}

// InKb - inline keyboard
type InKb [][]InKbButton

// InKbButton - inline keyboard button
type InKbButton struct {
	Text         string
	Url          string
	CallbackData string
}

// StKb - standard keyboard
type StKb [][]StKbButton

// StKbButton - standard keyboard button
type StKbButton struct {
	Text      string
	WebAppURL string
}

// CliMsg - client message and error.
// If need log an error, but for client send another message
type CliMsg struct {
	Message string
	Err     error
}

func (e *CliMsg) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}
