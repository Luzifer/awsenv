package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Luzifer/awsenv/credentials"
	"github.com/Luzifer/awsenv/security"
	"github.com/Luzifer/awsenv/shellsupport"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	_ "github.com/Luzifer/awsenv/shellsupport/bash"
	_ "github.com/Luzifer/awsenv/shellsupport/fish"
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

	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		if len(c.Args()) == 0 || (len(c.Args()) > 0 && c.Args()[0] != "unlock" && c.Args()[0] != "lockagent") {
			if len(c.String("password")) == 0 {
				filename := fmt.Sprintf("%s.lock", c.GlobalString("database"))
				if _, err := os.Stat(filename); os.IsNotExist(err) {
					log.Errorf("No password is available. Use 'unlock' or provide --password.")
					return fmt.Errorf("No password is available. Use 'unlock' or provide --password.")
				}
				pwd, err := security.LoadDatabasePasswordFromLockagent(filename)
				if err != nil {
					log.Errorf("Could not load password from lock-file:\n%s", err)
					return fmt.Errorf("Could not load password from lock-file:\n%s", err)
				}
				password = pwd
			} else {
				password = security.LoadDatabasePasswordFromInput(c.String("password"))
			}

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
		{
			Name:    "list",
			Aliases: []string{"l", "ls"},
			Usage:   "list available AWS environments",
			Action: func(c *cli.Context) {
				for k := range awsCredentials.Credentials {
					fmt.Println(k)
				}
			},
		},
		{
			Name:  "get",
			Usage: "print the AWS crentitals in human readable format",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) {
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
			},
		},
		{
			Name:  "add",
			Usage: "add new AWS credentials",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "access-key,a",
					Usage: "the AWSAccessKey",
				},
				cli.StringFlag{
					Name:  "secret-access-key,s",
					Usage: "the AWSSecretAccessKey",
				},
				cli.StringFlag{
					Name:  "region,r",
					Usage: "the default AWS EC2 region for this credential set",
					Value: "us-east-1",
				},
			},
			Action: func(c *cli.Context) {
				if !c.Args().Present() {
					cli.ShowCommandHelp(c, "add")
					log.Error("Please specify at least the name")
					os.Exit(1)
				}
				keyID := c.String("access-key")
				secretKey := c.String("secret-access-key")
				region := c.String("region")
				var err error
				if len(keyID) == 0 {
					keyID, err = readStdinLine("AWS Access-Key: ")
					if err != nil {
						log.Errorf("An error ocurred: %s", err)
						os.Exit(1)
					}
				}
				if len(secretKey) == 0 {
					secretKey, err = readStdinLine("AWS Secret-Access-Key: ")
					if err != nil {
						log.Errorf("An error ocurred: %s", err)
						os.Exit(1)
					}
				}

				cred := credentials.AWSCredential{
					AWSAccessKeyID:     keyID,
					AWSSecretAccessKey: secretKey,
					AWSRegion:          region,
				}
				awsCredentials.Credentials[c.Args().First()] = cred
				awsCredentials.SaveToFile()
				log.Infof("Credential '%s' has been created", c.Args().First())
			},
		},
		{
			Name:  "delete",
			Usage: "delete an AWS environment",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) {
				if !c.Args().Present() {
					cli.ShowCommandHelp(c, "delete")
					log.Error("Please specify the name of the environment to delete")
					os.Exit(1)
				}
				if _, ok := awsCredentials.Credentials[c.Args().First()]; ok {
					delete(awsCredentials.Credentials, c.Args().First())
					log.Infof("AWS environment '%s' was successfully deleted.", c.Args().First())
				} else {
					log.Errorf("AWS environment '%s' was not found.", c.Args().First())
					os.Exit(1)
				}
			},
		},
		{
			Name:  "shell",
			Usage: "print the AWS credentials in a format for your shell to eval()",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "shell,s",
					Value:  "",
					Usage:  "name of the shell to export for",
					EnvVar: "SHELL",
				},
				cli.BoolTFlag{
					Name:  "export,x",
					Usage: "Adds proper export options for your shell",
				},
			},
			Action: func(c *cli.Context) {
				if len(c.String("shell")) == 0 {
					log.Errorf("Could not determine your shell. Please provide --shell")
					os.Exit(1)
				}
				s := strings.Split(c.String("shell"), "/")
				shell := s[len(s)-1]

				log.Debugf("Found shell '%s'", shell)

				handler, err := shellsupport.GetShellHandler(shell)
				if err != nil {
					log.Errorf("Could not find a handler for '%s' shell", shell)
					os.Exit(1)
				}

				if !c.Args().Present() {
					log.Errorf("Please specify the enviroment to load")
					os.Exit(1)
				}

				if a, ok := awsCredentials.Credentials[c.Args().First()]; ok {
					fmt.Println(strings.Join(handler(a, c.Bool("export")), "\n"))
					os.Exit(0)
				}

				log.Errorf("Could not find environment '%s'", c.Args().First())
				os.Exit(1)
			},
		},
		{
			Name:  "lock",
			Usage: "lock the database",
			Action: func(c *cli.Context) {
				password.KillLockAgent(fmt.Sprintf("%s.lock", c.GlobalString("database")))
			},
		},
		{
			Name:  "unlock",
			Usage: "unlock the database (this will store your password on the disk!)",
			Action: func(c *cli.Context) {
				var pwd *security.DatabasePassword
				var err error

				if len(c.GlobalString("password")) > 0 {
					pwd = security.LoadDatabasePasswordFromInput(c.String("password"))
				} else {
					line, err := readStdinLine("Password: ")
					if err != nil {
						log.Errorln(err)
					}
					pwd = security.LoadDatabasePasswordFromInput(line)
				}

				err = pwd.SpawnLockAgent(fmt.Sprintf("%s.lock", c.GlobalString("database")))
				if err != nil {
					log.Errorf("Unable to spawn lockagent: %s", err)
					return
				}

				fmt.Println("Database unlocked.")
			},
		},
		{
			Name:  "console",
			Usage: "prints a sign-in URL for the AWS web console",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "duration,d",
					Value: 8 * 60,
					Usage: "time in minutes the sign-in is valid",
				},
			},
			Action: func(c *cli.Context) {
				if !c.Args().Present() {
					cli.ShowCommandHelp(c, "console")
					log.Error("Please specify the name of the environment to create the sign-in URL for")
					os.Exit(1)
				}

				if c.Int("duration")*60 < 900 || c.Int("duration")*60 > 129600 {
					log.Errorln("Duration parameter must between 15 and 2160 minutes.")
					os.Exit(1)
				}

				loginURL, err := awsCredentials.GetConsoleLoginURL(c.Args().First(), c.Int("duration")*60)
				if err != nil {
					log.Errorf("An error ocurred: %s", err)
					os.Exit(1)
				}
				fmt.Println(loginURL)
			},
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		switch command {
		case "lockagent":
			RunLockagent()
		default:
			fmt.Fprintf(c.App.Writer, "No help topic for '%v'\n", command)
			return
		}
	}

	app.Run(os.Args)
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
