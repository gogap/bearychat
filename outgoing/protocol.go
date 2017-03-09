package outgoing

type Request struct {
	Token       string `json:"token"`
	Timestamp   int    `json:"ts"`
	Text        string `json:"text"`
	TriggerWord string `json:"trigger_word"`
	Subdomain   string `json:"subdomain"`
	ChannelName string `json:"channel_name"`
	UserName    string `json:"user_name"`
}

type Image struct {
	URL string `json:"url"`
}

type Attachment struct {
	Title  string  `json:"title"`
	Text   string  `json:"text"`
	Color  string  `json:"color"`
	Images []Image `json:"images"`
}

type Response struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}
