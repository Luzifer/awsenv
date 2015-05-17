package bash // import "github.com/Luzifer/awsenv/shellsupport/bash"

import (
	"fmt"

	"github.com/Luzifer/awsenv/credentials"
	"github.com/Luzifer/awsenv/shellsupport"
)

func init() {
	shellsupport.RegisterShellHandler("bash", bashShellHandler)
}

func bashShellHandler(c credentials.AWSCredential, export bool) []string {
	flags := ""
	if export {
		flags = "export"
	}

	return []string{
		fmt.Sprintf("%s AWS_ACCESS_KEY=%s;", flags, c.AWSAccessKeyID),
		fmt.Sprintf("%s AWS_ACCESS_KEY_ID=%s;", flags, c.AWSAccessKeyID),
		fmt.Sprintf("%s AWS_SECRET_KEY=%s;", flags, c.AWSSecretAccessKey),
		fmt.Sprintf("%s AWS_SECRET_ACCESS_KEY=%s;", flags, c.AWSSecretAccessKey),
		fmt.Sprintf("%s EC2_REGION=%s;", flags, c.AWSRegion),
	}
}
