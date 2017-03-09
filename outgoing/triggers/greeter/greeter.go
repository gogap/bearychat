package greeter

import (
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

func (p *Greeter) Handle(req *outgoing.Request) outgoing.Response {

	switch req.TriggerWord {
	case "!hello":
		{
			return outgoing.Response{
				Text: "Hello " + req.UserName + " I am " + p.name,
				Attachments: []outgoing.Attachment{
					{
						Images: []outgoing.Image{
							{URL: p.image},
						},
					},
				},
			}
		}
	case "!morning":
		{
			return outgoing.Response{
				Text: "Morning " + req.UserName + " I am " + p.name,
				Attachments: []outgoing.Attachment{
					{
						Images: []outgoing.Image{
							{URL: p.image},
						},
					},
				},
			}
		}
	}

	return outgoing.Response{
		Text: "Unknown TriggerWord",
	}

}
