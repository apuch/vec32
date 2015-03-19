package vec32

import (
	"io"
	"strings"
	"testing"
)

func newReader(s string) io.Reader {
	return strings.NewReader(s)
}

func TestEmptyString(t *testing.T) {
	testError(t, "", "EOF")
	testError(t, "plx\r", "expected file-magic ply\\r")
	testError(t, "ply\rformat ascii 1.1\r", "only format ascii 1.0 supported")
}

func testError(t *testing.T, mesh, errTest string) {
	_, err := ReadPLY(newReader(mesh))
	if err == nil || Strcmp(errTest, err.Error()) != 0 {
		errString := "<nil>"
		if err != nil {
			errString = err.Error()
		}
		t.Errorf("Error on testing error response. expected: \"%s\", got \"%s\"",
			errTest, errString)
	}
}
