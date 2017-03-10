package commands

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat/outgoing"
)

type _CMD struct {
	cmd string
	cwd string
}

type Commands struct {
	namepath map[string]_CMD
	timeout  time.Duration

	defaultCWD string
}

func init() {
	outgoing.RegisterTriggerDriver("gogap-commands", NewCommands)
}

func NewCommands(word string, config *configuration.Config) (outgoing.Trigger, error) {

	commandConfig := config.GetConfig("commands")
	namePath := make(map[string]_CMD)

	if commandConfig != nil {
		commandNames := commandConfig.Root().GetObject().GetKeys()

		for i := 0; i < len(commandNames); i++ {
			conf := commandConfig.GetConfig(commandNames[i])
			if conf == nil {
				return nil, errors.New("command of " + commandNames[i] + "'s config not exist")
			}

			cmd := _CMD{
				cmd: conf.GetString("cmd"),
				cwd: conf.GetString("cwd"),
			}

			namePath[commandNames[i]] = cmd
		}
	}

	cwd, _ := os.Getwd()
	defaultCWD := config.GetString("cwd", cwd)

	return &Commands{
		namepath:   namePath,
		timeout:    config.GetTimeDuration("timeout", 30),
		defaultCWD: defaultCWD,
	}, nil
}

func (p *Commands) Handle(req *outgoing.Request, resp *outgoing.Response) error {

	args := strings.Split(req.Text, " ")
	if len(args) <= 1 {
		return errors.New("command argument is too less")
	}

	var newArgs []string
	for i := 0; i < len(args); i++ {
		s := strings.TrimSpace(args[i])
		if len(s) != 0 {
			newArgs = append(newArgs, s)
		}
	}

	commandName := newArgs[1]

	cmd, exist := p.namepath[commandName]
	if !exist {
		return errors.New("command not exist")
	}

	cwd := cmd.cwd
	if len(cwd) == 0 {
		cwd = p.defaultCWD
	}

	result, err := execCommand(p.timeout, cwd, cmd.cmd, newArgs[2:]...)

	if err != nil {
		return err
	}

	resp.Text = result

	return nil
}

func execCommand(timeout time.Duration, cwd string, name string, args ...string) (result string, err error) {

	cmd := exec.Command(name, args...)
	cmd.Dir = cwd

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
