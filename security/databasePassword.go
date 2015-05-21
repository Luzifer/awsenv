package security // import "github.com/Luzifer/awsenv/security"

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

type DatabasePassword struct {
	password string
}

func LoadDatabasePasswordFromInput(input string) *DatabasePassword {
	return &DatabasePassword{
		password: input,
	}
}

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
	defer res.Body.Close()

	pb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &DatabasePassword{
		password: strings.TrimSpace(string(pb)),
	}, nil
}

func (p *DatabasePassword) SpawnLockAgent(filename string) error {
	var err error
	proc := exec.Command(os.Args[0], "lockagent")
	proc.Env = []string{
		fmt.Sprintf("DBFILE=%s", filename),
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
	defer res.Body.Close()

	return nil
}

func (p *DatabasePassword) Encrypt(in []byte) ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	key := fmt.Sprintf("%x", sha256.Sum256([]byte(p.password)))[17 : 17+32]
	ivInt := rand.Intn(63 - aes.BlockSize)
	iv := fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewV4().String())))[ivInt : ivInt+aes.BlockSize]

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return []byte{}, err
	}

	encrypter := cipher.NewCFBEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(in))
	encrypter.XORKeyStream(encrypted, in)

	out := append([]byte(iv), encrypted...)

	return out, nil
}

func (p *DatabasePassword) Decrypt(in []byte) ([]byte, error) {
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
