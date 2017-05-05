package user_filter

import (
	"errors"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat"
)

type UserFilter struct {
	users map[string]bool
}

func init() {
	bearychat.RegisterTriggerDriver("gogap-user-filter", NewUserFilter)
}

func NewUserFilter(word string, config *configuration.Config) (bearychat.Trigger, error) {
	if config == nil {
		return &UserFilter{
			users: make(map[string]bool),
		}, nil
	}

	users := config.GetStringList("users")

	usersMap := make(map[string]bool, len(users))

	for _, user := range users {
		usersMap[user] = true
	}

	uf := &UserFilter{
		users: usersMap,
	}

	return uf, nil
}

func (p *UserFilter) Handle(req *bearychat.OutgoingRequest, msg *bearychat.Message) (err error) {

	if !p.users[req.UserName] {
		err = errors.New("gogap-user-filter: permision denied.")
		return
	}

	return
}
