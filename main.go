package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Luzifer/awsenv/credentials"
	"github.com/Luzifer/awsenv/security"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

var password *security.DatabasePassword
var awsCredentials *credentials.AWSCredentialStore

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	app := cli.NewApp()
	app.Name = "awsenv"
	app.Usage = "manage different AWS envs on your system"
	app.Version = "0.3.0"

	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		// Load the password if command is not unlock or lockagent
		if len(c.Args()) == 0 || (!strings.Contains("unlock||lockagent", c.Args()[0])) {

			if len(c.GlobalString("password")) > 0 {
				// If a password was provided, use that one
				password = security.LoadDatabasePasswordFromInput(c.String("password"))
			} else {
				// If the token file exists a lockagent should be running and we can use
				// the password stored in that logagent
				filename := fmt.Sprintf("%s.lock", c.GlobalString("database"))
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					log.Errorf("No password is available. Use 'unlock' or provide --password.")
					return err
				}
				pwd, err := security.LoadDatabasePasswordFromLockagent(filename)
				if err != nil {
					log.Errorf("Could not load password from lock-file:\n%s", err)
					return err
				}
				password = pwd
			}

			// As we got a password now try to load the database with that password or
			// Create a new one if the encrypted storage file is not available
			if _, err := os.Stat(c.GlobalString("database")); os.IsNotExist(err) {
				awsCredentials = credentials.New(c.GlobalString("database"), password)
			} else {
				s, err := credentials.FromFile(c.GlobalString("database"), password)
				if err != nil {
					log.Error("Unable to read credential database")
					return err
				}
				awsCredentials = s
			}
		}

		return nil
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug,d",
			Usage: "print debug information",
		},
		cli.StringFlag{
			Name:  "database",
			Value: strings.Join([]string{os.Getenv("HOME"), ".config/awsenv"}, "/"),
			Usage: "storage location of the database",
		},
		cli.StringFlag{
			Name:   "password,p",
			Value:  "",
			Usage:  "password to en/decrypt the database",
			EnvVar: "AWSENV_PASSWORD",
		},
	}

	app.Commands = []cli.Command{
		getCmdList(),
		getCmdGet(),
		getCmdAdd(),
		getCmdDelete(),
		getCmdShell(),
		getCmdLock(),
		getCmdUnlock(),
		getCmdConsole(),
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		switch command {
		case "lockagent":
			runLockagent()
		default:
			fmt.Fprintf(c.App.Writer, "No help topic for '%v'\n", command)
			return
		}
	}

	_ = app.Run(os.Args)
}

func readStdinLine(prompt string) (string, error) {
	fmt.Printf(prompt)
	bio := bufio.NewReader(os.Stdin)
	line, _, err := bio.ReadLine()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(line)), nil
}
