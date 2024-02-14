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
	Recipient      ResponseRecipient      `json:"recipient"`
	ResponseAction ResponseResponseAction `json:"attachment"`
}

// id: sender_id string
type ResponseRecipient struct {
	ID string `json:"id"`
}
