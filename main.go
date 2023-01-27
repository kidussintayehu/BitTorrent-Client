package main

import (
	torrentfile "BitTorrent-Client/torrent_file"
	"log"
	"os"
)

func main() {
	inPath := os.Args[1]

	_, err := torrentfile.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}

}
