package gostl

import (
	"bufio"
	"github.com/duoyume/gobox"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type StlAscii struct {
	*Stl
}

func NewStlAscii(fd io.ReadCloser, rdr *bufio.Reader, ct *gobox.Converter) *StlAscii {
	return &StlAscii{
		&Stl{
			fd:  fd,
			rdr: rdr,
			box: gobox.NewBox(ct),
		},
	}
}

func (stl *StlAscii) ReadVertex(line []byte) (*gobox.Vertex3, error) {
	var base int = 32
	reg, err := regexp.Compile(`[\t\n\v\f\r ]+`)
	if err != nil {
		return nil, err
	}

	bb := reg.Split(string(line), -1)
	if len(bb) != 4 || !strings.EqualFold(bb[0], "vertex") {
		return nil, ErrInvalidLine
	}

	vertex := &gobox.Vertex3{}
	x, err := strconv.ParseFloat(bb[1], base)
	if err != nil {
		return nil, err
	}
	vertex.X = float64(x)

	y, err := strconv.ParseFloat(bb[2], base)
	if err != nil {
		return nil, err
	}
	vertex.Y = float64(y)

	z, err := strconv.ParseFloat(bb[3], base)
	if err != nil {
		return nil, err
	}
	vertex.Z = float64(z)

	return vertex, nil
}

func (stl *StlAscii) ReadTriangle(line1, line2, line3 []byte) (*gobox.Triangle, error) {
	v1, err := stl.ReadVertex(line1)
	if err != nil {
		return nil, err
	}

	v2, err := stl.ReadVertex(line2)
	if err != nil {
		return nil, err
	}

	v3, err := stl.ReadVertex(line3)
	if err != nil {
		return nil, err
	}

	return &gobox.Triangle{V1: v1, V2: v2, V3: v3}, nil
}

func (stl *StlAscii) ReadAll(cb Callback) error {
	defer stl.fd.Close()

	// 三角形个数
	var triangle_count int32 = 0

	first_line, err := stl.ReadLine(false)
	if err != nil {
		return err
	}
	stl.box.SetModelName(stl.CleanModelName(first_line))
	//fmt.Println(string(stl.box.GetModelName()))

	for {
		// 读取三行
		line1, err := stl.ReadLine(true)
		if err != nil {
			if err == io.EOF {
				break
			}

			if err == ErrInvalidLine {
				continue
			}

			return err
		}

		line2, err := stl.ReadLine(true)
		if err != nil {
			break
		}
		line3, err := stl.ReadLine(true)
		if err != nil {
			break
		}
		// end

		triangle, err := stl.ReadTriangle(line1, line2, line3)
		if err != nil {
			return err
		}

		triangle_count += 1
		// 三角形个数
		stl.box.SetTriangleCount(triangle_count)
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

func (stl *StlAscii) GetBox() *gobox.Box {
	return stl.box
}
