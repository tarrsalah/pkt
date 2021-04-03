package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/tarrsalah/pkt"
	"github.com/tarrsalah/pkt/store/bolt"
)

var (
	helpFlag = flag.Bool("h", false, "show usage")
)

func usage() {
	fmt.Fprintf(os.Stderr, `pkt is tool for managing pocket items.

Usage:
	pkt <command>
The commands are:
	auth	authenticate via getpocket API
	show    sync and show getpocket dashboard

The default command is %s
`, "`show`")

	os.Exit(2)
}

func main() {
	log.SetPrefix("pkt: ")
	log.SetFlags(0)

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 && *helpFlag {
		usage()
		return
	}

	if len(args) == 0 {
		show()
		return
	}

	if args[0] == "show" {
		show()
		return
	}

	if args[0] == "auth" {
		auth()
		return
	}

	fmt.Fprintf(os.Stderr, `pkt %s: unknown ocmmand
run %s for usage.
`, args[0], "`pkt -h`")
}

func show() {
	boltPath := getBoltPath()
	configPath := getConfigPath()

	db := bolt.NewDB(boltPath)
	defer db.Close()

	auth := loadAuth(configPath)
	client := pkt.NewClient(auth)

	oldItems := db.Get()
	after := ""
	if len(oldItems) > 0 {
		after = oldItems[0].AddedAt
	}

	newItems, err := client.RetrieveAll(after)
	if err != nil {
		log.Fatal(err)
	}

	db.Put(newItems)

	items := db.Get()
	draw(items)
}

func auth() {
	configPath := getConfigPath()
	r := bufio.NewReader(os.Stdin)
	fmt.Print("pkt: enter your consumer key: ")
	key, _ := r.ReadString('\n')

	client := pkt.NewClient(nil)
	auth := client.Authenticate(strings.TrimSpace(key))
	saveAuth(auth, configPath)
	log.Println("authorized!")

}

func loadAuth(configPath string) *pkt.Auth {
	auth := &pkt.Auth{}
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(configFile, auth)
	if err != nil {
		log.Fatal(err)
	}

	return auth
}

func saveAuth(auth *pkt.Auth, configPath string) {
	configFile, err := json.MarshalIndent(auth, " ", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(configPath, configFile, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
