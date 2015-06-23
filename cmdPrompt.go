package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func getCmdPrompt() *cobra.Command {
	cmd := cobra.Command{
		Use:   "prompt",
		Short: "echos the name of the currently set env for use in prompts",
		Run:   actionCmdPrompt,
	}
	return &cmd
}

func actionCmdPrompt(cmd *cobra.Command, args []string) {
	for k, v := range awsCredentials.Credentials {
		if v.AWSAccessKeyID == os.Getenv("AWS_ACCESS_KEY") {
			fmt.Printf(k)
		}
	}
}
