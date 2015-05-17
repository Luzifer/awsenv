package security

import (
	"crypto/aes"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadDatabasePasswordFromInput(t *testing.T) {
	password := "test1235"
	p := LoadDatabasePasswordFromInput(password)
	if p.password != password {
		t.Error("Password in store did not match input password")
	}
}

func TestIVVariance(t *testing.T) {
	seenIV := []string{}
	password := "test1234"
	message := "Hello, I'm a test."
	p := LoadDatabasePasswordFromInput(password)

	for i := 0; i < 500; i++ {
		enc, err := p.Encrypt([]byte(message))
		if err != nil {
			t.Errorf("It errored: %s", err)
		}
		iv := string(enc[0:aes.BlockSize])
		for _, v := range seenIV {
			if iv == v {
				t.Errorf("There is a duplicate IV: %s", iv)
			}
		}
		seenIV = append(seenIV, iv)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	password := "test1234"
	message := "Hello, I'm a test."
	p := LoadDatabasePasswordFromInput(password)

	for i := 0; i < 10; i++ {
		enc, err := p.Encrypt([]byte(message))
		if err != nil {
			t.Errorf("Encrypt errored: %s", err)
		}
		dec, err := p.Decrypt(enc)
		if err != nil {
			t.Errorf("Decrypt errored: %s", err)
		}
		if string(dec) != message {
			t.Errorf("Messages did not match: in='%s' out='%s'", message, string(dec))
		}
	}
}

func TestPredefinedDecrypt(t *testing.T) {
	enc := []byte{51, 57, 100, 48, 48, 53, 98, 53, 102, 57, 101, 99, 50, 99, 56, 101, 85, 174, 161, 148, 98, 169, 245, 9, 64, 101}
	password := "test1234"
	message := "Hallo Welt"

	p := LoadDatabasePasswordFromInput(password)
	dec, err := p.Decrypt(enc)

	if err != nil {
		t.Errorf("Decrypt errored: %s", err)
	}

	if string(dec) != message {
		t.Errorf("Messages did not match: in='%s' out='%s'", message, string(dec))
	}
}

func TestMessageIsEncrypted(t *testing.T) {
	password := "test1234"
	message := "Hallo Welt"

	p := LoadDatabasePasswordFromInput(password)
	enc, err := p.Encrypt([]byte(message))
	if err != nil {
		t.Errorf("Encrypt errored: %s", err)
	}

	msg := enc[aes.BlockSize:]
	if string(msg) == message {
		t.Error("Output was input message")
	}

	if len(string(msg)) != len(message) {
		t.Error("Messages had different lengths")
	}
}

func TestSaveAndLoad(t *testing.T) {
	password := "test1234"

	file, _ := ioutil.TempFile(os.TempDir(), "prefix")
	defer os.Remove(file.Name())

	p := LoadDatabasePasswordFromInput(password)
	p.SaveToFile(file.Name())

	p2, err := LoadDatabasePasswordFromFile(file.Name())
	if err != nil {
		t.Errorf("LoadDatabasePasswordFromFile errored: %s", err)
	}

	if p2.password != password {
		t.Error("Loaded password does not match")
	}

}
