package gostl

import (
	"bufio"
	"bytes"
	"math/rand"
	"time"
)

func RandByte(n int) []byte {
	chars := []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u',
		'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P',
		'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '_',
	}

	b := make([]byte, n)
	for i, l := 0, len(chars); i < n; i++ {
		rand.Seed(time.Now().UnixNano())
		b[i] = chars[rand.Intn(l)]
	}

	return b
}

func peek_line(rdr *bufio.Reader) ([]byte, error) {
	var buf *bytes.Buffer = bytes.NewBuffer(nil)
	buf.Grow(MAX_READ_LINE_SIZE)

	var n int = 1
	var size int

	for {
		size = n * READ_LINE_SIZE
		if size > MAX_READ_LINE_SIZE {
			return nil, bufio.ErrBufferFull
		}

		buf.Reset()
		b, err := rdr.Peek(size)
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}

		if pos := bytes.Index(buf.Bytes(), []byte("\n")); pos != -1 {
			buf.Truncate(pos)
			return buf.Bytes(), nil
		}

		n += 1
	}
}
