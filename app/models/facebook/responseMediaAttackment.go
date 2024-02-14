package facebook

// the body of response attachment
type MediaMessage struct {
	Attachment ResponseMediaAttachment `json:"attachment"`
}

type ResponseMediaAttachment struct {
	// image | audio | video | file
	Type    string                       `json:"type"`
	Payload WebhookBodyAttachmentPayload `json:"payload"`
}
