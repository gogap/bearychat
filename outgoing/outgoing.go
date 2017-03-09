package outgoing

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/go-akka/configuration"
)

const (
	OUTGOING = "gogap-outgoing"
)

var (
	triggerFuncs = make(map[string]NewTriggerFunc)
)

var (
	ErrNilOutgoingName                = errors.New("nil outgoing name")
	ErrNilRequest                     = errors.New("nil request")
	ErrEmptyTriggerWord               = errors.New("empty trigger word")
	ErrTriggerDriverAlreadyRegistered = errors.New("trigger driver already registered")
	ErrNewTriggerFuncIsNil            = errors.New("trigger func is nil")
	ErrNilTrigger                     = errors.New("trigger is nil")
	ErrBadAuthToken                   = errors.New("bad auth token")
)

type ErrorHandlerFunc func(cause error) Response

type NewTriggerFunc func(word string, config *configuration.Config) (Trigger, error)

type Outgoing struct {
	triggers map[string]Trigger // map[word]Trigger
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
		triggers: make(map[string]Trigger),
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

func (p *Outgoing) GetTrigger(triggerWord string) (trigger Trigger, exist bool) {
	trigger, exist = p.triggers[triggerWord]
	return
}

func (p *Outgoing) BindTrigger(triggerName string, triggerWords ...string) *Outgoing {
	if len(triggerWords) == 0 {
		return p
	}

	if len(triggerName) == 0 {
		panic("trigger name could not be empty")
	}

	words := removeDuplicates(triggerWords)

	for i := 0; i < len(words); i++ {
		word := strings.TrimSpace(words[i])
		if len(word) == 0 {
			panic("word could not be empty")
		}

		triggerDriver, exist := triggerFuncs[triggerName]
		if !exist {
			panic(fmt.Errorf("the trigger of %s did not exist", triggerName))
		}

		trigger, err := triggerDriver(word, p.config.GetConfig("words").GetConfig(word).GetConfig("options"))
		if err != nil {
			panic(err)
		}

		p.triggers[word] = trigger
	}

	return p
}

func (p *Outgoing) Words() []string {
	var ret []string
	for k, _ := range p.triggers {
		ret = append(ret, k)
	}

	sort.Sort(sort.StringSlice(ret))

	return ret
}

func (p *Outgoing) BindTriggerDirect(trigger Trigger, triggerWords ...string) *Outgoing {

	if trigger == nil {
		panic(ErrNilTrigger)
	}

	words := removeDuplicates(triggerWords)

	for i := 0; i < len(words); i++ {
		word := strings.TrimSpace(words[i])
		if len(word) == 0 {
			panic("word could not be empty")
		}

		p.triggers[word] = trigger
	}

	return p
}

func (p *Outgoing) UnbindTrigger(triggerWords ...string) {
	for i := 0; i < len(triggerWords); i++ {
		delete(p.triggers, triggerWords[i])
	}
}

func (p *Outgoing) SetErrorHandler(handler ErrorHandlerFunc) {
	p.errorHandler = handler
	if p.errorHandler == nil {
		p.errorHandler = p.handleError
	}
}

func (p *Outgoing) Handle(req *Request) Response {
	if err := p.validateRequest(req); err != nil {
		return p.errorHandler(err)
	}

	word := strings.TrimLeft(req.TriggerWord, "!")
	trigger, exist := p.triggers[word]
	if !exist {
		err := fmt.Errorf("trigger of %s not exist!", word)
		return p.errorHandler(err)
	}

	return trigger.Handle(req)
}

func (p *Outgoing) HandleHttpRequest(rw http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	decoder.UseNumber()

	triggerReq := &Request{}
	err := decoder.Decode(triggerReq)

	var resp Response
	if err != nil {
		resp = p.errorHandler(err)
	} else {
		resp = p.Handle(triggerReq)
	}

	jsonResp, _ := json.Marshal(resp)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonResp)
}

func (p *Outgoing) validateRequest(req *Request) error {
	if req == nil {
		return ErrNilRequest
	}

	if len(strings.TrimSpace(req.TriggerWord)) == 0 {
		return ErrEmptyTriggerWord
	}

	if p.settings.ValidateToken {
		if len(p.settings.Tokens) > 0 {
			for i := 0; i < len(p.settings.Tokens); i++ {
				if req.Token == p.settings.Tokens[i] {
					return nil
				}
			}

			return ErrBadAuthToken
		}
	}

	return nil
}

func (p *Outgoing) handleError(cause error) Response {
	return Response{
		Text: cause.Error(),
	}
}

func (p *Outgoing) autoBind(config *configuration.Config) {
	if config == nil {
		return
	}

	wordsConfig := config.GetConfig("words")

	words := wordsConfig.Root().GetObject().GetKeys()

	for i := 0; i < len(words); i++ {
		driver := wordsConfig.GetConfig(words[i]).GetString("driver")
		p.BindTrigger(driver, words[i])
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
