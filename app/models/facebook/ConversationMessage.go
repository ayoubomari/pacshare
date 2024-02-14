package facebook

type ConversationMessage struct {
	From        ConversationFrom        `json:"from,omitempty"`
	Message     string                  `json:"message,omitempty"`
	Attachments ConversationAttachments `json:"attachments,omitempty"`
	ID          string                  `json:"id,omitempty"`
}

type ConversationFrom struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	ID    string `json:"id"`
}

type ConversationAttachments struct {
	Data   []ConversationAttachmentData `json:"data"`
	Paging ConversationPaging           `json:"paging"`
}

type ConversationPaging struct {
	Cursors ConversationCursors `json:"cursors"`
	Next    string              `json:"next"`
}

type ConversationAttachmentData struct {
	ID        string                          `json:"id"`
	MimeType  string                          `json:"mime_type"`
	Name      string                          `json:"name"`
	Size      int                             `json:"size"`
	ImageData ConversationAttachmentImageData `json:"image_data"`
}

type ConversationAttachmentImageData struct {
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	MaxWidth        int    `json:"max_width"`
	MaxHeight       int    `json:"max_height"`
	URL             string `json:"url"`
	PreviewURL      string `json:"preview_url"`
	ImageType       int    `json:"image_type"`
	RenderAsSticker bool   `json:"render_as_sticker"`
}

type ConversationCursors struct {
	Before string `json:"before"`
	After  string `json:"after"`
}
