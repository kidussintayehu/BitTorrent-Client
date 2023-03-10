package torrent

import (
	"crypto/rand"
	"log"
	"os"

	"github.com/kidussintayehu/BitTorrent-Client/utilities"
	"github.com/kidussintayehu/BitTorrent-Client/swarm"
	"github.com/kidussintayehu/BitTorrent-Client/bitfield_torrent"
)

type Torrent struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
	File        *os.File
	Bitfield    bitfield.Bitfield
}
// TorrentFile holds the metadata from a .torrent file, parsed from bencode
type TorrentFile struct {
	TrackerBaseURL string
	InfoHash       [20]byte   
	PieceHashes    [][20]byte 
	PieceLength    int       
	Length         int        
	Name           string
}

func Deserialize(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("Opening torrent file failed")
		return TorrentFile{}, err
	}
	defer file.Close()

	torrentMeta, err := bencodeUtils.ParseTorrent(file)
	if err != nil {
		log.Fatalln("Parsing torrent file content failed")
		return TorrentFile{}, err
	}

	// send the hash of Info to tracker to identify the file we want to download
	infoHash, err := torrentMeta.Info.Hash()
	if err != nil {
		log.Fatalln("Extracting torrent hash failed")
		return TorrentFile{}, err
	}

	// get hashes of each piece of the file for integrity check
	pieceHashes, err := torrentMeta.Info.SplitPieceHashes()
	if err != nil {
		log.Fatalln("Extracting hashes of pieces failed")
		return TorrentFile{}, err
	}

	
	// store in flatter struct for ease of use
	t := TorrentFile{
		TrackerBaseURL: torrentMeta.Announce,
		InfoHash:       infoHash,
		PieceHashes:    pieceHashes,
		PieceLength:    torrentMeta.Info.PieceLength,
		Length:         torrentMeta.Info.Length,
		Name:           torrentMeta.Info.Name,
	}
	return t, nil
}

func (t *TorrentFile) DownloadToFile(path string) error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:]) // use a random ID to identify ourselves to tracker
	if err != nil {
		return err
	}

	peers, err := t.requestForPeers(peerID, Port)
	if err != nil {
		return err
	}
	log.Printf("Got %d peers\n", len(peers))

	torrent := swarm.DownloadMeta{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceSize:   t.PieceLength,
		FileSize:    t.Length,
		Name:        t.Name,
	}
	buf, err := torrent.Download()
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = outFile.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

