package confirm

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat/outgoing"
)

type Confirm struct {
	word   string
	prompt string

	beforeReq        *outgoing.Request
	randomNumber     int32
	randomNumberTime time.Time
	currentHandlers  map[string]outgoing.TriggerHandleFunc
	defaultHandler   func(req *outgoing.Request, resp *outgoing.Response) (err error)

	sync.Mutex
}

func init() {
	outgoing.RegisterTriggerDriver("gogap-confirm", NewConfirm)
}

func NewConfirm(word string, config *configuration.Config) (outgoing.Trigger, error) {

	confirm := &Confirm{
		word:            word,
		currentHandlers: make(map[string]outgoing.TriggerHandleFunc),
	}

	if config != nil {
		confirm.prompt = config.GetString("prompt")
	} else {
		confirm.prompt = "please input number for comfirm"
	}

	confirm.defaultHandler = confirm.randomHandle

	return confirm, nil
}

func (p *Confirm) Handle(req *outgoing.Request, resp *outgoing.Response) (err error) {
	if handler, exist := p.currentHandlers[req.UserName]; exist {
		return handler(req, resp)
	}

	return p.defaultHandler(req, resp)
}

func (p *Confirm) randomHandle(req *outgoing.Request, resp *outgoing.Response) (err error) {

	rnd := rand.Int31n(99999)
	p.beforeReq = req

	resp.Text = fmt.Sprintf("%s: %d", p.prompt, rnd)

	p.Lock()
	p.currentHandlers[req.UserName] = p.generateComfirmRandomHandle(rnd, *p.beforeReq)
	p.Unlock()

	return outgoing.ErrBreakOnly
}

func (p *Confirm) generateComfirmRandomHandle(number int32, before outgoing.Request) func(*outgoing.Request, *outgoing.Response) error {
	now := time.Now()
	num := number
	originalReq := before

	fn := func(req *outgoing.Request, resp *outgoing.Response) error {
		defer func() {
			p.Lock()
			delete(p.currentHandlers, req.UserName)
			p.Unlock()
		}()

		if req.UserName != originalReq.UserName {
			return p.defaultHandler(req, resp)
		}

		if time.Now().Sub(now).Seconds() > 30 {
			return p.defaultHandler(req, resp)
		}

		strNum := strings.TrimSpace(strings.TrimPrefix(req.Text, p.word))

		n, err := strconv.Atoi(strNum)
		if err != nil {
			return errors.New("please input number")
		}

		if n == 0 {
			return errors.New("please input number")
		}

		if num != int32(n) {
			return errors.New("bad comfirm number")
		}

		*req = before
		return nil
	}

	return fn
}
