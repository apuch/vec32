package vec32

// Get the boundingBox around a triangle
func (tri *Triangle) OrthoBox() OrthoBox {
	return OrthoBox{Vec3{
		Min(tri.P1.X, Min(tri.P2.X, tri.P3.X)),
		Min(tri.P1.Y, Min(tri.P2.Y, tri.P3.Y)),
		Min(tri.P1.Z, Min(tri.P2.Z, tri.P3.Z)),
	}, Vec3{
		Max(tri.P1.X, Max(tri.P2.X, tri.P3.X)),
		Max(tri.P1.Y, Max(tri.P2.Y, tri.P3.Y)),
		Max(tri.P1.Z, Max(tri.P2.Z, tri.P3.Z)),
	}}
}
