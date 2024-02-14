package facebook

// the body of response attachment
type TemplateMessage struct {
	Attachment ResponseTemplateAttachment `json:"attachment"`
}

type ResponseTemplateAttachment struct {
	// image | audio | video | file
	Type    string                    `json:"type"`
	Payload TemplateAttachmentPayload `json:"payload"`
}

type TemplateAttachmentPayload struct {
	TemplateType string                      `json:"template_type"`
	Elements     []TemplateAttachmentElement `json:"elements,omitempty"`
	Text         string                      `json:"text,omitempty"`
	Buttons      []TemplateButtonButton      `json:"buttons,omitempty"`
}

type TemplateAttachmentElement struct {
	Title    string                     `json:"title"`
	Subtitle string                     `json:"subtitle"`
	ImageURL string                     `json:"image_url"`
	Buttons  []TemplateAttachmentButton `json:"buttons"`
}

type TemplateAttachmentButton struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Payload string `json:"payload"`
}

type TemplateButtonButton struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Payload string `json:"payload"`
}
