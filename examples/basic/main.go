package main

import (
	"flag"
	"fmt"
	"github.com/goraft/raft"
	"github.com/prataprc/go-failsafe"
	"log"
	"math/rand"
	"os"
	"time"
)

var options struct {
	name     string
	listAddr int
	join     string
	trace    bool
	debug    bool
}

func init() {
	flag.StringVar(&options.name, "name", "basic0", "server's unique name")
	flag.StringVar(&options.listAddr, "s", "localhost:4001", "host:port listen")
	flag.StringVar(&options.join, "join", "", "host:port of leader to join")
	flag.BoolVar(&options.trace, "trace", false, "Raft trace debugging")
	flag.BoolVar(&options.debug, "debug", false, "Raft debugging")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] <data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	if options.trace {
		failsafe.SetLogLevel(raft.Trace)
	} else if options.debug {
		failsafe.SetLogLevel(raft.Debug)
	}

	failsafe.RegisterCommands() // Setup commands.

	// path
	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatal("Data path argument required")
	}
	listAddr := options.listAddr
	name, leader, path := options.name, options.join, flag.Arg(0)
	log.SetFlags(log.LstdFlags)

	killch, quitch := make(chan []interface{}), make(chan []interface{})
	failsafe.StartDemoServer(name, path, listAddr, leader, quitch, killch)
	time.Sleep(1 * time.Second)

	connAddr := fmt.Sprintf("%v:%v", options.host, options.port)
	client := failsafe.NewSafeDictClient("http://" + connAddr)

	CAS, err := client.GetCAS()
	handleError(err)
	fmt.Println("Got initial CAS", CAS)

	CAS, err = client.SetCAS("/eyeColor", "brown", CAS)
	handleError(err)
	fmt.Println("Set /eyeColor gave nextCAS as", CAS)

	value, CAS, err := client.Get("/eyeColor")
	handleError(err)
	fmt.Printf("Get /eyeColor returned %v with CAS %v\n", value, CAS)

	killch <- []interface{}{failsafe.DemoCmdQuit}
	<-quitch
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
