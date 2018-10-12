package render

// Rasterizator ...
type Rasterizator struct {
	shader *Shader
}

// NewRasterizator ...
func NewRasterizator(shader *Shader) *Rasterizator {
	r := Rasterizator{}
	r.shader = shader
	return &r
}

// DrawMesh ...
func (r *Rasterizator) DrawMesh(mesh *Mesh) {
	r.shader.Use()

	time := float32(getTime() - getStartTime()) / 1000

	r.shader.Set1f("iTime", time)

	mesh.Draw()
}
