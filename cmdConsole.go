package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getCmdConsole() *cobra.Command {
	cmd := cobra.Command{
		Use:   "console [environment] [subconsole]",
		Short: "prints a sign-in URL for the AWS web console",
		Run:   actionCmdConsole,
	}

	cmd.Flags().IntVarP(&cfg.Console.Duration, "duration", "d", 8*60, "time in minutes the sign-in is valid")

	return &cmd
}

func actionCmdConsole(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		log.Error("Please specify the name of the environment to create the sign-in URL for")
		os.Exit(1)
	}

	if cfg.Console.Duration*60 < 900 || cfg.Console.Duration*60 > 129600 {
		log.Errorln("Duration parameter must between 15 and 2160 minutes.")
		os.Exit(1)
	}

	subconsole := "console"
	if len(args) == 2 {
		subconsole = args[1]
	}

	loginURL, err := awsCredentials.GetConsoleLoginURL(args[0], cfg.Console.Duration*60, subconsole)
	if err != nil {
		log.Errorf("An error ocurred: %s", err)
		os.Exit(1)
	}
	fmt.Println(loginURL)
}
