package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func getCmdConsole() cli.Command {
	return cli.Command{
		Name:  "console",
		Usage: "prints a sign-in URL for the AWS web console",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "duration,d",
				Value: 8 * 60,
				Usage: "time in minutes the sign-in is valid",
			},
		},
		Action: actionCmdConsole,
	}
}

func actionCmdConsole(c *cli.Context) {
	if !c.Args().Present() {
		cli.ShowCommandHelp(c, "console")
		log.Error("Please specify the name of the environment to create the sign-in URL for")
		os.Exit(1)
	}

	if c.Int("duration")*60 < 900 || c.Int("duration")*60 > 129600 {
		log.Errorln("Duration parameter must between 15 and 2160 minutes.")
		os.Exit(1)
	}

	subconsole := "console"
	if len(c.Args()) == 2 {
		subconsole = c.Args()[1]
	}

	loginURL, err := awsCredentials.GetConsoleLoginURL(c.Args().First(), c.Int("duration")*60, subconsole)
	if err != nil {
		log.Errorf("An error ocurred: %s", err)
		os.Exit(1)
	}
	fmt.Println(loginURL)
}
