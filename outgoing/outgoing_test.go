package outgoing

import (
	"testing"

	"github.com/go-akka/configuration"
)

type Greeter struct {
	name string
	word string
}

func init() {
	RegisterTriggerDriver("test-greeter", NewGreeter)
}

func NewGreeter(word string, config *configuration.Config) (Trigger, error) {
	return &Greeter{
		word: word,
		name: config.GetString("name"),
	}, nil
}

func (p *Greeter) Handle(req *Request) Response {

	switch req.TriggerWord {
	case "!hello":
		{
			return Response{
				Text: "Hello " + req.UserName + " I am " + p.name,
			}
		}
	case "!morning":
		{
			return Response{
				Text: "Morning " + req.UserName + " I am " + p.name,
			}
		}
	}

	return Response{
		Text: "Unknown TriggerWord",
	}

}

func TestOutgoingAuthTokenValiedate(t *testing.T) {

	confStr := `
	{
		validate_token = true
		tokens = ["abc"]

		words {
			hello = {
				driver = test-greeter
				options.name = "robot A"
			}
		}
	}

	`

	config := configuration.ParseString(confStr)

	outgoing, err := NewOutgoing(config)

	if outgoing == nil {
		t.Error(err)
		return
	}

	req1 := &Request{
		Text:        "!hello my name is zeal",
		UserName:    "zeal",
		TriggerWord: "!hello",
	}

	resp1 := outgoing.Handle(req1)
	if resp1.Text != "bad auth token" {
		t.Error("bad !hello trigger response: " + resp1.Text)
		return
	}

	req2 := &Request{
		Text:        "!hello my name is zeal",
		UserName:    "zeal",
		TriggerWord: "!hello",
		Token:       "abc",
	}

	resp2 := outgoing.Handle(req2)
	if resp2.Text != "Hello zeal I am robot A" {
		t.Error("bad !hello trigger response: " + resp2.Text)
		return
	}
}

func TestOutgoingAuthTokenValiedate2(t *testing.T) {

	confStr := `
	{
		words {
			hello = {
				driver = gogap-outgoing
				options = {
					validate_token = true
					tokens = ["bcd"]

					words {
						hello = {	
							driver = test-greeter
							options.name = "robot A"
						}
					}
				}
			}
		}
	}

	`

	config := configuration.ParseString(confStr)

	outgoing, err := NewOutgoing(config)

	if err != nil {
		t.Error(err)
		return
	}

	req := &Request{
		Text:        "!hello my name is zeal",
		UserName:    "zeal",
		TriggerWord: "!hello",
		Token:       "abc",
	}

	resp := outgoing.Handle(req)
	if resp.Text != "bad auth token" {
		t.Error("bad !hello trigger response: " + resp.Text)
		return
	}

	req2 := &Request{
		Text:        "!hello my name is zeal",
		UserName:    "zeal",
		TriggerWord: "!hello",
		Token:       "bcd",
	}

	resp2 := outgoing.Handle(req2)
	if resp2.Text != "Hello zeal I am robot A" {
		t.Error("bad !hello trigger response: " + resp2.Text)
		return
	}
}

func TestBasicOutgoingRequest(t *testing.T) {

	confStr := `
	{
		options = {

		}

		words {

			hello = {
				driver = test-greeter
				options.name = "robot A"
			}

			morning = {
				driver = test-greeter
				options.name = "robot B"
			}
		}
	}

	`

	config := configuration.ParseString(confStr)

	outgoing, err := NewOutgoing(config)

	if outgoing == nil {
		t.Error(err)
		return
	}

	req1 := &Request{
		Text:        "!hello my name is zeal",
		UserName:    "zeal",
		TriggerWord: "!hello",
	}

	req2 := &Request{
		Text:        "!hello my name is gogap",
		UserName:    "gogap",
		TriggerWord: "!morning",
	}

	resp1 := outgoing.Handle(req1)
	if resp1.Text != "Hello zeal I am robot A" {
		t.Error("bad !hello trigger response: " + resp1.Text)
		return
	}

	resp2 := outgoing.Handle(req2)
	if resp2.Text != "Morning gogap I am robot B" {
		t.Error("bad !morning trigger response: " + resp2.Text)
		return
	}
}
