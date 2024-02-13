package facebook

type ResponseAttachment struct {
	Type    string            `json:"type"`
	Payload AttachmentPayload `json:"payload"`
}

type AttachmentPayload struct {
	TemplateType string    `json:"template_type"`
	Elements     []Element `json:"elements"`
}

type Element struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	ImageURL string   `json:"image_url"`
	Buttons  []Button `json:"buttons"`
}

type Button struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Payload string `json:"payload"`
}
