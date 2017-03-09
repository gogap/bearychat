package outgoing

import (
	"github.com/go-akka/configuration"
)

type OutgoingOption func(*OutgoingSettings)

type OutgoingSettings struct {
	ValidateToken bool
	Tokens        []string
}

func NewOutgoingSettings(config *configuration.Config) *OutgoingSettings {
	if config == nil {
		return &OutgoingSettings{}
	}

	return &OutgoingSettings{
		ValidateToken: config.GetBoolean("validate_token", false),
		Tokens:        config.GetStringList("tokens"),
	}
}
