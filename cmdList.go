package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func getCmdList() cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"l", "ls"},
		Usage:   "list available AWS environments",
		Action: func(c *cli.Context) {
			for k := range awsCredentials.Credentials {
				fmt.Println(k)
			}
		},
	}
}
