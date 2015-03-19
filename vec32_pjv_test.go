package vec32

import (
	"io"
	"strings"
	"testing"
)

func newReader(s string) io.Reader {
	return strings.NewReader(s)
}

func TestHeaderBasic(t *testing.T) {
	headerStart := "ply\r" +
		"format ascii 1.0\r"
	headerEnd := "end_header\r"
	testError(t, "", "EOF")
	testError(t, "plx\r", "expected file-magic ply\\r")
	testError(t, "ply\r"+
		"format ascii 1.1\r", "only format ascii 1.0 supported")
	testError(t, headerStart+"foobar baz\r", "unexpected line in header: foobar baz")
	testError(t, headerStart+headerEnd, "")
	testError(t, headerStart+"comment foo bar baz\r"+headerEnd, "")
	testError(t, headerStart+"element vertex 0\r"+headerEnd, "")
	testError(t, headerStart+"element face 0\r"+headerEnd, "")
	testError(t, headerStart+"element vertex x\r"+headerEnd, "failed to parse number of vertices")
	testError(t, headerStart+"element face x\r"+headerEnd, "failed to parse number of faces")

}

func testError(t *testing.T, mesh, errTest string) {
	_, err := ReadPLY(newReader(mesh))
	if err == nil && errTest == "" {
		return
	}
	if err == nil || Strcmp(errTest, err.Error()) != 0 {
		errString := "<nil>"
		if err != nil {
			errString = err.Error()
		}
		t.Errorf("Error on testing error response. expected: \"%s\", got \"%s\"",
			errTest, errString)
	}
}
