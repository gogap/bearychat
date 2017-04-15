package outgoing

import (
	"strings"
)

type Request struct {
	Token       string   `json:"token"`
	Timestamp   int      `json:"ts"`
	Text        string   `json:"text"`
	TriggerWord string   `json:"trigger_word"`
	Subdomain   string   `json:"subdomain"`
	ChannelName string   `json:"channel_name"`
	UserName    string   `json:"user_name"`
	Commands    []string `json:"-"`
}

func (p *Request) Args() []string {

	strArgs := strings.TrimSpace(strings.TrimPrefix(p.Text, p.TriggerWord))

	for i := 0; i < len(p.Commands); i++ {
		strArgs = strings.TrimSpace(strings.TrimPrefix(strArgs, p.Commands[i]))
	}

	if len(strArgs) == 0 {
		return nil
	}

	return strings.Fields(strArgs)
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
	Markdown    bool         `json:"markdown"`
	Attachments []Attachment `json:"attachments"`
}
