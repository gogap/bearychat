package main

import (
	_ "git.rd.rijin.com/automation/bearychat-triggers/deploy/app"
	_ "git.rd.rijin.com/automation/bearychat-triggers/errors_finder"
	_ "git.rd.rijin.com/automation/bearychat-triggers/mns/peeker"
	_ "github.com/gogap/bearychat/outgoing/triggers/auth"
	_ "github.com/gogap/bearychat/outgoing/triggers/commands"
	_ "github.com/gogap/bearychat/outgoing/triggers/confirm"
	_ "github.com/gogap/bearychat/outgoing/triggers/greeter"
	_ "github.com/gogap/bearychat/outgoing/triggers/sensitive_filter"
)
