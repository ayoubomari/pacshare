package models

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
	MID  string `json:"mid"`
	Text string `json:"text"`
}

type PostBack struct {
	Title   string `json:"title"`
	Payload string `json:"payload"`
}
