package vec32

import (
	"bufio"
	"bytes"
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
	Tris  []Triangle
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

const (
	elementIgnore = iota
	elementVert   = iota
	elementFace   = iota
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
	currElement uint
	triIdx      int
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
	mb.scanner = bufio.NewScanner(mb.rd)
	mb.scanner.Split(scanLines)
	if err = mb.readHeader(); err != nil {
		return nil, err
	}
	mb.scanner.Split(bufio.ScanWords)
	Info.Printf("Read header, start to read values (%d verts, %d faces)",
		mb.nVerts, mb.nFaces)
	mb.mesh.Verts = make([]Vec3, mb.nVerts)
	mb.mesh.Tris = make([]Triangle, mb.nFaces)
	if err = mb.readVerts(); err != nil {
		return nil, err
	}
	if err = mb.readFaces(); err != nil {
		return nil, err
	}
	return mb.mesh, nil
}

func (mb *meshBuilder) readHeader() error {
	var line string
	var err error
	if line, err = mb.nextLine(); err != nil {
		return err
	}
	if line != "ply" {
		return newErrorMesh("expected file-magic ply\\r")
	}
	if line, err = mb.nextLine(); err != nil {
		return err
	}
	if line != "format ascii 1.0" {
		return newErrorMesh("only format ascii 1.0 supported")
	}
	for {
		line, err = mb.nextLine()
		if err != nil {
			return newErrorMesh(err.Error())
		}
		if line == "end_header" {
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
		} else if strings.HasPrefix(line, "element ") {
			mb.currElement = elementIgnore
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
	mb.currElement = elementFace
	return nil
}

func (mb *meshBuilder) readElementVertex(line string) error {
	if n, e := fmt.Sscanf(line, "element vertex %d\r", &mb.nVerts); e != nil || n != 1 || mb.nVerts < 0 {
		return newErrorMesh("failed to parse number of vertices")
	}
	mb.currElement = elementVert
	mb.haveVerts = mb.nVerts > 0
	return nil
}

func (mb *meshBuilder) readProperty(line string) error {
	if mb.currElement == elementIgnore {
		return nil
	}
	if mb.currElement == elementVert {
		var propTypeName, name string
		n, e := fmt.Sscanf(line, "property %s %s", &propTypeName, &name)
		if n != 2 || e != nil {
			return newErrorMesh("failed to parse: " + line)
		}
		propType, ok := typeMap[propTypeName]
		if !ok {
			return newErrorMesh("unknown property type: " + propTypeName)
		}
		if !mb.facesPhase {
			mb.addVertProperty(propType, name)
		}
	} else if mb.currElement == elementFace {
		if !strings.HasPrefix(line, "property list ") {
			return newErrorMesh("only lists supported for faces")
		}
		// should be a std-list
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
		for j := 0; j < mb.vertPropIdx; j++ {
			var val float32
			var err error
			if val, err = mb.nextFloat(); err != nil {
				return err
			}
			mb.addVertProp(i, int(mb.vertProp[j].propIdx), float32(val))
		}
	}
	return nil
}

func (mb *meshBuilder) readFaces() error {
	for i := 0; i < mb.nFaces; i++ {
		// we only expect lists here
		if !mb.scanner.Scan() {
			return newErrorMesh("unexpected end of file")
		}
		var cnt int64
		var err error
		if cnt, err = strconv.ParseInt(mb.scanner.Text(), 10, 32); err != nil {
			return newErrorMesh("could not convert `" + mb.scanner.Text() + "` to int")
		}
		if cnt < 3 {
			return newErrorMesh("a face must have at least 3 indices")
		}
		var p0, p1, p2 int
		if p0, err = mb.nextInt(); err != nil {
			return err
		}
		if p1, err = mb.nextInt(); err != nil {
			return err
		}
		for j := 2; j < int(cnt); j++ {
			if p2, err = mb.nextInt(); err != nil {
				return err
			}
			if err = mb.addTriangle(p0, p1, p2); err != nil {
				return err
			}
			p1 = p2
		}
	}
	mb.mesh.Tris = mb.mesh.Tris[0:mb.triIdx]
	return nil
}

func (mb *meshBuilder) addTriangle(p0, p1, p2 int) error {
	if mb.triIdx+1 == cap(mb.mesh.Tris) {
		tmp := make([]Triangle, 3*cap(mb.mesh.Tris)/2+1)
		copy(tmp, mb.mesh.Tris)
		mb.mesh.Tris = tmp
	}
	if p0 >= len(mb.mesh.Verts) || p1 >= len(mb.mesh.Verts) || p2 >= len(mb.mesh.Verts) {
		return newErrorMesh("vertex index out of range")
	}
	mb.mesh.Tris[mb.triIdx].P1 = &mb.mesh.Verts[p0]
	mb.mesh.Tris[mb.triIdx].P2 = &mb.mesh.Verts[p1]
	mb.mesh.Tris[mb.triIdx].P3 = &mb.mesh.Verts[p2]
	mb.triIdx += 1
	return nil
}

func (mb *meshBuilder) nextLine() (val string, err error) {
	if !mb.scanner.Scan() {
		return "", newErrorMesh("unexpected end of file")
	}
	//Trace.Printf(strings.Trim(mb.scanner.Text(), " \t\r\n"))
	return strings.Trim(mb.scanner.Text(), " \t\r\n"), nil
}

func (mb *meshBuilder) nextFloat() (val float32, err error) {
	var token string
	if token, err = mb.getToken(); err != nil {
		return 0, err
	}
	var result float64
	if result, err = strconv.ParseFloat(token, 32); err != nil {
		return 0, newErrorMesh("could not convert `" + token + "` to float")
	}
	return float32(result), nil
}

func (mb *meshBuilder) nextInt() (val int, err error) {
	var token string
	if token, err = mb.getToken(); err != nil {
		return 0, err
	}
	var result int64
	if result, err = strconv.ParseInt(token, 10, 32); err != nil {
		return 0, newErrorMesh("could not convert `" + token + "` to int")
	}
	return int(result), nil
}

func (mb *meshBuilder) getToken() (token string, err error) {
	if !mb.scanner.Scan() {
		return "", newErrorMesh("unexpected end of file")
	}
	return mb.scanner.Text(), nil
}

func (mb *meshBuilder) addVertProp(vertIdx, propIdx int, value float32) {
	switch propIdx {
	case propX:
		mb.mesh.Verts[vertIdx].X = value
	case propY:
		mb.mesh.Verts[vertIdx].Y = value
	case propZ:
		mb.mesh.Verts[vertIdx].Z = value
	}
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

// copied and modified from src/bufio/scan.go
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
