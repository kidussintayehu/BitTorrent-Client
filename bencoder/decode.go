package bencode

import (
	"bufio"
	"io"
)


func Decode(reader io.Reader) (data interface{}, err error) {
	
	bufioReader, ok := reader.(*bufio.Reader)
	if !ok {
		bufioReader = newBufioReader(reader)
		defer bufioReaderPool.Put(bufioReader)
	}

	return decodeFromReader(bufioReader)
}
