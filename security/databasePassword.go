package security // import "github.com/Luzifer/awsenv/security"

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/rand"
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

func LoadDatabasePasswordFromFile(filename string) (*DatabasePassword, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	t := &DatabasePassword{}
	err = t.decryptPassword(buf)
	return t, err
}

func (p *DatabasePassword) SaveToFile(filename string) error {
	buf, err := p.encryptPassword()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, buf, 0600)
	return err
}

func (p *DatabasePassword) encryptPassword() ([]byte, error) {
	return []byte(p.password), nil
}

func (p *DatabasePassword) decryptPassword(encrypted []byte) error {
	p.password = string(encrypted)
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
