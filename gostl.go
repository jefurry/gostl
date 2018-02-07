package gostl

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/jefurry/gobox"
	"io"
	"regexp"
	"syscall"
)

const (
	READ_LINE_SIZE = 128

	MAX_READ_LINE_SIZE = 1024

	// 打开stl文件权限
	OPEN_STL_FILE_PERM = 0666

	// 顶点个数
	TRIANGLE_VERTEX_COUNT = 3

	// 跳过头部字节数
	BINARY_SKIP_HEADER_SIZE = 80

	// 三角形个字数占用字节数
	BINARY_TRIANGLES_COUNT_SIZE = 4

	// 块大小
	BINARY_CHUNK_SIZE = 50
)

var (
	ErrUnexpectedFormat = errors.New("unexpected stl format.")
	ErrStlParser        = errors.New("stl parse error.")
	ErrUnexpectedCount  = errors.New("unexpected count.")
	ErrInvalidLine      = errors.New("invalid stl line")
)

type Callback func(*gobox.Triangle, StlParser)

type StlParser interface {
	ReadAll(Callback) error
	GetBox() *gobox.Box
}

type Stl struct {
	fd  io.ReadCloser
	rdr *bufio.Reader
	box *gobox.Box
}

func (stl *Stl) Read(n int) (int, []byte, error) {
	if n <= 0 {
		return 0, nil, nil
	}

	buf := bytes.NewBuffer(nil)
	buf.Grow(n)

	var nn int = 0
	for {
		p := make([]byte, n-nn)
		_n, err := stl.rdr.Read(p)

		if _n > 0 {
			nn += _n
			buf.Write(p[0:_n])
		}

		if err != nil {
			if err != syscall.EINTR && err != syscall.EAGAIN && err != syscall.EWOULDBLOCK {
				return nn, buf.Bytes(), err
			}
		}

		if nn >= n {
			break
		}
	}

	return nn, buf.Bytes(), nil
}

// check 是否验证是否有效的包含顶点数据的行
func (stl *Stl) ReadLine(check bool) ([]byte, error) {
	line, err := stl.rdr.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	line = bytes.TrimSpace(line)
	if bytes.Equal(line, []byte{}) || bytes.Index(line, []byte("endsolid")) == 0 {
		return nil, io.EOF
	}

	if check {
		if bytes.Index(line, []byte("vertex")) == 0 {
			return line, nil
		}

		return nil, ErrInvalidLine
	}

	return line, nil
}

func (stl *Stl) CleanModelName(modelName []byte) []byte {
	randBytes := RandBytes(20)
	var p []byte
	p = bytes.TrimSpace(modelName)
	p = bytes.Replace(p, []byte("solid "), []byte(""), -1)

	reg, err := regexp.Compile("[^_a-zA-Z0-9]+")
	if err != nil {
		return randBytes
	}

	p = reg.ReplaceAll(p, []byte("_"))
	if bytes.Equal(p, []byte{}) {
		return randBytes
	}

	return p
}
