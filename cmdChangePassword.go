package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bgentry/speakeasy"
	"github.com/spf13/cobra"
)

func getCmdChangePassword() *cobra.Command {
	cmd := cobra.Command{
		Use:   "passwd",
		Short: "change the password of the database",
		Run:   actionCmdChangePassword,
	}

	return &cmd
}

func actionCmdChangePassword(cmd *cobra.Command, args []string) {
	passwd, err := speakeasy.Ask("New Password: ")
	if err != nil {
		log.Errorf("Unable to read password: %s", err)
	}
	repeatPasswd, err := speakeasy.Ask("Please repeat: ")
	if err != nil {
		log.Errorf("Unable to read password: %s", err)
	}

	if passwd != repeatPasswd {
		log.Errorln("Passwords does not match. Not executing change.")
		return
	}

	if err := awsCredentials.UpdatePassword(passwd); err != nil {
		log.Error("Unable to change the password.")
		return
	}

	log.Println("Password has been changed, locking database now.")
	actionCmdLock(cmd, args)
}
