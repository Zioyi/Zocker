package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call is outside",
	Action: func(context *cli.Context) error {
		log.Infof("")
		return nil
	},
}
