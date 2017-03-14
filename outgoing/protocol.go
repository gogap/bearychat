package outgoing

import (
	"strings"
)

type Request struct {
	Token       string `json:"token"`
	Timestamp   int    `json:"ts"`
	Text        string `json:"text"`
	TriggerWord string `json:"trigger_word"`
	Subdomain   string `json:"subdomain"`
	ChannelName string `json:"channel_name"`
	UserName    string `json:"user_name"`
}

func (p *Request) Args() []string {

	strArgs := strings.TrimSpace(strings.TrimLeft(p.Text, p.TriggerWord))

	if len(strArgs) == 0 {
		return nil
	}

	var retArgs []string
	args := strings.Split(strArgs, " ")

	for i := 0; i < len(args); i++ {
		s := strings.TrimSpace(args[i])
		if len(s) != 0 {
			retArgs = append(retArgs, s)
		}
	}
	return retArgs
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
