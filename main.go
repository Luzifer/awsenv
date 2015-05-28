package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Luzifer/awsenv/credentials"
	"github.com/Luzifer/awsenv/security"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	version = "0.4.1"
)

var (
	password       *security.DatabasePassword
	awsCredentials *credentials.AWSCredentialStore
	cfg            = &config{}
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	// Do not route special commands into cobra logic
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "lockagent":
			runLockagent()
			os.Exit(0)
		}
	}

	app := cobra.Command{
		Use:   "awsenv",
		Short: "manage different AWS envs on your system",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cfg.Debug {
				log.SetLevel(log.DebugLevel)
			}

			// Load the password if command is not unlock
			if !strings.Contains("unlock version", cmd.Name()) {

				if len(cfg.Password) > 0 {
					// If a password was provided, use that one
					password = security.LoadDatabasePasswordFromInput(cfg.Password)
				} else {
					// If the token file exists a lockagent should be running and we can use
					// the password stored in that logagent
					filename := fmt.Sprintf("%s.lock", cfg.Database)
					if _, err := os.Stat(filename); os.IsNotExist(err) {
						log.Errorf("No password is available. Use 'unlock' or provide --password.")
						os.Exit(1)
					}
					pwd, err := security.LoadDatabasePasswordFromLockagent(filename)
					if err != nil {
						log.Errorf("Could not load password from lock-file:\n%s", err)
						os.Exit(1)
					}
					password = pwd
				}

				// As we got a password now try to load the database with that password or
				// Create a new one if the encrypted storage file is not available
				if _, err := os.Stat(cfg.Database); os.IsNotExist(err) {
					awsCredentials = credentials.New(cfg.Database, password)
				} else {
					s, err := credentials.FromFile(cfg.Database, password)
					if err != nil {
						log.Error("Unable to read credential database")
						os.Exit(1)
					}
					awsCredentials = s
				}
			}
		},
	}

	app.Flags().StringVarP(&cfg.Password, "password", "p", os.Getenv("AWSENV_PASSWORD"), "password to en/decrypt the database")
	app.Flags().StringVar(&cfg.Database, "database", strings.Join([]string{os.Getenv("HOME"), ".config/awsenv"}, "/"), "storage location of the database")
	app.Flags().BoolVarP(&cfg.Debug, "debug", "d", false, "print debug information")

	app.AddCommand(
		getCmdAdd(),
		getCmdConsole(),
		getCmdDelete(),
		getCmdGet(),
		getCmdList(),
		getCmdLock(),
		getCmdShell(),
		getCmdUnlock(),
		getCmdVersion(),
	)

	_ = app.Execute()
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
