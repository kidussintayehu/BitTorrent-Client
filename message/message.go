
package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MessageID uint8
type Message struct {
	ID      MessageID
	Payload []byte
}

const (
	MsgChoke MessageID = 0
	MsgUnchoke MessageID = 1
	MsgInterested MessageID = 2
	MsgNotInterested MessageID = 3
	MsgHave MessageID = 4
	MsgBitfield MessageID = 5
	MsgRequest MessageID = 6
	MsgPiece MessageID = 7
	MsgCancel MessageID = 8
)



// FormatRequest creates a REQUEST message
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: MsgRequest, Payload: payload}
}

func FormatChoke() *Message {
	return &Message{ID: MsgChoke}
}

func FormatPiece(Payload []byte) *Message {
	return &Message{ID: MsgBitfield, Payload: Payload}
}

// FormatHave creates a HAVE message
func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: MsgHave, Payload: payload}
}

func ParseUnchoke(msg *Message) (int, error) {
	if msg.ID != MsgUnchoke {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", MsgUnchoke, msg.ID)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

// ParsePiece parses a PIECE message and copies its payload into a buffer
func ParsePiece(index int, buf []byte, msg *Message) (int) {
	if msg.ID != MsgPiece {
		return 0
	}
	if len(msg.Payload) < 8 {
		return 0
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0
	}
	copy(buf[begin:], data)
	return len(data)
}

// ParseHave parses a HAVE message
func ParseHave(msg *Message) (int, error) {
	if msg.ID != MsgHave {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", MsgHave, msg.ID)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}


func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1) 
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	return buf
}

func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)
	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		ID:      MessageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return &m, nil
}

func (m *Message) name() string {
	if m == nil {
		return "KeepAlive"
	}
	switch m.ID {
	case MsgNotInterested:
		return "NotInterested"
	case MsgHave:
		return "Have"
	case MsgBitfield:
		return "Bitfield"
	case MsgRequest:
		return "Request"
	case MsgPiece:
		return "Piece"
	case MsgChoke:
		return "Choke"
	case MsgUnchoke:
		return "Unchoke"
	case MsgInterested:
		return "Interested"
	
	case MsgCancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.ID)
	}
}

func (m *Message) String() string {
	if m == nil {
		return m.name()
	}
	return fmt.Sprintf("%s [%d]", m.name(), len(m.Payload))
}
