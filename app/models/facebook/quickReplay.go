package facebook

type QuickReplyMessage struct {
	Text         string               `json:"text"`
	QuickReplies []QuickReplyResponse `json:"quick_replies"`
}

type QuickReplyResponse struct {
	// text
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Payload     string `json:"payload"`
	ImageURL    string `json:"image_url"`
}
