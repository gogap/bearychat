package commands

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat/outgoing"
)

type Commands struct {
	namepath map[string]string
	timeout  time.Duration
}

func init() {
	outgoing.RegisterTriggerDriver("gogap-commands", NewCommands)
}

func NewCommands(word string, config *configuration.Config) (outgoing.Trigger, error) {

	commandConfig := config.GetConfig("commands")
	namePath := make(map[string]string)

	if commandConfig != nil {
		commandNames := commandConfig.Root().GetObject().GetKeys()

		for i := 0; i < len(commandNames); i++ {
			namePath[commandNames[i]] = commandConfig.GetString(commandNames[i])
		}
	}

	return &Commands{
		namepath: namePath,
		timeout:  config.GetTimeDuration("timeout", 30),
	}, nil
}

func (p *Commands) Handle(req *outgoing.Request) outgoing.Response {

	args := strings.Split(req.Text, " ")
	if len(args) <= 1 {
		return outgoing.Response{Text: "command argument is too less"}
	}

	var newArgs []string
	for i := 0; i < len(args); i++ {
		s := strings.TrimSpace(args[i])
		if len(s) != 0 {
			newArgs = append(newArgs, s)
		}
	}

	commandName := newArgs[1]

	binPath, exist := p.namepath[commandName]
	if !exist {
		return outgoing.Response{Text: "command not exist"}
	}

	result, err := execCommand(p.timeout, binPath, newArgs[2:]...)

	if err != nil {
		return outgoing.Response{Text: err.Error()}
	}

	return outgoing.Response{Text: result}
}

func execCommand(timeout time.Duration, name string, args ...string) (result string, err error) {

	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)

	err = cmd.Start()

	if err != nil {
		return
	}

	go io.Copy(outBuf, stdout)
	go io.Copy(errBuf, stderr)

	ch := make(chan struct{})

	go func(cmd *exec.Cmd) {
		defer close(ch)
		cmd.Wait()
	}(cmd)

	select {
	case <-ch:
	case <-time.After(timeout):
		cmd.Process.Kill()
		err = errors.New("execute timeout")
		return
	}

	errStr := errBuf.String()
	outString := outBuf.String()

	if len(errStr) > 0 {
		return "", errors.New(errStr)
	}

	return outString, nil
}
