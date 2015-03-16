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

// Intersect a triangle by a ray
func (tri *Triangle) Intersect(r *Ray, i *Intersection) float32 {
	var e1, e2, P, Q, T Vec3
	var t, det, inv_det, u, v float32
	Sub3(tri.P2, tri.P1, &e1)
	Sub3(tri.P3, tri.P1, &e2)
	Cross3(&r.N, &e2, &P)
	det = e1.Dot(&P)
	// culling backside
	if det > 100*EPS {
		return Inf(1)
	}
	inv_det = 1 / det
	Sub3(&r.P0, tri.P1, &T)

	u = T.Dot(&P) * inv_det
	if u < 0 || u > 1 {
		return Inf(1)
	}

	Cross3(&T, &e1, &Q)
	v = r.N.Dot(&Q) * inv_det
	if v < 0 || u+v > 1 {
		return Inf(1)
	}
	t = e2.Dot(&Q) * inv_det

	if t > 100*EPS {
		i.u = u
		i.v = v
		i.t = t
		return t
	} else {
		return Inf(1)
	}
}
