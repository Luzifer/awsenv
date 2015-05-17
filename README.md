# Luzifer / awsenv

awsenv is intended as a local credential store for people using more than one AWS account at the same time.

## Features
- Secure storage of credentials (AES256)
- No more access when credential store is "locked"
- Export credentials for your shells eval function

## Planned features
- More secure unlock mechanism
  - No more storing the password in plain text on the disk
- Amazon STS support to open the web-console without login-hazzle

## Supported shells

- bash

```bash
eval $(awsenv shell private)
```

- fish

```fish
eval (awsenv shell private)
```
