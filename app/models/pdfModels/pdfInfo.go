package pdfModels

type ApkInfo struct {
	Cover     string `json:"cover,omitempty"`
	Name      string `json:"name,omitempty"`
	Author    string `json:"author,omitempty"`
	FileSize  string `json:"fileSize,omitempty"`
	PagesNum  string `json:"pagesNum,omitempty"`
	Language  string `json:"language,omitempty"`
	SessionID string `json:"sessionID,omitempty"`
}
