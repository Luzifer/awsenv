package main

import (
	"os"

	"github.com/Luzifer/awsenv/credentials"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getCmdAdd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "add [env name]",
		Short: "add a new AWS credential environment",
		Run:   actionCmdAdd,
	}

	cmd.Flags().StringVarP(&cfg.Add.AccessKey, "access-key", "a", "", "the AWSAccessKey")
	cmd.Flags().StringVarP(&cfg.Add.SecretAccessKey, "secret-access-key", "s", "", "the AWSSecretAccessKey")
	cmd.Flags().StringVarP(&cfg.Add.Region, "cfg.Add.Region", "r", "us-east-1", "the default AWS EC2 cfg.Add.Region for this credential set")

	return &cmd
}

func actionCmdAdd(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmd.Usage()
		log.Error("Please specify at least the name")
		os.Exit(1)
	}

	var err error
	if len(cfg.Add.AccessKey) == 0 {
		cfg.Add.AccessKey, err = readStdinLine("AWS Access-Key: ")
		if err != nil {
			log.Errorf("An error ocurred: %s", err)
			os.Exit(1)
		}
	}
	if len(cfg.Add.SecretAccessKey) == 0 {
		cfg.Add.SecretAccessKey, err = readStdinLine("AWS Secret-Access-Key: ")
		if err != nil {
			log.Errorf("An error ocurred: %s", err)
			os.Exit(1)
		}
	}

	cred := credentials.AWSCredential{
		AWSAccessKeyID:     cfg.Add.AccessKey,
		AWSSecretAccessKey: cfg.Add.SecretAccessKey,
		AWSRegion:          cfg.Add.Region,
	}
	awsCredentials.Credentials[args[0]] = cred
	err = awsCredentials.SaveToFile()
	if err != nil {
		log.Errorf("Unable to save the credentials to database '%s'", cfg.Database)
		os.Exit(1)
	}
	log.Infof("Credential '%s' has been created", args[0])
}
