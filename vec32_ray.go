package vec32

// Create a new ray
func NewRay(p0, p1 *Vec3) *Ray {
	r := new(Ray)
	r.P0 = *p0
	r.N = p1.Sub(*p0).Normalize()
	return r
}

// Where will the ray be at value t
func (r *Ray) At(t float32, v *Vec3) {
	v1 := r.N.Scale(t)
	Add3(&r.P0, &v1, v)
}
