package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/lacasian/ethwheels/bestblock"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("starting best block tracker")
	httpRPC := flag.String("http-rpc-url", "", "HTTP RPC URL")
	wsRPC := flag.String("ws-rpc-url", "", "Websockets RPC URL")
	flag.Parse()

	if *httpRPC == "" && *wsRPC == "" {
		log.Fatal("please specify a \"(http|ws)-rpc-url\"")
	}

	tracker, err := bestblock.NewTracker(bestblock.Config{
		HTTP:         *httpRPC,
		WS:           *wsRPC,
		PollInterval: time.Second * 12,
	})
	if err != nil {
		log.Fatal(err)
	}

	go tracker.Run()

	subChan := tracker.Subscribe()
	errChan := tracker.Err()
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)

loop:
	for {
		select {
		case bb := <-subChan:
			log.Info(bb)
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case <-exitChan:
			break loop
		}
	}

	tracker.Close()
}
