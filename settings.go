package bearychat

import (
	"github.com/go-akka/configuration"
)

type OutgoingOption func(*OutgoingSettings)

type OutgoingSettings struct {
}

func NewOutgoingSettings(config *configuration.Config) *OutgoingSettings {
	return &OutgoingSettings{}
}
