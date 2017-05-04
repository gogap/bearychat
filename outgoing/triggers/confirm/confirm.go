package confirm

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat"
)

type Confirm struct {
	word   string
	prompt string

	beforeReq        *bearychat.OutgoingRequest
	randomNumber     int32
	randomNumberTime time.Time
	currentHandlers  map[string]bearychat.TriggerHandleFunc
	defaultHandler   func(req *bearychat.OutgoingRequest, msg *bearychat.Message) (err error)

	sync.Mutex
}

func init() {
	bearychat.RegisterTriggerDriver("gogap-confirm", NewConfirm)
}

func NewConfirm(word string, config *configuration.Config) (bearychat.Trigger, error) {

	confirm := &Confirm{
		word:            word,
		currentHandlers: make(map[string]bearychat.TriggerHandleFunc),
	}

	if config != nil {
		confirm.prompt = config.GetString("prompt", "please input number for comfirm")
	} else {
		confirm.prompt = "please input number for comfirm"
	}

	confirm.defaultHandler = confirm.randomHandle

	return confirm, nil
}

func (p *Confirm) Handle(req *bearychat.OutgoingRequest, msg *bearychat.Message) (err error) {
	if handler, exist := p.currentHandlers[req.UserName]; exist {
		return handler(req, msg)
	}

	return p.defaultHandler(req, msg)
}

func (p *Confirm) randomHandle(req *bearychat.OutgoingRequest, resp *bearychat.Message) (err error) {

	rnd := rand.Int31n(99999)
	p.beforeReq = req

	resp.Text = fmt.Sprintf("%s: %d", p.prompt, rnd)

	p.Lock()
	p.currentHandlers[req.UserName] = p.generateComfirmRandomHandle(rnd, *p.beforeReq)
	p.Unlock()

	return bearychat.ErrBreakOnly
}

func (p *Confirm) generateComfirmRandomHandle(number int32, before bearychat.OutgoingRequest) func(*bearychat.OutgoingRequest, *bearychat.Message) error {
	now := time.Now()
	num := number
	originalReq := before

	fn := func(req *bearychat.OutgoingRequest, resp *bearychat.Message) error {
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

		args := req.Args()

		if len(args) != 1 {
			return errors.New("please input numbers")
		}

		strNum := args[0]

		n, err := strconv.Atoi(strNum)
		if err != nil {
			fmt.Println(err, strNum)
			return errors.New("please input numbers")
		}

		if n == 0 {
			return errors.New("please input numbers")
		}

		if num != int32(n) {
			return errors.New("bad comfirm numbers")
		}

		*req = before
		return nil
	}

	return fn
}
