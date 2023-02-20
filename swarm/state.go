package swarm

import (

	"github.com/kidussintayehu/BitTorrent-Client/message"
	"github.com/kidussintayehu/BitTorrent-Client/worker"
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
	msg, err := state.worker.Read() // blocking call
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

