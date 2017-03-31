package credentials // import "github.com/Luzifer/awsenv/credentials"

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/Luzifer/awsenv/security"
	"gopkg.in/yaml.v2"
)

// AWSCredentialStore represents a storage for all the credentials
type AWSCredentialStore struct {
	Credentials      map[string]AWSCredential
	databasePassword *security.DatabasePassword `yaml:"-"`
	storageFile      string                     `yaml:"-"`
}

// AWSCredential holds the credential set for an environment
type AWSCredential struct {
    AWSProfile         string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
}

// SaveToFile stores the encrypted version of the AWSCredentialStore to the file
// the store has been loaded from
func (a *AWSCredentialStore) SaveToFile() error {
	t, err := yaml.Marshal(a)
	if err != nil {
		return err
	}

	enc, err := a.databasePassword.Encrypt(t)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Dir(a.storageFile), 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(a.storageFile, enc, 0600)
	return err
}

// UpdatePassword changes the password of the store and saves the store encrypted
// with the new password back to its file
func (a *AWSCredentialStore) UpdatePassword(passwd string) error {
	a.databasePassword = security.LoadDatabasePasswordFromInput(passwd)
	return a.SaveToFile()
}

// FromFile loads an AWSCredentialStore from the given file and decrypts it
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

// New creates an empty credential store and sets the storage location
func New(storefile string, pass *security.DatabasePassword) *AWSCredentialStore {
	return &AWSCredentialStore{
		databasePassword: pass,
		storageFile:      storefile,
		Credentials:      make(map[string]AWSCredential),
	}
}
