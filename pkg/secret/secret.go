package secret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func fileExists(fname string) bool {
	info, err := os.Stat(fname)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

type Vault struct {
	key  []byte
	path string
}

func FileVault(secret, path string) Vault {
	key := createHash(secret)
	v := Vault{
		key:  key,
		path: path,
	}
	if !fileExists(path) {
		kvmap := make(map[string]string)
		data, _ := json.Marshal(kvmap)
		if err := v.writeEncryptedFile(data); err != nil {
			panic(err)
		}
	}
	return v
}

func (v Vault) Set(k, val string) error {
	plaintext, err := v.readEncryptedFile()
	if err != nil {
		return err
	}
	var kvmap map[string]string
	err = json.Unmarshal(plaintext, &kvmap)
	kvmap[k] = val

	plaintext, err = json.Marshal(kvmap)
	if err != nil {
		return err
	}
	return v.writeEncryptedFile(plaintext)

}

func (v Vault) Get(k string) (string, error) {
	plaintext, err := v.readEncryptedFile()
	if err != nil {
		return "", err
	}
	var kvmap map[string]string
	err = json.Unmarshal(plaintext, &kvmap)
	if val, ok := kvmap[k]; ok {
		return val, nil
	}
	return "", fmt.Errorf("Key not found: %s", k)
}

func (v Vault) readEncryptedFile() (plaintext []byte, err error) {
	data, err := ioutil.ReadFile(v.path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	plaintext = decrypt(data, v.key)
	return plaintext, err
}

func (v Vault) writeEncryptedFile(data []byte) error {
	ciphertext := encrypt(data, v.key)
	f, err := os.OpenFile(v.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(ciphertext)
	return err
}
