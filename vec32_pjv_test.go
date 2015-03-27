package vec32

import (
	"io"
	"os"
	"strings"
	"testing"
)

const (
	headerStart = "ply\r" +
		"format ascii 1.0\r"
	headerEnd      = "end_header\r"
	header1Vert    = headerStart + "element vertex 1\r"
	header2Vert    = headerStart + "element vertex 2\r"
	header6Vert    = headerStart + "element vertex 6\r"
	validVertCoord = "property float x\r" +
		"property float y\r" +
		"property float z\r"
	valid1VertHeader = header1Vert + validVertCoord + headerEnd
	valid2VertHeader = header2Vert + validVertCoord + headerEnd
	valid1FaceHeader = header6Vert + validVertCoord +
		"element face 1\rproperty list uchar int\r" + headerEnd
	valid3FaceHeader = header6Vert + validVertCoord +
		"element face 3\rproperty list uchar int\r" + headerEnd
	valid6Coords = "0 0 0\r" +
		"0 1 0\r" +
		"1 1 0\r" +
		"1 0 0\r" +
		"0 2 0\r" +
		"1 2 0\r"
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
		{header1Vert +
			"property float foobar\r" +
			validVertCoord +
			headerEnd + "0 0 0 0\r", ""},
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

func TestFaces(t *testing.T) {
	header := valid1FaceHeader + valid6Coords
	header3 := valid3FaceHeader + valid6Coords
	var cases = []struct {
		inputStr string
		err      string
		numFaces int
	}{
		{header + "3 0 1 2\r", "", 1},
		{header, "unexpected end of file", 0},
		{header + "3\r", "unexpected end of file", 0},
		{header + "3 0\r", "unexpected end of file", 0},
		{header + "3 0 1\r", "unexpected end of file", 0},
		{header + "2 0 1\r", "a face must have at least 3 indices", 0},
		{header3 + "3 0 1 2\r3 0 2 3\r3 1 4 2\r", "", 3},
		{header + "4 0 1 2 3\r", "", 2},
		{header + "5 0 1 2 3 4\r", "", 3},
		{header + "6 0 1 2 3 4 5\r", "", 4},
		{header + "3 0 1 6\r", "vertex index out of range", 4},
	}
	for i, tc := range cases {
		m := testError(t, i, tc.inputStr, tc.err)
		if m == nil {
			continue
		}
		if len(m.Tris) != tc.numFaces {
			t.Errorf("tc %d: expected %d faces, got %d", i+1, tc.numFaces, len(m.Tris))
		}
	}
}

func TestRealSamples(t *testing.T) {
	var cases = []struct {
		file     string
		numVerts int
		numTris  int
	}{
		{"paulbourke.net.sample1.ply", 6, 12},
	}
	for i, tc := range cases {
		var f *os.File
		var err error
		var m *Mesh
		if f, err = os.Open("test/ply/" + tc.file); err != nil {
			t.Errorf("tc %d: could not open file %s: %s ", i, tc.file, err.Error())
			continue
		}
		if m, err = ReadPLY(f); err != nil {
			t.Errorf("tc %d: unexpected error on reading mesh: %s", i+1, err.Error())
		}
		if m == nil {
			continue
		}
		if len(m.Verts) != tc.numVerts {
			t.Errorf("tc %d: expected %d verts, got %d", i+1, tc.numVerts, len(m.Verts))
		}
		if len(m.Tris) != tc.numTris {
			t.Errorf("tc %d: expected %d faces, got %d", i+1, tc.numTris, len(m.Tris))
		}
	}
}

func testError(t *testing.T, tc int, mesh, errTest string) *Mesh {
	m, err := ReadPLY(newReader(mesh))
	if err == nil && errTest == "" {
		return m
	}
	if err == nil || Strcmp(errTest, err.Error()) != 0 {
		errString := "<nil>"
		if err != nil {
			errString = err.Error()
		}
		t.Errorf("tc %d: Error on testing error response. expected: \"%s\", got \"%s\"",
			tc+1, errTest, errString)
	}
	return m
}
