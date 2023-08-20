package lib

type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) Scale(factor float64) Vector3 {
	return Vector3{X: v.X * factor, Y: v.Y * factor, Z: v.Z * factor}
}

func (v Vector3) Add(vector Vector3) Vector3 {
	return Vector3{X: v.X + vector.X, Y: v.Y + vector.Y, Z: v.Z + vector.Z}
}

func (v Vector3) To2() Vector2 {
	return Vector2{X: v.X, Y: v.Y}
}

type Vector2 struct {
	X, Y float64
}

func (v Vector2) Scale(factor float64) Vector2 {
	return Vector2{X: v.X * factor, Y: v.Y * factor}
}

func (v Vector2) Add(vector Vector2) Vector2 {
	return Vector2{X: v.X + vector.X, Y: v.Y + vector.Y}
}
