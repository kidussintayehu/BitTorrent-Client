package main

import (
	"fmt"
	"os"
	"time"

	"github.com/anacrolix/torrent"
)

func seed() {
	// Open the torrent file
	torrentFile, err := os.Open("example.torrent")
	if err != nil {
		panic(err)
	}
	defer torrentFile.Close()

	// Create a new Torrent client
	clientConfig := torrent.NewDefaultClientConfig()
	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Add the torrent to the client
	t, err := client.AddTorrent(torrentFile)
	if err != nil {
		panic(err)
	}

	// Wait for the torrent to be ready
	<-t.GotInfo()
	fmt.Printf("Torrent ready: %s\n", t.Info().Name)

	// Seed the file as a new peer
	seed := t.Seed()
	if seed == nil {
		panic("Failed to start seeding")
	}

	// Wait for the seed to finish
	<-seed.GotInfo()
	fmt.Printf("Seeding finished: %s\n", t.Info().Name)

	// Sleep for a few seconds to keep the program running
	time.Sleep(5 * time.Second)
}
