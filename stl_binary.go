package gostl

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/duoyume/gobox"
	"io"
)

type StlBinary struct {
	*Stl
}

func NewStlBinary(fd io.ReadCloser, rdr *bufio.Reader, ct *gobox.Converter) *StlBinary {
	return &StlBinary{
		&Stl{
			fd:  fd,
			rdr: rdr,
			box: gobox.NewBox(ct),
		},
	}
}

func (stl *StlBinary) ReadVertex(rdr *bytes.Reader) (*gobox.Vertex3, error) {
	var err error
	var v float32
	vertex := &gobox.Vertex3{}

	err = binary.Read(rdr, binary.LittleEndian, &v)
	if err != nil {
		return nil, err
	}
	vertex.X = float64(v)

	err = binary.Read(rdr, binary.LittleEndian, &v)
	if err != nil {
		return nil, err
	}
	vertex.Y = float64(v)

	err = binary.Read(rdr, binary.LittleEndian, &v)
	if err != nil {
		return nil, err
	}
	vertex.Z = float64(v)

	return vertex, nil
}

func (stl *StlBinary) ReadTriangle() (*gobox.Triangle, error) {
	n, p, err := stl.Read(BINARY_CHUNK_SIZE)
	if err != nil {
		return nil, err
	}
	if n != BINARY_CHUNK_SIZE {
		return nil, ErrUnexpectedCount
	}

	rdr := bytes.NewReader(p)
	// 跳过起始12个字节
	rdr.Seek(12, io.SeekStart)

	v1, err := stl.ReadVertex(rdr)
	if err != nil {
		return nil, err
	}

	v2, err := stl.ReadVertex(rdr)
	if err != nil {
		return nil, err
	}

	v3, err := stl.ReadVertex(rdr)
	if err != nil {
		return nil, err
	}

	// 跳过结尾2个字节
	rdr.Seek(2, io.SeekCurrent)

	return &gobox.Triangle{
		V1: v1,
		V2: v2,
		V3: v3,
	}, nil
}

func (stl *StlBinary) ReadAll(cb Callback) error {
	defer stl.fd.Close()

	var n int
	var p []byte
	var err error

	// 跳过起始80个字节
	n, p, err = stl.Read(BINARY_SKIP_HEADER_SIZE)
	if err != nil {
		return err
	}
	stl.box.SetModelName(stl.CleanModelName(p))
	//fmt.Println(string(stl.box.GetModelName()))

	if n != BINARY_SKIP_HEADER_SIZE {
		return ErrUnexpectedCount
	}
	// end

	// 读取三角形个数
	n, p, err = stl.Read(BINARY_TRIANGLES_COUNT_SIZE)
	if err != nil {
		return err
	}
	if n != BINARY_TRIANGLES_COUNT_SIZE {
		return ErrUnexpectedCount
	}
	var triangle_count int32
	err = binary.Read(bytes.NewReader(p), binary.LittleEndian, &triangle_count)
	if err != nil {
		return err
	}
	stl.box.SetTriangleCount(triangle_count)
	// end

	for {
		triangle, err := stl.ReadTriangle()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		// 表面积
		stl.box.SeekArea(triangle)
		// 体积
		stl.box.SeekVolume(triangle)
		// 长宽高
		stl.box.SeekBoundsBox(triangle)

		if cb != nil {
			cb(triangle, stl)
		}
	}

	return nil
}

func (stl *StlBinary) GetBox() *gobox.Box {
	return stl.box
}
