package main

import (
	"fmt"
	"os"
	"time"

	"github.com/anacrolix/torrent"
)

func seedafterdownload() {
	// Open the existing torrent file
	torrentFile, err := os.Open("testfiles/debian-10.2.0-amd64-netinst.iso.torrent")
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

	// Add the existing torrent to the client
	t, err := client.AddTorrent(torrentFile)
	if err != nil {
		panic(err)
	}

	// Wait for the torrent to be ready
	<-t.GotInfo()
	fmt.Printf("Torrent ready: %s\n", t.Info().Name)

	// Download the file and wait for it to complete
	file := t.Files()[0]
	err = file.Priority(torrent.PiecePriorityStandard)
	if err != nil {
		panic(err)
	}
	err = t.DownloadPieces(file.PieceLength() * (file.NumPieces() - 1))
	if err != nil {
		panic(err)
	}

	// Seed the downloaded file as a new peer
	seed, err := t.SeedFile(file)
	if err != nil {
		panic(err)
	}

	// Wait for the seed to finish
	<-seed.GotInfo()
	fmt.Printf("Seeding finished: %s\n", file.Path())

	// Sleep for a few seconds to keep the program running
	time.Sleep(5 * time.Second)
}
