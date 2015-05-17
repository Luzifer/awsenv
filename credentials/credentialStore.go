package credentials // import "github.com/Luzifer/awsenv/credentials"

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/Luzifer/awsenv/security"
	"gopkg.in/yaml.v2"
)

type AWSCredentialStore struct {
	Credentials      map[string]AWSCredential
	databasePassword *security.DatabasePassword `yaml:"-"`
	storageFile      string                     `yaml:"-"`
}

type AWSCredential struct {
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
}

func (a *AWSCredentialStore) SaveToFile() error {
	t, err := yaml.Marshal(a)
	if err != nil {
		return err
	}

	enc, err := a.databasePassword.Encrypt(t)
	if err != nil {
		return err
	}

	os.MkdirAll(path.Dir(a.storageFile), 0755)
	err = ioutil.WriteFile(a.storageFile, enc, 0600)
	return err
}

func FromFile(filename string, pass *security.DatabasePassword) (*AWSCredentialStore, error) {
	enc, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	dec, err := pass.Decrypt(enc)
	if err != nil {
		return nil, err
	}

	t := &AWSCredentialStore{
		databasePassword: pass,
		storageFile:      filename,
	}
	err = yaml.Unmarshal(dec, t)
	return t, err
}

func New(storefile string, pass *security.DatabasePassword) *AWSCredentialStore {
	return &AWSCredentialStore{
		databasePassword: pass,
		storageFile:      storefile,
		Credentials:      make(map[string]AWSCredential),
	}
}
