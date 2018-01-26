package gostl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/duoyume/gobox"
	"io"
	"os"
)

type Engine struct {
	// 字节大小
	size int64
	fd   io.ReadCloser
	rdr  *bufio.Reader
	ct   *gobox.Converter
}

func NewEngine(size int64, fd io.ReadCloser, ct *gobox.Converter) *Engine {
	return &Engine{
		size: size,
		fd:   fd,
		rdr:  bufio.NewReader(fd),
		ct:   ct,
	}
}

func NewFileEngine(file string, ct *gobox.Converter) (*Engine, error) {
	fd, err := os.OpenFile(file, os.O_RDONLY, OPEN_STL_FILE_PERM)
	if err != nil {
		return nil, err
	}

	fileInfo, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	size := fileInfo.Size()

	return NewEngine(size, fd, ct), nil
}

func (e *Engine) IsBinary() bool {
	var err error

	b, err := e.rdr.Peek(BINARY_SKIP_HEADER_SIZE + BINARY_TRIANGLES_COUNT_SIZE)
	if err != nil {
		panic(err)
	}

	rdr := bytes.NewReader(b[BINARY_SKIP_HEADER_SIZE:])
	/*
		_, err = rdr.Seek(BINARY_SKIP_HEADER_SIZE, io.SeekStart)
		if err != nil {
			panic(err)
		}
	*/
	var triangle_count int32
	err = binary.Read(rdr, binary.LittleEndian, &triangle_count)
	if err != nil {
		panic(err)
	}

	// 通过字节数判断
	size := int64(triangle_count*BINARY_CHUNK_SIZE + (BINARY_SKIP_HEADER_SIZE + BINARY_TRIANGLES_COUNT_SIZE))
	if size == e.size {
		return true
	}

	return false
}

func (e *Engine) IsAscii() bool {
	if e.IsBinary() {
		return false
	}

	fline, err := peek_line(e.rdr)
	if err != nil {
		panic(err)
	}

	sline, err := peek_line(e.rdr)
	if err != nil {
		panic(err)
	}

	if bytes.Index(fline, []byte("solid")) != 0 && bytes.Index(sline, []byte("facet normal")) == -1 {
		return false
	}

	return true
}

func (e *Engine) GetParser() (parser StlParser, err error) {
	defer func() {
		if er := recover(); er != nil {
			parser = nil
			//err = ErrStlParser
			if e, ok := er.(error); ok {
				//err = er.(error)
				err = e
			} else {
				err = ErrStlParser
			}
		}
	}()

	if e.IsBinary() {
		return NewStlBinary(e.fd, e.rdr, e.ct), nil
	}

	if e.IsAscii() {
		return NewStlAscii(e.fd, e.rdr, e.ct), nil
	}

	e.fd.Close()

	return nil, ErrUnexpectedFormat
}

func (e *Engine) SetSize(size int64) {
	e.size = size
}

func (e *Engine) GetSize() int64 {
	return e.size
}
