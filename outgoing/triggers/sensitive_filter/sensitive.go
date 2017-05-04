package sensitive_filter

import (
	"regexp"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat"
)

type Sensitive struct {
	expressions []*regexp.Regexp
}

func init() {
	bearychat.RegisterTriggerDriver("gogap-sensitive-filter", NewSensitive)
}

func NewSensitive(word string, config *configuration.Config) (bearychat.Trigger, error) {

	if config == nil {
		return &Sensitive{}, nil
	}

	expressions := config.GetStringList("expressions")

	var regExprs []*regexp.Regexp

	for i := 0; i < len(expressions); i++ {
		expr, err := regexp.Compile(expressions[i])
		if err != nil {
			return nil, err
		}

		regExprs = append(regExprs, expr)
	}

	return &Sensitive{expressions: regExprs}, nil
}

func (p *Sensitive) Handle(req *bearychat.OutgoingRequest, msg *bearychat.Message) (err error) {

	if len(p.expressions) == 0 {
		return
	}

	txt := msg.Text
	for i := 0; i < len(p.expressions); i++ {
		txt = p.expressions[i].ReplaceAllString(txt, "******")
	}

	msg.Text = txt

	return
}
