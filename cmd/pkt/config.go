package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func getBoltPath() string {
	return filepath.Join(getConfigDir(), "pkt.bolt")
}

func getConfigPath() string {
	return filepath.Join(getConfigDir(), "config.json")
}

func getConfigDir() string {
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
