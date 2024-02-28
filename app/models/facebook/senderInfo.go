package facebook

type SenderInfo struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ProfilePic string `json:"profile_pic"`

	// -------- those need autorization request from facebook ------
	// Locale     string `json:"locale"`
	// Timezone   int    `json:"timezone"`
	// Gender     string `json:"gender"`
}
