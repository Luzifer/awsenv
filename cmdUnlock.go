package main

import (
	"fmt"

	"github.com/Luzifer/awsenv/security"
	log "github.com/Sirupsen/logrus"
	"github.com/bgentry/speakeasy"
	"github.com/spf13/cobra"
)

func getCmdUnlock() *cobra.Command {
	cmd := cobra.Command{
		Use:   "unlock",
		Short: "unlock the database",
		Run:   actionCmdUnlock,
	}

	cmd.Flags().StringVar(&cfg.LockAgent.Timeout, "timeout", "30m", "Lock awsenv after this time (set to 0 to disable)")

	return &cmd
}

func actionCmdUnlock(cmd *cobra.Command, args []string) {
	var pwd *security.DatabasePassword
	var err error

	if len(cfg.Password) > 0 {
		pwd = security.LoadDatabasePasswordFromInput(cfg.Password)
	} else {
		line, err := speakeasy.Ask("Password: ")
		if err != nil {
			log.Errorln(err)
		}
		pwd = security.LoadDatabasePasswordFromInput(line)
	}

	err = pwd.SpawnLockAgent(fmt.Sprintf("%s.lock", cfg.Database), cfg.LockAgent.Timeout)
	if err != nil {
		log.Errorf("Unable to spawn lockagent: %s", err)
		return
	}

	fmt.Println("Database unlocked.")
}
