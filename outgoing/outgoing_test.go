package outgoing

import (
	"errors"
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

func (p *Greeter) Handle(req *Request, resp *Response) error {
	switch req.TriggerWord {
	case "!hello":
		{
			resp.Text = "Hello " + req.UserName + " I am " + p.name
			return nil
		}
	case "!morning":
		{
			resp.Text = "Morning " + req.UserName + " I am " + p.name
			return nil
		}
	}

	return errors.New("Unknown TriggerWord")
}

func TestBasicOutgoingRequest(t *testing.T) {

	confStr := `
	{
 		hello = {
 			word = "!hello"
			drivers = [test-greeter]

			test-greeter = {
				name = "robot A"
			}
		}

		morning = {
			word = "!morning"
			drivers = [test-greeter]

			test-greeter = {
				name = "robot B"
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

	resp1 := Response{}
	err = outgoing.Handle(req1, &resp1)

	if err != nil {
		t.Error(err)
		return
	}

	if resp1.Text != "Hello zeal I am robot A" {
		t.Error("bad !hello trigger response: " + resp1.Text)
		return
	}

	resp2 := Response{}
	err = outgoing.Handle(req2, &resp2)
	if err != nil {
		t.Error(err)
		return
	}

	if resp2.Text != "Morning gogap I am robot B" {
		t.Error("bad !morning trigger response: " + resp2.Text)
		return
	}
}
