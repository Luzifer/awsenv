package main

import (
	"fmt"

	"github.com/Luzifer/awsenv/security"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func getCmdUnlock() cli.Command {
	return cli.Command{
		Name:   "unlock",
		Usage:  "unlock the database (this will store your password on the disk!)",
		Action: actionCmdUnlock,
	}
}

func actionCmdUnlock(c *cli.Context) {
	var pwd *security.DatabasePassword
	var err error

	if len(c.GlobalString("password")) > 0 {
		pwd = security.LoadDatabasePasswordFromInput(c.String("password"))
	} else {
		line, err := readStdinLine("Password: ")
		if err != nil {
			log.Errorln(err)
		}
		pwd = security.LoadDatabasePasswordFromInput(line)
	}

	err = pwd.SpawnLockAgent(fmt.Sprintf("%s.lock", c.GlobalString("database")))
	if err != nil {
		log.Errorf("Unable to spawn lockagent: %s", err)
		return
	}

	fmt.Println("Database unlocked.")
}
