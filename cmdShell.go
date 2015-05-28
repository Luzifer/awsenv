package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Luzifer/awsenv/shellsupport"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	_ "github.com/Luzifer/awsenv/shellsupport/bash"
	_ "github.com/Luzifer/awsenv/shellsupport/fish"
)

func getCmdShell() *cobra.Command {
	cmd := cobra.Command{
		Use:   "shell [environment]",
		Short: "print the AWS credentials in a format for your shell to eval()",
		Run:   actionCmdShell,
	}

	cmd.Flags().StringVarP(&cfg.Shell.Shell, "shell", "s", os.Getenv("SHELL"), "name of the shell to export for")
	cmd.Flags().BoolVarP(&cfg.Shell.Export, "export", "x", true, "Adds proper export options for your shell")

	return &cmd
}

func actionCmdShell(cmd *cobra.Command, args []string) {
	if len(cfg.Shell.Shell) == 0 {
		log.Errorf("Could not determine your shell. Please provide --shell")
		os.Exit(1)
	}
	s := strings.Split(cfg.Shell.Shell, "/")
	shell := s[len(s)-1]

	log.Debugf("Found shell '%s'", shell)

	handler, err := shellsupport.GetShellHandler(shell)
	if err != nil {
		log.Errorf("Could not find a handler for '%s' shell", shell)
		os.Exit(1)
	}

	if len(args) < 1 {
		log.Errorf("Please specify the enviroment to load")
		os.Exit(1)
	}

	if a, ok := awsCredentials.Credentials[args[0]]; ok {
		fmt.Println(strings.Join(handler(a, cfg.Shell.Export), "\n"))
		os.Exit(0)
	}

	log.Errorf("Could not find environment '%s'", args[0])
	os.Exit(1)
}
