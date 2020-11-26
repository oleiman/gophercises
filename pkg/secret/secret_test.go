package secret

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewVaultNoFile(t *testing.T) {
	v := FileVault("foo", "bar")
	defer os.Remove(v.path)
	if !fileExists(v.path) {
		t.Errorf("FileVault constructor didn't create the vault file")
	}
}

func TestGetSet(t *testing.T) {
	v := FileVault("foo", "secrets")
	defer os.Remove(v.path)

	if _, err := v.Get("baz"); err == nil {
		t.Error("Unexpected value from empty vault")
	}

	if err := v.Set("baz", "qux"); err != nil {
		t.Errorf("v.Set failed: %s", err)
	}

	val, err := v.Get("baz")
	if err != nil {
		t.Errorf("v.Get failed: %s", err)
	} else if val != "qux" {
		t.Errorf(`v.Get("baz"): expected "qux", got %s`, val)
	}
}

func TestEncrypted(t *testing.T) {
	v := FileVault("foo", "secrets")
	defer os.Remove(v.path)
	v.Set("bar", "baz")
	v.Set("qux", "flux")
	data, _ := ioutil.ReadFile(v.path)
	var kvmap map[string]string
	err := json.Unmarshal(data, &kvmap)
	if err == nil {
		t.Errorf("Unexpectedly unmarshaled encrypted file!")
	}
}

func TestBadPassword(t *testing.T) {
	v := FileVault("foo", "secrets")
	defer os.Remove(v.path)
	v.Set("bar", "baz")

	v2 := FileVault("WRONG", "secrets")

	_, err := v2.Get("bar")
	if err == nil {
		t.Errorf("Unexpected v.Get success with incorrect password")
	}
}
