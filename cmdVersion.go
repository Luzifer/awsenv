package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func getCmdVersion() *cobra.Command {
	cmd := cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "prints the current version of awsenv",
		Run:     actionCmdVersion,
	}
  
	return &cmd
}

func actionCmdVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("awsenv version %s\n", version)
}
