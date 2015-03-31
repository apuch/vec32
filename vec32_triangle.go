package vec32

// Get the boundingBox around a triangle
func (tri *Triangle) OrthoBox(b *OrthoBox) {
	b.P0.X = Min(tri.P1.X, Min(tri.P2.X, tri.P3.X))
	b.P0.Y = Min(tri.P1.Y, Min(tri.P2.Y, tri.P3.Y))
	b.P0.Z = Min(tri.P1.Z, Min(tri.P2.Z, tri.P3.Z))
	b.P1.X = Max(tri.P1.X, Max(tri.P2.X, tri.P3.X))
	b.P1.Y = Max(tri.P1.Y, Max(tri.P2.Y, tri.P3.Y))
	b.P1.Z = Max(tri.P1.Z, Max(tri.P2.Z, tri.P3.Z))
}
