package vec32

import (
	"io"
	"strings"
	"testing"
)

const (
	headerStart = "ply\r" +
		"format ascii 1.0\r"
	headerEnd = "end_header\r"
)

func newReader(s string) io.Reader {
	return strings.NewReader(s)
}

func TestHeaderBasic(t *testing.T) {
	var cases = []struct {
		inputStr string
		errorRsp string
	}{
		{"", "EOF"},
		{"plx\r", "expected file-magic ply\\r"},
		{"ply\r" + "format ascii 1.1\r", "only format ascii 1.0 supported"},
		{headerStart + "foobar baz\r", "unexpected line in header: foobar baz"},
		{headerStart + headerEnd, ""},
		{headerStart + "comment foo bar baz\r" + headerEnd, ""},
		{headerStart + "element vertex 0\r" + headerEnd, ""},
		{headerStart + "element face 0\r" + headerEnd, ""},
		{headerStart + "element vertex x\r" + headerEnd, "failed to parse number of vertices"},
		{headerStart + "element face x\r" + headerEnd, "failed to parse number of faces"},
	}
	for i, tc := range cases {
		testError(t, i, tc.inputStr, tc.errorRsp)
	}
}

func TestVerts(t *testing.T) {
	header1Vert := headerStart + "element vertex 1\r"
	var cases = []struct {
		inputStr string
		errorRsp string
	}{
		{header1Vert +
			"property float x\r" +
			"property float y\r" +
			"property float z\r" +
			headerEnd + "0 0 0\r", ""},
		{header1Vert +
			"property float x\r" +
			"property float y\r" +
			headerEnd + "0 0 0\r", "invalid vertex definition (missing coordinate)"},
		{header1Vert +
			"property float x\r" +
			"property float y\r" +
			"property float z\r" +
			headerEnd, "unexpected end of file"},
		{header1Vert +
			"property float x\r" +
			"property float y\r" +
			"property float z\r" +
			headerEnd + "0 0\r", "unexpected end of line"},
		{header1Vert +
			"property float x\r" +
			"property float y\r" +
			"property float z\r" +
			headerEnd + "0 0 0\r", ""},
		{header1Vert +
			"property float x\r" +
			"property float y\r" +
			"property float z\r" +
			headerEnd + "0 0 0", "unexpected end of file"},
	}
	for i, tc := range cases {
		testError(t, i, tc.inputStr, tc.errorRsp)
	}
}

func testError(t *testing.T, tc int, mesh, errTest string) {
	_, err := ReadPLY(newReader(mesh))
	if err == nil && errTest == "" {
		return
	}
	if err == nil || Strcmp(errTest, err.Error()) != 0 {
		errString := "<nil>"
		if err != nil {
			errString = err.Error()
		}
		t.Errorf("tc %d: Error on testing error response. expected: \"%s\", got \"%s\"",
			tc+1, errTest, errString)
	}
}
