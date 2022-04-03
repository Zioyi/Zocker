package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `zocker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
               Enjoy it, just fo fun.`

func main() {
	app := cli.NewApp()
	app.Name = "zocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		commitCommand,
		listCommand,
		logCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetReportCaller(true)
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:12:04",
		})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
