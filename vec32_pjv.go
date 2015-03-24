package vec32

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Struct holding all informations of a mesh
//
// this struct will move somewhere else as soon as we support a second
// formar. Until then...
type Mesh struct {
	// The vertices we have
	Verts []Vec3
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

const (
	propIgnore = iota
	propX      = iota
	propY      = iota
	propZ      = iota
	propCount  = iota
)

var propMap = map[string]uint{
	"x": propX,
	"y": propY,
	"z": propZ,
}

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
	rd          *bufio.Reader
	mesh        *Mesh
	facesPhase  bool
	nVerts      int
	nFaces      int
	vertProp    []property
	vertPropIdx int
	faceProp    []property
	facePropIdx int
	haveVerts   bool
	scanner     *bufio.Scanner
}

type property struct {
	propType uint
	propIdx  uint
	name     string
}

func ReadPLY(r io.Reader) (m *Mesh, err error) {
	rd := bufio.NewReader(r)
	mb := meshBuilder{}
	mb.rd = rd
	mb.mesh = &Mesh{}
	mb.vertProp = make([]property, 64)
	mb.faceProp = make([]property, 64)
	if err = mb.readHeader(); err != nil {
		return nil, err
	}
	mb.scanner = bufio.NewScanner(mb.rd)
	mb.scanner.Split(bufio.ScanWords)
	if err = mb.readVerts(); err != nil {
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
			return mb.validateHeader()
		} else if strings.HasPrefix(line, "comment ") {
			// pass
		} else if strings.HasPrefix(line, "element vertex ") {
			if e := mb.readElementVertex(line); e != nil {
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

func (mb *meshBuilder) readElementVertex(line string) error {
	if n, e := fmt.Sscanf(line, "element vertex %d\r", &mb.nVerts); e != nil || n != 1 || mb.nVerts < 0 {
		return newErrorMesh("failed to parse number of vertices")
	}
	mb.haveVerts = mb.nVerts > 0
	return nil
}

func (mb *meshBuilder) readProperty(line string) error {
	var propTypeName, name string
	if n, e := fmt.Sscanf(line, "property %s %s", &propTypeName, &name); n != 2 || e != nil {
		return newErrorMesh("failed to parse: " + line)
	}
	propType, ok := typeMap[propTypeName]
	if !ok {
		return newErrorMesh("unknown property type: " + propTypeName)
	}
	if !mb.facesPhase {
		mb.addVertProperty(propType, name)
	}
	return nil
}

func (mb *meshBuilder) addVertProperty(propType uint, name string) {
	prop := &mb.vertProp[mb.vertPropIdx]
	mb.vertPropIdx += 1
	prop.name = name
	prop.propType = propType
	prop.propIdx, _ = propMap[name]
}

func (mb *meshBuilder) validateHeader() error {
	if !mb.haveVerts {
		return nil
	}
	var haveX, haveY, haveZ bool
	for i := 0; i < mb.vertPropIdx; i++ {
		switch mb.vertProp[i].propIdx {
		case propX:
			haveX = true
		case propY:
			haveY = true
		case propZ:
			haveZ = true
		}
	}
	if !(haveX && haveY && haveZ) {
		return newErrorMesh("invalid vertex definition (missing coordinate)")
	}
	return nil
}

func (mb *meshBuilder) readVerts() error {
	for i := 0; i < mb.nVerts; i++ {
		mb.scanner.Split(splitSpace)
		for j := 0; j < mb.vertPropIdx; j++ {
			if !mb.scanner.Scan() {
				return mb.scanner.Err()
			}
			if _, err := strconv.ParseFloat(mb.scanner.Text(), 32); err != nil {
				return newErrorMesh("could not convert `" + mb.scanner.Text() + "` to float")
			}
		}
		mb.scanner.Split(splitBreak)
	}
	return nil
}

func splitSpace(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var start, end int
	for ; start < len(data); start++ {
		if data[start] != ' ' {
			break
		}
		if data[start] == '\r' {
			return 0, nil, newErrorMesh("unexpected end of line")
		}
	}
	if start < len(data) && data[start] == '\r' {
		return 0, nil, newErrorMesh("unexpected end of line")
	}
	if start == len(data)-1 {
		return 0, nil, newErrorMesh("unexpected end of file")
	}
	for end = start + 1; end < len(data); end++ {
		if data[end] == ' ' || data[end] == '\r' {
			return end, data[start:end], nil
		}
	}
	return 0, nil, newErrorMesh("unexpected end of file")
}

func splitBreak(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var start, end int
	for ; start < len(data); start++ {
		if data[start] != '\r' {
			break
		}
		if data[start] != ' ' {
			return 0, nil, newErrorMesh("unexpected token")
		}
	}
	start += 1
	if start == len(data)-1 {
		return 0, nil, newErrorMesh("unexpected end of file")
	}
	for end = start + 1; end < len(data); end++ {
		if data[end] == ' ' || data[end] == '\r' {
			return end, data[start:end], nil
		}
	}
	return 0, nil, newErrorMesh("unexpected end of file")
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
