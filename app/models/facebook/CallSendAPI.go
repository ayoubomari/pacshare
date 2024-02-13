package facebook

// the back response that you got after sending a message to facebook
type CallSendAPIResonse struct {
	Recipient_id string `json:"Recipient_id,omitempty"`
	Message_id   string `json:"message_id,omitempty"`
	Error        struct {
		Message string `json:"message,omitempty"`
	} `json:"error,omitempty"`
}
