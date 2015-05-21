# Security considerations

## Database unlock
The database password is stored in a tiny HTTP server listening on a random port on localhost after the database is unlocked. This server generates a security token every 10 seconds and writes it to the disk next to the password database. This way a backup does include only the security token which was exchanged in the meantime but not the password itself.

Also this method requires more effort from anyone trying to extract the password itself as he has to get the security token and query the password from the server within those 10s. A working method to extract the password might be to read the memory store of the server and extract the password stored in memory.

## Multi-Factor-Authentication
For the command `awsenv console` the AWS [GetFederationToken](http://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html) API is used. This API does not have the chance to require any MFA token from the user. Instead the URL is generated and after opening the URL the user is logged in directly.

If the MFA device should be requested the application would need to have a fixed IAM role in the target account which can be assumed as the MFA token login is only possible with [AssumeRole](http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html). In that case anyone using awsenv would receive the same rights the fixed IAM role has in that account. Should the user be able to choose the role to assume in order to gain different rights he would be able to guess role names with higher access like "admin" or "superuser" and use that instead of the "user" role he should assume.

Additional the usage of MFA devices is not trivial in AWS as one account can has multiple devices assigned. There is an API to retrieve the IDs of all associated MFA devices but in the request there has to be a matching MFA device ID to the token the user just entered. The token can not be validated against any API.

Because of this the `awsenv` command does not support login to AWS web console using MFA devices. For that reason the database should kept in a locked state as everyone gaining physical access to the machine otherwise can impersonate the owner of the credentials and gain access to that AWS account.

One solution against this would be to secure the whole `awsenv` command using an own MFA token to be entered with every request made by the user.

## Encryption on disk
The credential database stored on the users computer at `~/.config/awsenv` is encrypted using an AES encryption with a 32 byte key commonly known as AES256. The initialization vector is generated randomly at every save of the database and does not make direct attacks to the credential store possible.

Also it is not possible to decrypt the raw value stored in the credential database with just knowing the password itself as the key used to encrypt the database is not the password itself but [derived](https://github.com/Luzifer/awsenv/blob/master/security/databasePassword.go#L56) from the password.
