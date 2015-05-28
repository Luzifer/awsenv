package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getCmdList() *cobra.Command {
	cmd := cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "list available AWS environments",
		Run:     actionCmdList,
	}
	return &cmd
}

func actionCmdList(cmd *cobra.Command, args []string) {
	for k := range awsCredentials.Credentials {
		fmt.Println(k)
	}
}
