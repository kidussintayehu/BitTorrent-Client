package worker

import "github.com/kidussintayehu/BitTorrent-Client/message"


func (w *Worker) SendInterested() error {
	msg := message.Message{ID: message.MsgInterested}
	_, err := w.Conn.Write(msg.Serialize())
	return err
}

func (w *Worker) SendNotInterested() error {
	msg := message.Message{ID: message.MsgNotInterested}
	_, err := w.Conn.Write(msg.Serialize())
	return err
}

func (w *Worker) SendUnchoke() error {
	msg := message.Message{ID: message.MsgUnchoke}
	_, err := w.Conn.Write(msg.Serialize())
	return err
}

func (w *Worker) Read() (*message.Message, error) {
	msg, err := message.Read(w.Conn)
	return msg, err
}

func (w *Worker) SendRequest(index, begin, length int) error {
	req := message.FormatRequest(index, begin, length)
	_, err := w.Conn.Write(req.Serialize())
	return err
}

func (w *Worker) SendHave(index int) error {
	msg := message.FormatHave(index)
	_, err := w.Conn.Write(msg.Serialize())
	return err
}
