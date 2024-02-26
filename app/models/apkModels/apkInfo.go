package apkModels

type ApkInfo struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Uname   string `json:"uname,omitempty"`
	Icon    string `json:"icon,omitempty"`
	Package string `json:"package,omitempty"`

	File *struct {
		Filesize int    `json:"filesize,omitempty"`
		Path     string `json:"path,omitempty"`
	} `json:"file,omitempty"`

	Obb *struct {
		Main *struct {
			Filesize int    `json:"filesize,omitempty"`
			Path     string `json:"path,omitempty"`
			Filename string `json:"filename,omitempty"`
		} `json:"Main,omitempty"`
	} `json:"obb,omitempty"`

	Media *struct {
		Description string `json:"description,omitempty"`
		// Screenshots []struct {
		// 	Url string `json:"url,omitempty"`
		// } `json:"screenshots,omitempty"`
	} `json:"media,omitempty"`

	// Aab *interface{} `json:"aab,omitempty"`
}
