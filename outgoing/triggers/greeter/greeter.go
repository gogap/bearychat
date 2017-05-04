package greeter

import (
	"errors"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat"
)

type Greeter struct {
	name  string
	word  string
	image string
}

func init() {
	bearychat.RegisterTriggerDriver("gogap-greeter", NewGreeter)
}

func NewGreeter(word string, config *configuration.Config) (bearychat.Trigger, error) {
	return &Greeter{
		word:  word,
		name:  config.GetString("name"),
		image: config.GetString("image"),
	}, nil
}

func (p *Greeter) Handle(req *bearychat.OutgoingRequest, msg *bearychat.Message) error {

	switch req.TriggerWord {
	case "!hello":
		{
			msg.Text = "Hello " + req.UserName + " I am " + p.name
			msg.Attachments = []bearychat.Attachment{
				{
					Images: []bearychat.Image{
						{URL: p.image},
					},
				},
			}
		}
	case "!morning":
		{
			msg.Text = "Morning " + req.UserName + " I am " + p.name
			msg.Attachments = []bearychat.Attachment{
				{
					Images: []bearychat.Image{
						{URL: p.image},
					},
				},
			}
		}
	}

	return errors.New("Known TriggerWord")
}
