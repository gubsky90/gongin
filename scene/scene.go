package scene

import (
	// "fmt"
	"github.com/go-gl/mathgl/mgl32"
)

type Node struct {
	position mgl32.Vec3
}

func (node *Node) GetPosition() mgl32.Vec3 {
	return node.position
}

// ------------------------

type Container struct {
	Node
	children []interface{}
}

func (container *Container) Append(child interface{}) {
	container.children = append(container.children, child)
}

func (container *Container) GetChildren()  []interface{} {
	return container.children
}

// ------------------------

type Root struct {
	Container
}

func NewRoot() *Root {
	return &Root{}
}


