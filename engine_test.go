package gostl

import (
	"fmt"
	"github.com/jefurry/gobox"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, err := NewFileEngine("./examples/binary.stl", gobox.DefaultConverter)
	//engine, err := NewFileEngine("./examples/ascii.stl", gobox.DefaultConverter)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer engine.fd.Close()
	//fmt.Println(engine.IsBinary())
	//fmt.Println(engine.IsAscii())

	parser, err := engine.GetParser()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = parser.ReadAll(func(triangle *gobox.Triangle, stl StlParser) {
		//fmt.Println(triangle, stl)
		//fmt.Println(stl.GetBox().GetBoundsBox(gobox.UNIT_CM, 2))
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(parser.GetBox().GetArea(gobox.UNIT_CM, 2))
	fmt.Println(parser.GetBox().GetVolume(gobox.UNIT_CM, 2))
	fmt.Println(parser.GetBox().GetBoundsBox(gobox.UNIT_CM, 2))
	fmt.Println(parser.GetBox().GetTriangleCount())
}
