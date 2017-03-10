package greeter

import (
	"errors"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat/outgoing"
)

type Greeter struct {
	name  string
	word  string
	image string
}

func init() {
	outgoing.RegisterTriggerDriver("gogap-greeter", NewGreeter)
}

func NewGreeter(word string, config *configuration.Config) (outgoing.Trigger, error) {
	return &Greeter{
		word:  word,
		name:  config.GetString("name"),
		image: config.GetString("image"),
	}, nil
}

func (p *Greeter) Handle(req *outgoing.Request, resp *outgoing.Response) error {

	switch req.TriggerWord {
	case "!hello":
		{
			resp.Text = "Hello " + req.UserName + " I am " + p.name
			resp.Attachments = []outgoing.Attachment{
				{
					Images: []outgoing.Image{
						{URL: p.image},
					},
				},
			}
		}
	case "!morning":
		{
			resp.Text = "Morning " + req.UserName + " I am " + p.name
			resp.Attachments = []outgoing.Attachment{
				{
					Images: []outgoing.Image{
						{URL: p.image},
					},
				},
			}
		}
	}

	return errors.New("Known TriggerWord")
}
