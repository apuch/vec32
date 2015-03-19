package vec32

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Mesh struct {
	Points []Vec3
}

const (
	typeIgnore = iota
	typeChar   = iota
	typeUchar  = iota
	typeShort  = iota
	typeUshort = iota
	typeInt    = iota
	typeUint   = iota
	typeFloat  = iota
	typeDouble = iota
)

var typeMap = map[string]uint{
	"char":   typeChar,
	"uchar":  typeUchar,
	"short":  typeUshort,
	"int":    typeInt,
	"uint":   typeUint,
	"float":  typeFloat,
	"double": typeDouble,
}

type meshBuilder struct {
	rd         *bufio.Reader
	mesh       *Mesh
	facesPhase bool
	nVerts     int
	nFaces     int
	vertFmt    string
	faceFmt    string
	vertProp   []property
	faceProp   []property
}

type property struct {
	dataType uint
	Name     string
}

func ReadPLY(r io.Reader) (m *Mesh, err error) {
	rd := bufio.NewReader(r)
	mb := meshBuilder{}
	mb.rd = rd
	mb.mesh = &Mesh{}
	mb.vertFmt = ""
	mb.faceFmt = ""
	mb.vertProp = make([]property, 64)
	mb.faceProp = make([]property, 64)
	err = mb.readHeader()
	if err != nil {
		return nil, err
	}
	return mb.mesh, nil
}

func (mb *meshBuilder) readHeader() error {
	line, err := mb.rd.ReadString('\r')
	if err != nil {
		return newErrorMesh(err.Error())
	}
	if line != "ply\r" {
		return newErrorMesh("expected file-magic ply\\r")
	}
	line, err = mb.rd.ReadString('\r')
	if err != nil {
		return newErrorMesh(err.Error())
	}
	if line != "format ascii 1.0\r" {
		return newErrorMesh("only format ascii 1.0 supported")
	}
	for {
		line, err = mb.rd.ReadString('\r')
		if err != nil {
			return newErrorMesh(err.Error())
		}
		if line == "end_header\r" {
			return nil
		} else if strings.HasPrefix(line, "comment ") {
			// pass
		} else if strings.HasPrefix(line, "element vertex ") {
			if e := mb.readElement(line); e != nil {
				return e
			}
		} else if strings.HasPrefix(line, "element face ") {
			if e := mb.readFace(line); e != nil {
				return e
			}
		} else if strings.HasPrefix(line, "property ") {
			if e := mb.readProperty(line); e != nil {
				return e
			}
		} else {
			return newErrorMesh("unexpected line in header: " + strings.TrimRight(line, "\r"))
		}
	}
}

func (mb *meshBuilder) readFace(line string) error {
	if n, e := fmt.Sscanf(line, "element face %d\r", &mb.nFaces); e != nil || n != 1 {
		return newErrorMesh("failed to parse number of faces")
	}
	return nil
}

func (mb *meshBuilder) readElement(line string) error {
	if n, e := fmt.Sscanf(line, "element vertex %d\r", &mb.nVerts); e != nil || n != 1 {
		return newErrorMesh("failed to parse number of vertices")
	}
	return nil
}

func (mb *meshBuilder) readProperty(line string) error {
	var propType, name string
	if n, e := fmt.Sscanf(line, "property %s %s", &propType, &name); n != 2 || e != nil {
		return newErrorMesh("failed to parse: " + line)
	}
	_, ok := typeMap[propType]
	if !ok {
		return newErrorMesh("unknown property type: " + propType)
	}
	return nil
}

func (mb *meshBuilder) addVertProperty(propType uint, name string) {

}

// Error in mesh creation or something
type ErrorMesh struct {
	what string
}

func newErrorMesh(what string) error {
	return &ErrorMesh{what}
}

func (e *ErrorMesh) Error() string {
	return e.what
}

func Strcmp(a, b string) int {
	var min = len(b)
	if len(a) < len(b) {
		min = len(a)
	}
	var diff int
	for i := 0; i < min && diff == 0; i++ {
		diff = int(a[i]) - int(b[i])
	}
	if diff == 0 {
		diff = len(a) - len(b)
	}
	return diff
}
