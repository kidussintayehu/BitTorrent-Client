package swarm

import (
	"time"
	"github.com/kidussintayehu/BitTorrent-Client/worker"
	"github.com/kidussintayehu/BitTorrent-Client/message"
)


type progressTracker struct {
	index      int
	worker     *worker.Worker
	buf        []byte
	downloaded int
	requested  int
	backlog    int 
}

func (state *progressTracker) readMessage() error {
	msg, err := state.worker.Read()
	if err != nil {
		return err
	}
	if msg == nil { // keep-alive
		return nil
	}

	switch msg.ID {
	case message.MsgUnchoke:
		state.worker.Choked = false
	case message.MsgChoke:
		state.worker.Choked = true
	case message.MsgHave:
		index, err := message.ParseHave(msg)
		if err != nil {
			return err
		}
		state.worker.Bitfield.SetPiece(index)
	case message.MsgPiece:
		n := message.ParsePiece(state.index, state.buf, msg)
		state.downloaded += n
		state.backlog--
	}
	return nil
}

const MaxBlockSize = 16384
const MaxBacklog = 5

func attemptDownload(w *worker.Worker, piece *pieceOfWork) ([]byte) {
	state := progressTracker{
		index:  piece.index,
		worker: w,
		buf:    make([]byte, piece.length),
	}

	w.Conn.SetDeadline(time.Now().Add(30 * time.Second))
	defer w.Conn.SetDeadline(time.Time{})

	for state.downloaded < piece.length {
		if !state.worker.Choked {
			for state.backlog < MaxBacklog && state.requested < piece.length {
				blockSize := MaxBlockSize
				if piece.length-state.requested < blockSize {
					blockSize = piece.length - state.requested
				}

				err := w.SendRequest(piece.index, state.requested, blockSize)
				if err != nil {
					return nil
				}
				state.backlog++
				state.requested += blockSize
			}
		}

		err := state.readMessage()
		if err != nil {
			return nil
		}
	}

	return state.buf
}
