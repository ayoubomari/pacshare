package facebook

// facebook message as a response
type ResponseWithMessage struct {
	Recipient ResponseRecipient `json:"recipient"`
	Message   ResponseMessage   `json:"message,omitempty"`
}

// facebook media attachment as a response
type ResponseWithMediaAttachment struct {
	Recipient ResponseRecipient `json:"recipient"`
	Message   MediaMessage      `json:"message"`
}

// facebook template attachment as a response
type ResponseWithTemplateAttachment struct {
	Recipient ResponseRecipient `json:"recipient"`
	Message   TemplateMessage   `json:"message"`
}

// facebook action as a response
type ResponseWithResponseAction struct {
	Recipient ResponseRecipient `json:"recipient"`
	// MARK_SEEN | TYPING_ON | TYPING_OFF | REACT | UNREACT
	Sender_action string `json:"sender_action"`
}

// facebook quick replay as a response
type ResponseWithQuickReplay struct {
	Recipient ResponseRecipient `json:"recipient"`
	Message   QuickReplyMessage `json:"message"`
}

// id: sender_id string
type ResponseRecipient struct {
	ID string `json:"id"`
}
