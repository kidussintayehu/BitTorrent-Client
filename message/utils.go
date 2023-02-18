package message

import (
	"encoding/binary"
	"fmt"
)

// FormatRequest creates a REQUEST message
// It takes three integers,
// and returns a pointer to a Message struct with the ID field set to Request,
// and the Payload field set to a byte slice containing the three integers in
// big-endian order
func FormatRequest(index, begin, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{ID: Request, Payload: payload}
}

// FormatHave creates a HAVE message
// FormatHave takes an integer index and returns a pointer to a Message with the ID Have and a payload
// of the index encoded as a 4-byte big-endian integer.
func FormatHave(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{ID: Have, Payload: payload}
}

// ParsePiece parses a PIECE message and copies its payload into a buffer
// It takes a buffer, a message, and an index, and it copies the payload of the message into the buffer
// at the correct offset
func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID != Piece {
		return 0, fmt.Errorf("Expected PIECE (ID %d), got ID %d", Piece, msg.ID)
	}
	if len(msg.Payload) < 8 {
		return 0, fmt.Errorf("Payload too short. %d < 8", len(msg.Payload))
	}
	parsedIndex := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
	if parsedIndex != index {
		return 0, fmt.Errorf("Expected index %d, got %d", index, parsedIndex)
	}
	begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
	if begin >= len(buf) {
		return 0, fmt.Errorf("Begin offset too high. %d >= %d", begin, len(buf))
	}
	data := msg.Payload[8:]
	if begin+len(data) > len(buf) {
		return 0, fmt.Errorf("Data too long [%d] for offset %d with length %d", len(data), begin, len(buf))
	}
	copy(buf[begin:], data)
	return len(data), nil
}

// ParseHave parses a HAVE message
// It takes a message, checks that it's a HAVE message, checks that the payload is 4 bytes long, and
// then returns the index of the piece that the peer has
func ParseHave(msg *Message) (int, error) {
	if msg.ID != Have {
		return 0, fmt.Errorf("Expected HAVE (ID %d), got ID %d", Have, msg.ID)
	}
	if len(msg.Payload) != 4 {
		return 0, fmt.Errorf("Expected payload length 4, got length %d", len(msg.Payload))
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}

// A method of the Message struct. It returns the name of the message.
func (m *Message) name() string {
	if m == nil {
		return "KeepAlive"
	}
	switch m.ID {
	case Choke:
		return "Choke"
	case Unchoke:
		return "Unchoke"
	case Interested:
		return "Interested"
	case NotInterested:
		return "NotInterested"
	case Have:
		return "Have"
	case Bitfield:
		return "Bitfield"
	case Request:
		return "Request"
	case Piece:
		return "Piece"
	case Cancel:
		return "Cancel"
	default:
		return fmt.Sprintf("Unknown#%d", m.ID)
	}
}

// A method of the Message struct. It returns the name of the message.
func (m *Message) String() string {
	if m == nil {
		return m.name()
	}
	return fmt.Sprintf("%s [%d]", m.name(), len(m.Payload))
}
