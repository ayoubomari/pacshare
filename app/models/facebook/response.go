package facebook

// facebook message as a response
type ResponseWithMessage struct {
	Recipient ResponseRecipient `json:"recipient"`
	Message   ResponseMessage   `json:"message"`
}

// facebook attachment as a response
type ResponseWithAttachment struct {
	Recipient  ResponseRecipient  `json:"recipient"`
	Attachment ResponseAttachment `json:"attachment"`
}

// id: sender_id string
type ResponseRecipient struct {
	ID string `json:"id"`
}
