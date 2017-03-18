package incoming

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Request struct {
	Text         string       `json:"text"`
	Notification string       `json:"notification"`
	Markdown     bool         `json:"markdown"`
	Channel      string       `json:"channel"`
	User         string       `json:"user"`
	Attachments  []Attachment `json:"attachments"`
}

type Response struct {
	Code   int         `json:"code"`
	Error  string      `json:"error"`
	Result interface{} `json:"result"`
}

func (p *Response) Err() error {
	if p.Code != 0 {
		return fmt.Errorf("code: %d, %s", p.Code, p.Error)
	}
	return nil
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

type Client struct {
	client *http.Client
}

type ClientOption func(*Client)

func NewClient(opts ...ClientOption) *Client {
	cli := &Client{
		client: &http.Client{},
	}

	cli.Options(opts...)

	return cli
}

func (p *Client) Send(robotId string, token string, req *Request) (resp *Response, err error) {

	if len(robotId) == 0 {
		err = errors.New("robot id is empty")
		return
	}

	if len(token) == 0 {
		err = errors.New("token is empty")
		return
	}

	if req == nil {
		err = errors.New("empty request")
		return
	}

	body, err := json.Marshal(req)
	if err != nil {
		return
	}

	url := fmt.Sprintf("https://hook.bearychat.com/%s/incoming/%s", robotId, token)

	httpResp, err := p.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}

	defer httpResp.Body.Close()

	r := Response{}

	decoder := json.NewDecoder(httpResp.Body)
	decoder.UseNumber()

	err = decoder.Decode(&r)

	if err != nil {
		return
	}

	resp = &r

	return

}

func (p *Client) Options(opts ...ClientOption) {
	for i := 0; i < len(opts); i++ {
		opts[i](p)
	}
}

func TransportOption(transport *http.Transport) ClientOption {
	return func(c *Client) {
		c.client.Transport = transport
	}
}

func TimeoutOption(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}
