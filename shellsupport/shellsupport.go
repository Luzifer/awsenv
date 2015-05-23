package shellsupport // import "github.com/Luzifer/awsenv/shellsupport"

import (
	"fmt"

	"github.com/Luzifer/awsenv/credentials"
)

var shellHandlers map[string]ShellHandler

// ShellHandler defines the method stub to implement for every supported shell
type ShellHandler func(credentials.AWSCredential, bool) []string

// RegisterShellHandler registers a ShellHandler for a given shell. The shell
// name must be the last part of the path. So for example for /bin/bash it would
// be "bash"
func RegisterShellHandler(shell string, fn ShellHandler) {
	if _, exists := shellHandlers[shell]; exists {
		panic(fmt.Errorf("Tried to overwrite '%s' handler.", shell))
	}
	shellHandlers[shell] = fn
}

func init() {
	shellHandlers = make(map[string]ShellHandler)
}

// GetShellHandler returns a ShellHandler for the given shell name
func GetShellHandler(shell string) (ShellHandler, error) {
	if handler, ok := shellHandlers[shell]; ok {
		return handler, nil
	}
	return nil, fmt.Errorf("Could not find a handler for '%s' shell", shell)
}
