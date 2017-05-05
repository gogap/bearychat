package channel_filter

import (
	"errors"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat"
)

type ChannelFilter struct {
	channels map[string]bool
}

func init() {
	bearychat.RegisterTriggerDriver("gogap-channel-filter", NewChannelFilter)
}

func NewChannelFilter(word string, config *configuration.Config) (bearychat.Trigger, error) {
	if config == nil {
		return &ChannelFilter{
			channels: make(map[string]bool),
		}, nil
	}

	channels := config.GetStringList("channels")

	channelsMap := make(map[string]bool, len(channels))

	for _, channel := range channels {
		channelsMap[channel] = true
	}

	uf := &ChannelFilter{
		channels: channelsMap,
	}

	return uf, nil
}

func (p *ChannelFilter) Handle(req *bearychat.OutgoingRequest, msg *bearychat.Message) (err error) {

	if !p.channels[req.ChannelName] {
		err = errors.New("gogap-channel-filter: Illegal channel.")
		return
	}

	return
}
