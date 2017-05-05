package confirm

import (
	"errors"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"sync"
	"time"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat"
)

var (
	UserTOTPSecretNotExist = errors.New("user totp secret not exist")
)

type TOTPConfirm struct {
	prompt string

	userSecret map[string]string

	beforeReq       *bearychat.OutgoingRequest
	currentHandlers map[string]bearychat.TriggerHandleFunc
	defaultHandler  bearychat.TriggerHandleFunc

	period int32

	sync.Mutex
}

func init() {
	bearychat.RegisterTriggerDriver("gogap-confirm-totp", NewTOTPConfirm)
}

func NewTOTPConfirm(word string, config *configuration.Config) (bearychat.Trigger, error) {

	confirm := &TOTPConfirm{
		userSecret:      make(map[string]string),
		currentHandlers: make(map[string]bearychat.TriggerHandleFunc),
	}

	if config != nil {
		confirm.prompt = config.GetString("prompt", "please input one time password for comfirm")

		secConfig := config.GetConfig("secrets")

		confirm.period = config.GetInt32("period", 30)

		if secConfig != nil {
			keys := secConfig.Root().GetObject().GetKeys()
			for _, key := range keys {
				user := secConfig.GetString(key + ".user")
				secret := secConfig.GetString(key + ".secret")
				confirm.userSecret[user] = secret
			}
		}
	} else {
		confirm.prompt = "please input one time password for comfirm"
	}

	confirm.defaultHandler = confirm.totpHandle

	return confirm, nil
}

func (p *TOTPConfirm) Handle(req *bearychat.OutgoingRequest, msg *bearychat.Message) (err error) {
	if handler, exist := p.currentHandlers[req.UserName]; exist {
		fmt.Println("exist", handler)
		return handler(req, msg)
	}

	return p.defaultHandler(req, msg)
}

func (p *TOTPConfirm) totpHandle(req *bearychat.OutgoingRequest, msg *bearychat.Message) (err error) {

	if secret := p.userSecret[req.UserName]; len(secret) == 0 {
		err = UserTOTPSecretNotExist
		return
	}

	p.beforeReq = req

	msg.Text = p.prompt

	p.Lock()
	p.currentHandlers[req.UserName] = p.generateComfirmHandle(*p.beforeReq)
	p.Unlock()

	return bearychat.ErrBreakOnly
}

func (p *TOTPConfirm) generateComfirmHandle(before bearychat.OutgoingRequest) bearychat.TriggerHandleFunc {
	now := time.Now()
	originalReq := before

	fn := func(req *bearychat.OutgoingRequest, msg *bearychat.Message) error {
		defer func() {
			p.Lock()
			delete(p.currentHandlers, req.UserName)
			p.Unlock()
		}()

		if req.UserName != originalReq.UserName {
			return p.defaultHandler(req, msg)
		}

		if time.Now().Sub(now).Seconds() > p.period {
			return p.defaultHandler(req, msg)
		}

		args := req.Args()

		if len(args) != 1 {
			return errors.New(p.period)
		}

		passcode := args[0]
		secret := p.userSecret[before.UserName]

		if rv, _ := totp.ValidateCustom(
			passcode,
			secret,
			time.Now().UTC(),
			totp.ValidateOpts{
				Period:    uint(p.period),
				Skew:      1,
				Digits:    otp.DigitsSix,
				Algorithm: otp.AlgorithmSHA1,
			},
		); !rv {
			return errors.New("bad one time password ")
		}

		*req = before
		return nil
	}

	return fn
}
