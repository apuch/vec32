package vec32

import (
	"io"
	"strings"
	"testing"
)

const (
	headerStart = "ply\r" +
		"format ascii 1.0\r"
	headerEnd      = "end_header\r"
	header1Vert    = headerStart + "element vertex 1\r"
	header2Vert    = headerStart + "element vertex 2\r"
	validVertCoord = "property float x\r" +
		"property float y\r" +
		"property float z\r"
	valid1VertHeader = header1Vert + validVertCoord + headerEnd
	valid2VertHeader = header2Vert + validVertCoord + headerEnd
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

func TestVertsHeader(t *testing.T) {
	var cases = []struct {
		inputStr string
		errorRsp string
	}{
		{valid1VertHeader + "0 0 0\r", ""},
		{header1Vert +
			"property float x\r" +
			"property float y\r" +
			headerEnd + "0 0 0\r", "invalid vertex definition (missing coordinate)"},
		{valid1VertHeader, "unexpected end of file"},
		{valid1VertHeader + "0 0\r", "unexpected end of file"},
		{valid1VertHeader + "0 0 0\r", ""},
		{valid1VertHeader + "0  0  0\r", ""},
		{valid1VertHeader + "0 0 0 \r", ""},
		{valid1VertHeader + " 0 0 0\r", ""},
		{valid2VertHeader + " 0 0 0\r0 0 0\r", ""},
		{valid1VertHeader + "0 0 0", ""},
		{valid1VertHeader + "0 0 X", "could not convert `X` to float"},
	}
	for i, tc := range cases {
		testError(t, i, tc.inputStr, tc.errorRsp)
	}
}

func TestVerts(t *testing.T) {
	m, err := ReadPLY(newReader(valid2VertHeader + "1 2 3\r4 5 6\r"))
	if err != nil {
		t.Errorf("error at reading PLY: %s", err.Error())
		return
	}
	if len(m.Verts) != 2 {
		t.Errorf("Expected 2 verts, got %d", len(m.Verts))
	}
	if !m.Verts[0].IsEqual(&Vec3{1, 2, 3}) || !m.Verts[1].IsEqual(&Vec3{4, 5, 6}) {
		t.Errorf("read wrong vectors (%s and %s)",
			m.Verts[0].String(), m.Verts[1].String())
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
