package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func getCmdGet() cli.Command {
	return cli.Command{
		Name:   "get",
		Usage:  "print the AWS crentitals in human readable format",
		Flags:  []cli.Flag{},
		Action: actionCmdGet,
	}
}

func actionCmdGet(c *cli.Context) {
	if !c.Args().Present() {
		cli.ShowCommandHelp(c, "get")
		log.Error("Please specify the name of the environment to get")
		os.Exit(1)
	}

	if a, ok := awsCredentials.Credentials[c.Args().First()]; ok {
		fmt.Printf("Credentials for the '%s' environment:\n", c.Args().First())
		fmt.Printf(" AWS Access-Key:        %s\n", a.AWSAccessKeyID)
		fmt.Printf(" AWS Secret-Access-Key: %s\n", a.AWSSecretAccessKey)
		fmt.Printf(" AWS EC2-Region:        %s\n", a.AWSRegion)
		os.Exit(0)
	}

	log.Errorf("Could not find environment '%s'", c.Args().First())
	os.Exit(1)
}
