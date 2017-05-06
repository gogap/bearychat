package bearychat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat/internal"
)

const (
	OUTGOING = "gogap-outgoing"
)

var (
	triggerFuncs = make(map[string]NewTriggerFunc)
)

var (
	ErrTriggerDriverAlreadyRegistered = errors.New("trigger driver already registered")
	ErrNewTriggerFuncIsNil            = errors.New("trigger func is nil")
	ErrBreakOnly                      = errors.New("break only")
	ErrNoContent                      = errors.New("no content")
)

type ErrorHandlerFunc func(cause error) Message

type NewTriggerFunc func(word string, config *configuration.Config) (Trigger, error)

type Outgoing struct {
	triggers map[string]*internal.Command // map[word]Command tree

	settings *OutgoingSettings

	config       *configuration.Config
	errorHandler ErrorHandlerFunc
}

func init() {
	RegisterTriggerDriver(OUTGOING, NewOutgoingTrigger)
}

func TriggerDrivers() []string {
	var ret []string
	for k, _ := range triggerFuncs {
		ret = append(ret, k)
	}

	sort.Sort(sort.StringSlice(ret))

	return ret
}

func RegisterTriggerDriver(name string, fn NewTriggerFunc) {
	if fn == nil {
		panic(ErrNewTriggerFuncIsNil)
	}

	_, exist := triggerFuncs[name]
	if exist {
		panic(ErrTriggerDriverAlreadyRegistered)
	}
	triggerFuncs[name] = fn
}

func NewOutgoing(config *configuration.Config) (*Outgoing, error) {

	outgoing := &Outgoing{
		triggers: make(map[string]*internal.Command),
		config:   config,
		settings: NewOutgoingSettings(config),
	}

	outgoing.errorHandler = outgoing.handleError

	outgoing.autoBind(config)

	return outgoing, nil
}

func NewOutgoingTrigger(word string, config *configuration.Config) (Trigger, error) {
	return NewOutgoing(config)
}

func (p *Outgoing) BindTrigger(config *configuration.Config) *Outgoing {
	triggerWord := config.GetString("word")
	triggerWord = strings.TrimSpace(triggerWord)

	commands := config.GetStringList("commands")

	drivers := config.GetStringList("drivers")

	if len(triggerWord) == 0 {
		return p
	}

	if len(drivers) == 0 {
		return p
	}

	names := removeDuplicates(drivers)

	var triggers []interface{}

	for i := 0; i < len(names); i++ {
		triggerDriver, exist := triggerFuncs[names[i]]
		if !exist {
			panic(fmt.Errorf("the trigger of %s did not exist", names[i]))
		}

		trigger, err := triggerDriver(triggerWord, config.GetConfig(names[i]))
		if err != nil {
			panic(err)
		}

		triggers = append(triggers, trigger)
	}

	root, exist := p.triggers[triggerWord]
	if !exist {
		root = &internal.Command{}
	}

	node := root.Match(commands...)

	if len(node.Values) > 0 {
		panic(fmt.Errorf("command alrady has triggers: %s", strings.Join(commands, " ")))
	}

	subCommands := commands[len(node.Commands()):]

	for i := 0; i < len(subCommands); i++ {

		child := &internal.Command{
			Name: subCommands[i],
		}

		if i+1 == len(subCommands) {
			child.Values = triggers
		}

		node.AddChild(child)
		node = child
	}

	p.triggers[triggerWord] = root

	return p
}

func (p *Outgoing) SetErrorHandler(handler ErrorHandlerFunc) {
	p.errorHandler = handler
	if p.errorHandler == nil {
		p.errorHandler = p.handleError
	}
}

func (p *Outgoing) Handle(req *OutgoingRequest, msg *Message) error {

	word := strings.TrimSpace(req.TriggerWord)
	treeRoot, exist := p.triggers[word]

	if !exist {
		return fmt.Errorf("trigger of %s not exist!", word)
	}

	args := req.Args()

	node := treeRoot.Match(args...)

	if node == treeRoot {
		return fmt.Errorf("unknown sub-command: %s", strings.Join(args, " "))
	}

	if len(node.Values) == 0 {
		return fmt.Errorf("unfinished sub-command")
	}

	req.Commands = node.Commands()

	for i := 0; i < len(node.Values); i++ {
		if err := node.Values[i].(Trigger).Handle(req, msg); err != nil {
			return err
		}
	}

	return nil
}

func (p *Outgoing) HandleHttpRequest(rw http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	decoder.UseNumber()

	triggerReq := &OutgoingRequest{}
	err := decoder.Decode(triggerReq)

	statusCode := 200

	msg := Message{}
	if err != nil {
		msg = p.errorHandler(err)
	} else {
		err = p.Handle(triggerReq, &msg)
		if err == ErrBreakOnly {
			err = nil
		} else if err == ErrNoContent {
			err = nil
			statusCode = 204
		}

		if err != nil {
			msg = p.errorHandler(err)
		}
	}

	jsonMsg, _ := json.Marshal(msg)

	rw.Header().Set("Content-Type", "application/json")

	if statusCode != 200 {
		rw.WriteHeader(statusCode)
	} else {
		rw.Write(jsonMsg)
	}
}

func (p *Outgoing) handleError(cause error) Message {
	return Message{
		Text: cause.Error(),
	}
}

func (p *Outgoing) autoBind(config *configuration.Config) {
	if config == nil {
		return
	}

	keys := p.config.Root().GetObject().GetKeys()

	for i := 0; i < len(keys); i++ {
		p.BindTrigger(p.config.GetConfig(keys[i]))
	}
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range elements {
		if _, exist := encountered[v]; !exist {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}
