package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getCmdDelete() *cobra.Command {
	cmd := cobra.Command{
		Use:   "delete [environment]",
		Short: "delete an AWS environment",
		Run:   actionCmdDelete,
	}

	return &cmd
}

func actionCmdDelete(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Error("Please specify the name of the environment to delete")
		os.Exit(1)
	}

	if _, ok := awsCredentials.Credentials[args[0]]; ok {
		delete(awsCredentials.Credentials, args[0])
		awsCredentials.SaveToFile()
		log.Infof("AWS environment '%s' was successfully deleted.", args[0])
	} else {
		log.Errorf("AWS environment '%s' was not found.", args[0])
		os.Exit(1)
	}
}
