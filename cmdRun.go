package main

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func getCmdRun() *cobra.Command {
	cmd := cobra.Command{
		Use:   "run [environment] [command]",
		Short: "runs the given command with populated environment variables",
		Run:   actionCmdRun,
	}
	return &cmd
}

func actionCmdRun(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		log.Error("Please specify the name of the environment to get and the command to run")
		os.Exit(1)
	}

	if a, ok := awsCredentials.Credentials[args[0]]; ok {
		cmd := exec.Command(args[1], args[2:]...)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		env := envListToMap(os.Environ())
		env["AWS_PROFILE"] = a.AWSProfile
		env["AWS_ACCESS_KEY"] = a.AWSAccessKeyID
		env["AWS_ACCESS_KEY_ID"] = a.AWSAccessKeyID
		env["AWS_SECRET_KEY"] = a.AWSSecretAccessKey
		env["AWS_SECRET_ACCESS_KEY"] = a.AWSSecretAccessKey
		env["EC2_REGION"] = a.AWSRegion
		env["AWS_REGION"] = a.AWSRegion
		env["AWS_DEFAULT_REGION"] = a.AWSRegion

		cmd.Env = envMapToList(env)

		err := cmd.Run()
		switch err.(type) {
		case nil:
			os.Exit(0)
		case *exec.ExitError:
			log.Println("Unclean exit with exit-code != 0")
			os.Exit(1)
		default:
			log.Printf("An unknown error ocurred: %s", err)
			os.Exit(2)
		}
	}

	log.Errorf("Could not find environment '%s'", args[0])
	os.Exit(1)
}

func envListToMap(list []string) map[string]string {
	out := map[string]string{}
	for _, entry := range list {
		if len(entry) == 0 {
			continue
		}

		parts := strings.SplitN(entry, "=", 2)
		out[parts[0]] = parts[1]
	}
	return out
}

func envMapToList(envMap map[string]string) []string {
	out := []string{}
	for k, v := range envMap {
		out = append(out, k+"="+v)
	}
	return out
}
