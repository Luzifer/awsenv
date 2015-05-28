package fish // import "github.com/Luzifer/awsenv/shellsupport/fish"

import (
	"fmt"

	"github.com/Luzifer/awsenv/credentials"
	"github.com/Luzifer/awsenv/shellsupport"
)

func init() {
	shellsupport.RegisterShellHandler("fish", fishShellHandler)
}

func fishShellHandler(c credentials.AWSCredential, export bool) []string {
	flags := "-g"
	if export {
		flags = "-gx"
	}

	return []string{
		fmt.Sprintf("set %s AWS_ACCESS_KEY %s;", flags, c.AWSAccessKeyID),
		fmt.Sprintf("set %s AWS_ACCESS_KEY_ID %s;", flags, c.AWSAccessKeyID),
		fmt.Sprintf("set %s AWS_SECRET_KEY %s;", flags, c.AWSSecretAccessKey),
		fmt.Sprintf("set %s AWS_SECRET_ACCESS_KEY %s;", flags, c.AWSSecretAccessKey),
		fmt.Sprintf("set %s EC2_REGION %s;", flags, c.AWSRegion),
		fmt.Sprintf("set %s AWS_REGION %s", flags, c.AWSRegion),
	}
}
