package main

import (
	"os"

	"github.com/Luzifer/awsenv/credentials"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func getCmdAdd() cli.Command {
	return cli.Command{
		Name:  "add",
		Usage: "add new AWS credentials",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "access-key,a",
				Usage: "the AWSAccessKey",
			},
			cli.StringFlag{
				Name:  "secret-access-key,s",
				Usage: "the AWSSecretAccessKey",
			},
			cli.StringFlag{
				Name:  "region,r",
				Usage: "the default AWS EC2 region for this credential set",
				Value: "us-east-1",
			},
		},
		Action: actionCmdAdd,
	}
}

func actionCmdAdd(c *cli.Context) {
	if !c.Args().Present() {
		cli.ShowCommandHelp(c, "add")
		log.Error("Please specify at least the name")
		os.Exit(1)
	}
	keyID := c.String("access-key")
	secretKey := c.String("secret-access-key")
	region := c.String("region")
	var err error
	if len(keyID) == 0 {
		keyID, err = readStdinLine("AWS Access-Key: ")
		if err != nil {
			log.Errorf("An error ocurred: %s", err)
			os.Exit(1)
		}
	}
	if len(secretKey) == 0 {
		secretKey, err = readStdinLine("AWS Secret-Access-Key: ")
		if err != nil {
			log.Errorf("An error ocurred: %s", err)
			os.Exit(1)
		}
	}

	cred := credentials.AWSCredential{
		AWSAccessKeyID:     keyID,
		AWSSecretAccessKey: secretKey,
		AWSRegion:          region,
	}
	awsCredentials.Credentials[c.Args().First()] = cred
	err = awsCredentials.SaveToFile()
	if err != nil {
		log.Errorf("Unable to save the credentials to database '%s'", c.GlobalString("database"))
		os.Exit(1)
	}
	log.Infof("Credential '%s' has been created", c.Args().First())
}
