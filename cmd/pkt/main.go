package main

import (
	"bufio"
	"flag"
	"fmt"
	pkt "github.com/tarrsalah/pkt/internal"
	"log"
	"os"
	"strings"
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

	if len(args) == 0 || args[0] == "show" {
		db := pkt.NewDB()
		defer db.Close()

		auth := pkt.GetAuth()
		client := pkt.NewClient(auth)

		oldItems, _ := db.Get()
		after := ""
		if len(oldItems) > 0 {
			after = oldItems[0].AddedAt
		}

		newItems, err := client.RetrieveAll(after)
		if err != nil {
			log.Fatal(err)
		}
		db.Put(newItems)

		list, _ := db.Get()

		ui := pkt.NewApp(list)
		ui.Run()
		return
	}

	if args[0] == "auth" {
		r := bufio.NewReader(os.Stdin)
		fmt.Print("pkt: enter your consumer key: ")
		key, _ := r.ReadString('\n')
		client := pkt.NewClient(nil)
		auth := client.Authenticate(strings.TrimSpace(key))
		pkt.PutAuth(auth)
		log.Println("authorized!")
		return
	}

	fmt.Fprintf(os.Stderr, `pkt %s: unknown command run %s for usage.
`, args[0], "`pkt -h`")
}
