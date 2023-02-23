package swarm

import (
	"log"
	"bytes"
	"crypto/sha1"
	"github.com/kidussintayehu/BitTorrent-Client/worker"
	"github.com/kidussintayehu/BitTorrent-Client/peers"
)



type pieceOfWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceOfResult struct {
	index int
	buf   []byte
}


type DownloadMeta struct {
	Peers       []peers.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceSize   int
	FileSize    int
	Name        string
}

func (meta *DownloadMeta) startDownloadWorker(peer peers.Peer, workQueue chan *pieceOfWork, results chan *pieceOfResult) {
	w, err := worker.New(peer, meta.PeerID, meta.InfoHash)
	if err != nil {
		log.Printf("Handshake with peer %s failed\n", peer.IP)
		return
	}
	log.Printf("Handshake with peer %s successful\n", peer.IP)

	defer w.Conn.Close()

	w.SendUnchoke()
	w.SendInterested()

	for piece := range workQueue {
		if !w.Bitfield.HasPiece(piece.index) {
			workQueue <- piece
			continue
		}

		buf := attemptDownload(w, piece)


		hash := sha1.Sum(buf)
 		if !bytes.Equal(hash[:], piece.hash[:]) {
			log.Printf("Piece #%d from %s failed integrity check, will retry\n", piece.index, peer.IP)
			workQueue <- piece
			continue
		}

		w.SendHave(piece.index)
		results <- &pieceOfResult{piece.index, buf}
	}
}

func (meta *DownloadMeta) Download() ([]byte, error) {
	log.Println("Downloading", meta.Name)

	workQueue := make(chan *pieceOfWork, len(meta.PieceHashes))
	results := make(chan *pieceOfResult)

	for index, hash := range meta.PieceHashes {
		begin := index * meta.PieceSize
		end := begin + meta.PieceSize
		if end > meta.FileSize {
			end = meta.FileSize
		}
		workQueue <- &pieceOfWork{index, hash, end-begin}
	}

	for _, peer := range meta.Peers {
		go meta.startDownloadWorker(peer, workQueue, results)
	}

	resultBuf := make([]byte, meta.FileSize)
	donePieces := 0
	for donePieces < len(meta.PieceHashes) {
		piece := <-results
		begin := piece.index * meta.PieceSize
		end := begin + meta.PieceSize
		if end > meta.FileSize {
			end = meta.FileSize
		}
		copy(resultBuf[begin:end], piece.buf)

		donePieces++
	}
	close(workQueue)
	return resultBuf, nil
}

