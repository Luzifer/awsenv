package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Luzifer/awsenv/shellsupport"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	_ "github.com/Luzifer/awsenv/shellsupport/bash"
	_ "github.com/Luzifer/awsenv/shellsupport/fish"
)

func getCmdShell() cli.Command {
	return cli.Command{
		Name:  "shell",
		Usage: "print the AWS credentials in a format for your shell to eval()",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "shell,s",
				Value:  "",
				Usage:  "name of the shell to export for",
				EnvVar: "SHELL",
			},
			cli.BoolTFlag{
				Name:  "export,x",
				Usage: "Adds proper export options for your shell",
			},
		},
		Action: actionCmdShell,
	}
}

func actionCmdShell(c *cli.Context) {
	if len(c.String("shell")) == 0 {
		log.Errorf("Could not determine your shell. Please provide --shell")
		os.Exit(1)
	}
	s := strings.Split(c.String("shell"), "/")
	shell := s[len(s)-1]

	log.Debugf("Found shell '%s'", shell)

	handler, err := shellsupport.GetShellHandler(shell)
	if err != nil {
		log.Errorf("Could not find a handler for '%s' shell", shell)
		os.Exit(1)
	}

	if !c.Args().Present() {
		log.Errorf("Please specify the enviroment to load")
		os.Exit(1)
	}

	if a, ok := awsCredentials.Credentials[c.Args().First()]; ok {
		fmt.Println(strings.Join(handler(a, c.Bool("export")), "\n"))
		os.Exit(0)
	}

	log.Errorf("Could not find environment '%s'", c.Args().First())
	os.Exit(1)
}
