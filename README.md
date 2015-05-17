# Luzifer / awsenv

awsenv is intended as a local credential store for people using more than one AWS account at the same time.

**Please remember: This is alpha-software!**

## Features
- Secure storage of credentials (AES256)
- No more access when credential store is "locked"
- Export credentials for your shells eval function

## Planned features
- More secure unlock mechanism
  - No more storing the password in plain text on the disk
- Amazon STS support to open the web-console without login-hazzle

## Installation

### From source

```
go get github.com/Luzifer/awsenv
```

### From binary

1. Go to the [GoBuilder builds page](https://gobuilder.me/github.com/Luzifer/awsenv)
2. Download the ZIP for your system
3. Unpack and put the `awsenv` file into your `$PATH`

## Supported shells

- bash
  - Put this function into your `~/.bashrc` and you can access your environments using `set_aws <name>`

```bash
function set_aws {
  eval $(awsenv shell $1)
}
```

- fish
  - Put this function into `~/.config/fish/functions/set_aws.fish` and you can access your environments using `set_aws <name>`

```fish
function set_aws --description 'Set the AWS environment variables' --argument AWS_ENV
	eval (awsenv shell $AWS_ENV)
end
```
