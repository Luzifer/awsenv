package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getCmdGet() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get [environment]",
		Short: "print the AWS crentitals in human readable format",
		Run:   actionCmdGet,
	}
	return &cmd
}

func actionCmdGet(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Error("Please specify the name of the environment to get")
		os.Exit(1)
	}

	if a, ok := awsCredentials.Credentials[args[0]]; ok {
		fmt.Printf("Credentials for the '%s' environment:\n", args[0])
		fmt.Printf(" AWS Profile:           %s\n", a.AWSProfile)
		fmt.Printf(" AWS Access-Key:        %s\n", a.AWSAccessKeyID)
		fmt.Printf(" AWS Secret-Access-Key: %s\n", a.AWSSecretAccessKey)
		fmt.Printf(" AWS EC2-Region:        %s\n", a.AWSRegion)
		os.Exit(0)
	}

	log.Errorf("Could not find environment '%s'", args[0])
	os.Exit(1)
}
