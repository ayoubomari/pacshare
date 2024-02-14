package facebook

// the body of response action
type ResponseResponseAction struct {
	// MARK_SEEN | TYPING_ON | TYPING_OFF | REACT | UNREACT
	Sender_action string `json:"sender_action,omitempty"`
}
