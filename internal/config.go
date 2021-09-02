package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// GetAuth returns the saved pocket credentials
func GetAuth() *Creds {
	auth := &Creds{}
	configFile, err := ioutil.ReadFile(configPath())
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
func PutAuth(auth *Creds) {
	content, err := json.MarshalIndent(auth, " ", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(configPath(), content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// configDir returns the path of the confi directory
func configDir() string {
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

// configPath returns the configPath of the config file
func configPath() string {
	return filepath.Join(configDir(), "config.json")
}
