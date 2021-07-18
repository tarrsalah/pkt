package config

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"encoding/json"
	"github.com/tarrsalah/pkt"
	"io/ioutil"
)

// GetAuth returns the saved pocket credentials
func GetAuth() *pkt.Auth {
	auth := &pkt.Auth{}
	configFile, err := ioutil.ReadFile(Path())
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(configFile, auth)
	if err != nil {
		log.Fatal(err)
	}

	return auth
}

// PutAuth saves the pocket credentials into a file
func PutAuth(auth *pkt.Auth) {
	content, err := json.MarshalIndent(auth, " ", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(Path(), content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Dir returns the path of the confi directory
func Dir() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configDir := filepath.Join(user.HomeDir, ".config", "pkt")
	err = os.MkdirAll(configDir, 0777)
	if err != nil {
		log.Fatal(err)
	}

	return configDir
}

// Path returns the Path of the config file
func Path() string {
	return filepath.Join(Dir(), "config.json")
}
