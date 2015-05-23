package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func getCmdLock() cli.Command {
	return cli.Command{
		Name:   "lock",
		Usage:  "lock the database",
		Action: actionCmdLock,
	}
}

func actionCmdLock(c *cli.Context) {
	err := password.KillLockAgent(fmt.Sprintf("%s.lock", c.GlobalString("database")))
	if err != nil {
		log.Error("Unable to kill the lock agent. Is it running?")
	}
}
