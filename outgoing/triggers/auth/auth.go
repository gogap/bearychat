package auth

import (
	"errors"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat/outgoing"
)

type Auth struct {
	word  string
	token string
}

func init() {
	outgoing.RegisterTriggerDriver("gogap-auth", NewAuth)
}

func NewAuth(word string, config *configuration.Config) (outgoing.Trigger, error) {
	return &Auth{
		word:  word,
		token: config.GetString("token"),
	}, nil
}

func (p *Auth) Handle(req *outgoing.Request, resp *outgoing.Response) (err error) {

	if req.TriggerWord != p.word {
		err = errors.New("bad request trigger word in gogap-auth")
		return
	}

	if req.Token != p.token {
		err = errors.New("error auth token")
		return
	}

	return
}
