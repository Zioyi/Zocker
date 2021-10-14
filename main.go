package main

import "github.com/urfave/cli"

const usage = `zocker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
               Enjoy it, just fo fun.`

func main() {
	app := cli.NewApp()
	app.Name = "zocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
	}
}
