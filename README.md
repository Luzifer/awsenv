[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/awsenv)](https://goreportcard.com/report/github.com/Luzifer/awsenv)
![](https://badges.fyi/github/license/Luzifer/awsenv)
![](https://badges.fyi/github/downloads/Luzifer/awsenv)
![](https://badges.fyi/github/latest-release/Luzifer/awsenv)

# Luzifer / awsenv

awsenv is intended as a local credential store for people using more than one AWS account at the same time.

For security considerations about this software please refer to the [security.md](https://github.com/Luzifer/awsenv/blob/master/security.md) file in this repository.

## Features
- Secure storage of credentials (AES256)
- No more access when credential store is "locked"
- Export credentials for your shells eval function
- Amazon STS support to open the web-console without login-hazzle

## Installation

### From source

```
go get -u github.com/Luzifer/awsenv
```

### From binary

1. Go to the [releases page](https://github.com/Luzifer/awsenv/releases)
2. Download the binary for your system and put into your `$PATH`

## Supported shells

- bash / zsh
  - Put this function into your `~/.bashrc` / `~/.zshrc` and you can access your environments using `set_aws <name>`

```bash
function set_aws {
  eval $(awsenv shell $1)
}
function login_aws {
  open $(awsenv console $1)
}
```

- fish
  - Put this function into `~/.config/fish/functions/set_aws.fish` and you can access your environments using `set_aws <name>`

```fish
function set_aws --description 'Set the AWS environment variables' --argument AWS_ENV
	eval (awsenv shell $AWS_ENV)
end
function login_aws --description 'Open browser with AWS console' --argument AWS_ENV
	open (awsenv console $AWS_ENV)
end
```

## Sample workthrough

### Installation
```bash
$ curl -sSLfo awsenv https://github.com/Luzifer/awsenv/releases/download/v0.11.1/awsenv_linux_amd64
$ chmod 0755 awsenv
$ sudo mv awsenv /usr/local/bin/
```

### Adding an environment and using it
```bash
# We can not list because the credentials are locked
$ awsenv list
ERRO[0000] No password is available. Use 'unlock' or provide --password.

# Unlock the credentials (now the password is set for later)
$ awsenv unlock
Password: demo

# We can now list without errors but have no environments
$ awsenv list

# Lets add an environment
$ awsenv add --region eu-west-1 demoenv
AWS Access-Key: myaccesskey
AWS Secret-Access-Key: mysecretkey
INFO[0010] Credential 'demoenv' has been created

# Now we can list the environment we just created
$ awsenv list
demoenv

# With the get command we can display the information
$ awsenv get demoenv
Credentials for the 'demoenv' environment:
 AWS Access-Key:        myaccesskey
 AWS Secret-Access-Key: mysecretkey
 AWS EC2-Region:        eu-west-1

# The lock command will secure the credentials again
$ awsenv lock
$ awsenv get demoenv
ERRO[0000] No password is available. Use 'unlock' or provide --password.

# We need to unlock it with the same credentials
$ awsenv unlock
Password: demo
$ awsenv get demoenv
Credentials for the 'demoenv' environment:
 AWS Access-Key:        myaccesskey
 AWS Secret-Access-Key: mysecretkey
 AWS EC2-Region:        eu-west-1

# We're currently working in a bash without AWS ENV vars
$ env | grep AWS

# But we can load them using the set_aws function
$ set_aws demoenv
$ env | grep AWS
AWS_SECRET_ACCESS_KEY=mysecretkey
AWS_ACCESS_KEY_ID=myaccesskey
AWS_ACCESS_KEY=myaccesskey
AWS_SECRET_KEY=mysecretkey

# Now the prompt command can tell you which env is set
$ awsenv prompt
demoenv

# You also can run commands with AWS crentials directly
$ awsenv run demoenv -- env | grep AWS
AWS_ACCESS_KEY_ID=myaccesskey
AWS_SECRET_ACCESS_KEY=mysecretkey
AWS_ACCESS_KEY=myaccesskey
AWS_SECRET_KEY=mysecretkey
AWS_REGION=us-east-1
AWS_DEFAULT_REGION=us-east-1

# Lets try to unlock with a wrong password
$ awsenv lock
$ awsenv unlock
Password: fooo

# The database is now not readable for us
$ awsenv l
ERRO[0000] Unable to read credential database

# As soon as we unlock with the right password it works again
$ awsenv unlock
Password: demo
$ awsenv l
demoenv
```
