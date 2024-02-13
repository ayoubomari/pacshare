package facebook

// message webhook (when some one send message to the page)
type FacebookWebhookBody struct {
	Object string      `json:"object"`
	Entry  []PageEntry `json:"entry"`
}

type PageEntry struct {
	ID        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp"`
	Message   Message   `json:"message,omitempty"`
	PostBack  PostBack  `json:"postback,omitempty"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type Message struct {
	MID         string       `json:"mid"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Type    string                       `json:"type"`
	Payload WebhookBodyAttachmentPayload `json:"payload"`
}

type WebhookBodyAttachmentPayload struct {
	URL         string                           `json:"url,omitempty"`
	Is_reusable bool                             `json:"is_reusable,omitempty"`
	Title       string                           `json:"title,omitempty"`
	Coordinates WebhookBodyAttachmentCoordinates `json:"coordinates,omitempty"`
}

type WebhookBodyAttachmentCoordinates struct {
	Lat  float64 `json:"lat,omitempty"`
	Long float64 `json:"long,omitempty"`
}

type PostBack struct {
	Title   string `json:"title"`
	Payload string `json:"payload"`
}
