package vec32

// Create a new ray
func NewRay(p0, p1 *Vec3) *Ray {
	r := new(Ray)
	r.P0 = *p0
	r.N = *p1.Sub(p0).Normalize()
	return r
}

// Where will the ray be at value t
func (r *Ray) At(t float32, v *Vec3) {
	Add3(&r.P0, r.N.Scale(t), v)
}

// ray-triangle-intersection by Möller–Trumbore
//
// returns inf if triangle isn't hit. Will always be greater zero
func (r *Ray) Intersect(tri *Triangle, i *Intersection, tMax float32) float32 {
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

	if t > 100*EPS && t < tMax {
		i.u = u
		i.v = v
		i.t = t
		return t
	} else {
		return Inf(1)
	}
}
