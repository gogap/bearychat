package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-akka/configuration"
	"github.com/gogap/bearychat/outgoing"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

func main() {

	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "outgoing"
	app.Usage = "a service for bearychat"
	app.HelpName = "outgoing"

	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "run bearychat outgoing service",
			Action: cmdRun,
			Flags:  []cli.Flag{ConfigFlag},
		},
	}

	app.Run(os.Args)
}

func cmdRun(c *cli.Context) (err error) {
	filename := c.String(ConfigFlag.Name)

	if len(filename) == 0 {
		filename = "outgoing.conf"
	}

	config := configuration.LoadConfig(filename)

	httpConfig := config.GetConfig("http")

	if httpConfig == nil {
		err = fmt.Errorf("config of http section did not set")
		return
	}

	outgoingConfig := config.GetConfig("outgoing")

	if outgoingConfig == nil {
		err = fmt.Errorf("config of outgoing section did not set")
		return
	}

	var out *outgoing.Outgoing
	out, err = initOutgoing(outgoingConfig)
	if err != nil {
		return
	}

	mux := http.NewServeMux()

	words := out.Words()

	for _, word := range words {
		mux.HandleFunc("/"+word, out.HandleHttpRequest)
	}

	n := negroni.Classic()
	n.UseHandler(mux)

	err = http.ListenAndServe(httpConfig.GetString("address", ":8080"), n)

	return err
}

func initOutgoing(config *configuration.Config) (*outgoing.Outgoing, error) {
	outgoing, err := outgoing.NewOutgoing(config)
	if err != nil {
		return nil, err
	}

	return outgoing, nil

}
