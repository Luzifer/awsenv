package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getCmdLock() *cobra.Command {
	cmd := cobra.Command{
		Use:   "lock",
		Short: "lock the database",
		Run:   actionCmdLock,
	}

	return &cmd
}

func actionCmdLock(cmd *cobra.Command, args []string) {
	err := password.KillLockAgent(fmt.Sprintf("%s.lock", cfg.Database))
	if err != nil {
		log.Error("Unable to kill the lock agent. Is it running?")
	}
}
