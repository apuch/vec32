package vec32

import (
	"bufio"
	"io"
)

type Mesh struct {
	Points []Vec3
}

func ReadPLY(r io.Reader) (m *Mesh, err error) {
	rd := bufio.NewReader(r)

	var line string
	line, err = rd.ReadString('\r')
	if err != nil {
		return nil, newErrorMesh(err.Error())
	}
	if Strcmp(line, "ply\r") != 0 {
		return nil, newErrorMesh("expected file-magic ply\\r")
	}
	line, err = rd.ReadString('\r')
	if err != nil {
		return nil, newErrorMesh(err.Error())
	}
	if line != "format ascii 1.0\r" {
		return nil, newErrorMesh("only format ascii 1.0 supported")
	}
	return nil, nil
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
