package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tarrsalah/pkt"
	"github.com/tarrsalah/pkt/internal/bolt"
	"github.com/tarrsalah/pkt/internal/config"
	"github.com/tarrsalah/pkt/internal/ui"
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
	db := bolt.NewDB()
	defer db.Close()

	auth := config.GetAuth()
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
	app := ui.NewWindow(items)
	app.Run()
}

func auth() {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("pkt: enter your consumer key: ")
	key, _ := r.ReadString('\n')

	client := pkt.NewClient(nil)
	auth := client.Authenticate(strings.TrimSpace(key))
	config.PutAuth(auth)
	log.Println("authorized!")

}
