package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var numPublishers = flag.Int("n", 1, "Number of publishers")

	flag.Parse()

	args := flag.Args()

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Publisher")}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	startPublisher := func(i int) {
		nc, err := nats.Connect(*urls, opts...)
		if err != nil {
			log.Printf("Error connecting: %v", err)
		}
		defer nc.Close()

		for {
			subj, msg := args[0], []byte(args[1])

			// log.Println("Connection status:", i, nc.IsConnected(),
			// nc.ConnectedClusterName(), nc.ConnectedServerName(), nc.ConnectedAddr())

			resp, err := nc.Request(subj, []byte(fmt.Sprintf("msg:%s from publisher %d", msg, i)), 1000*time.Millisecond)
			if err != nil {
				log.Printf("Error: %v\n", err)
			} else {
				log.Printf("Received response: %s\n", resp.Data)
			}
			nc.Flush()

			if err := nc.LastError(); err != nil {
				log.Printf("Error: %v\n", err)
			}
			sleep := rand.IntN(1000)
			// log.Printf("sleeping for %d milliseconds", sleep)
			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < *numPublishers; i++ {
		log.Printf("Starting publisher %d", i)
		wg.Add(1)
		go startPublisher(i)
	}
	wg.Wait()
}
