package security // import "github.com/Luzifer/awsenv/security"

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	openssl "github.com/Luzifer/go-openssl"
)

// DatabasePassword stores a non retrievable password and includes the
// encryption methods to store data protected with AES256
type DatabasePassword struct {
	password string
}

// LoadDatabasePasswordFromInput creates a new DatabasePassword store with the
// given password
func LoadDatabasePasswordFromInput(input string) *DatabasePassword {
	return &DatabasePassword{
		password: input,
	}
}

// LoadDatabasePasswordFromLockagent calls the lock agent of awsenv to get the
// password stored in memory for secure retrieval
func LoadDatabasePasswordFromLockagent(filename string) (*DatabasePassword, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	communication := strings.Split(string(buf), "::")
	port := communication[1]
	token := communication[0]

	r, _ := http.NewRequest("GET", fmt.Sprintf("http://127.0.0.1:%s/password", port), nil)
	r.Header.Add("Token", token)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	pb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &DatabasePassword{
		password: strings.TrimSpace(string(pb)),
	}, nil
}

// SpawnLockAgent spawns and detaches from a new lockagent
func (p *DatabasePassword) SpawnLockAgent(filename, agentTimeout string) error {
	var err error
	proc := exec.Command(os.Args[0], "lockagent")
	proc.Env = []string{
		fmt.Sprintf("DBFILE=%s", filename),
		fmt.Sprintf("TIMEOUT=%s", agentTimeout),
		fmt.Sprintf("PASSWD=%s", p.password),
	}
	err = proc.Start()
	if err != nil {
		return err
	}
	err = proc.Process.Release()
	if err != nil {
		return err
	}
	return nil
}

// KillLockAgent sends the lockagent a kill command using its HTTP interface
func (p *DatabasePassword) KillLockAgent(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	communication := strings.Split(string(buf), "::")
	port := communication[1]
	token := communication[0]

	r, _ := http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:%s/kill", port), nil)
	r.Header.Add("Token", token)
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()

	return nil
}

// Encrypt uses the password stored in the DatabasePassword storage to encrypt
// the given message with AES256 and returns the encrypted message for storing
func (p *DatabasePassword) Encrypt(in []byte) ([]byte, error) {
	return openssl.New().EncryptString(p.password, string(in))
}

// Decrypt takes an encrypted message from Encrypt and decrypts it with a
// corresponding method
func (p *DatabasePassword) Decrypt(in []byte) ([]byte, error) {
	out, err := openssl.New().DecryptString(p.password, string(in))
	if err == nil {
		// Decrypt with OpenSSL compatible encryption worked
		return out, nil
	}

	// Fallback to previous decryption method to read old formats
	return p.deprecatedDecrypt(in)
}

func (p *DatabasePassword) deprecatedDecrypt(in []byte) ([]byte, error) {
	key := fmt.Sprintf("%x", sha256.Sum256([]byte(p.password)))[17 : 17+32]
	iv := in[0:aes.BlockSize]

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return []byte{}, err
	}

	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypted := make([]byte, len(in)-aes.BlockSize)
	encrypted := in[aes.BlockSize:]
	decrypter.XORKeyStream(decrypted, encrypted)

	return decrypted, nil
}
