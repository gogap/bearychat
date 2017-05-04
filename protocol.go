package bearychat

import (
	"fmt"
	"strings"
)

type Image struct {
	URL string `json:"url"`
}

type Attachment struct {
	Title  string  `json:"title"`
	Text   string  `json:"text"`
	Color  string  `json:"color"`
	Images []Image `json:"images"`
}

type Message struct {
	Text         string       `json:"text"`
	Notification string       `json:"notification"`
	Markdown     bool         `json:"markdown"`
	Channel      string       `json:"channel"`
	User         string       `json:"user"`
	Attachments  []Attachment `json:"attachments"`
}

type OutgoingRequest struct {
	Token       string   `json:"token"`
	Timestamp   int      `json:"ts"`
	Text        string   `json:"text"`
	TriggerWord string   `json:"trigger_word"`
	Subdomain   string   `json:"subdomain"`
	ChannelName string   `json:"channel_name"`
	UserName    string   `json:"user_name"`
	Commands    []string `json:"-"`
}

func (p *OutgoingRequest) Args() []string {

	strArgs := strings.TrimSpace(strings.TrimPrefix(p.Text, p.TriggerWord))

	for i := 0; i < len(p.Commands); i++ {
		strArgs = strings.TrimSpace(strings.TrimPrefix(strArgs, p.Commands[i]))
	}

	if len(strArgs) == 0 {
		return nil
	}

	return strings.Fields(strArgs)
}

type IncomingResponse struct {
	Code   int         `json:"code"`
	Error  string      `json:"error"`
	Result interface{} `json:"result"`
}

func (p *IncomingResponse) Err() error {
	if p.Code != 0 {
		return fmt.Errorf("code: %d, %s", p.Code, p.Error)
	}
	return nil
}
