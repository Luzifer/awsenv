package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func getCmdDelete() cli.Command {
	return cli.Command{
		Name:   "delete",
		Usage:  "delete an AWS environment",
		Flags:  []cli.Flag{},
		Action: actionCmdDelete,
	}
}

func actionCmdDelete(c *cli.Context) {
	if !c.Args().Present() {
		cli.ShowCommandHelp(c, "delete")
		log.Error("Please specify the name of the environment to delete")
		os.Exit(1)
	}
	if _, ok := awsCredentials.Credentials[c.Args().First()]; ok {
		delete(awsCredentials.Credentials, c.Args().First())
		log.Infof("AWS environment '%s' was successfully deleted.", c.Args().First())
	} else {
		log.Errorf("AWS environment '%s' was not found.", c.Args().First())
		os.Exit(1)
	}
}
