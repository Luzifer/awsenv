package shellsupport // import "github.com/Luzifer/awsenv/shellsupport"

import (
	"fmt"

	"github.com/Luzifer/awsenv/credentials"
)

var shellHandlers map[string]ShellHandler

type ShellHandler func(credentials.AWSCredential, bool) []string

func RegisterShellHandler(shell string, fn ShellHandler) {
	if _, exists := shellHandlers[shell]; exists {
		panic(fmt.Errorf("Tried to overwrite '%s' handler.", shell))
	}
	shellHandlers[shell] = fn
}

func init() {
	shellHandlers = make(map[string]ShellHandler)
}

func GetShellHandler(shell string) (ShellHandler, error) {
	if handler, ok := shellHandlers[shell]; ok {
		return handler, nil
	}
	return nil, fmt.Errorf("Could not find a handler for '%s' shell", shell)
}
