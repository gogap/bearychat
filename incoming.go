package bearychat

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type IncomingClient struct {
	client *http.Client
}

type ClientOption func(*IncomingClient)

func NewIncomingClient(opts ...ClientOption) *IncomingClient {
	cli := &IncomingClient{
		client: &http.Client{},
	}

	cli.Options(opts...)

	return cli
}

func (p *IncomingClient) Send(url string, msg *Message) (resp *IncomingResponse, err error) {

	if len(url) == 0 {
		err = errors.New("url is empty")
		return
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return
	}

	httpResp, err := p.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return
	}

	defer httpResp.Body.Close()

	r := IncomingResponse{}

	decoder := json.NewDecoder(httpResp.Body)
	decoder.UseNumber()

	err = decoder.Decode(&r)

	if err != nil {
		return
	}

	resp = &r

	return
}

func (p *IncomingClient) Options(opts ...ClientOption) {
	for i := 0; i < len(opts); i++ {
		opts[i](p)
	}
}

func TransportOption(transport *http.Transport) ClientOption {
	return func(c *IncomingClient) {
		c.client.Transport = transport
	}
}

func TimeoutOption(timeout time.Duration) ClientOption {
	return func(c *IncomingClient) {
		c.client.Timeout = timeout
	}
}
