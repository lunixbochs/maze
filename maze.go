package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

const (
	UP = iota
	DOWN
	LEFT
	RIGHT
)

type Node struct {
	N, E, S, W *Node
	X, Y       int
}

func (n *Node) Dir(dir int) (int, int) {
	x, y := n.X, n.Y
	switch dir {
	case UP:
		y = n.Y - 1
	case DOWN:
		y = n.Y + 1
	case LEFT:
		x = n.X - 1
	case RIGHT:
		x = n.X + 1
	}
	return x, y
}

func (n *Node) Attach(dir int, other *Node) {
	switch dir {
	case UP:
		n.N = other
		other.S = n
	case DOWN:
		n.S = other
		other.N = n
	case LEFT:
		n.W = other
		other.E = n
	case RIGHT:
		n.E = other
		other.W = n
	}
}

type Maze struct {
	Graph  [][]*Node
	Start  *Node
	Height int
	Width  int
}

func (m *Maze) Bounds(x, y int) bool {
	return 0 <= x && x < m.Width && 0 <= y && y < m.Height
}

func (m *Maze) Neighbors(node *Node) (ret []int) {
	for dir := 0; dir < 4; dir++ {
		x, y := node.Dir(dir)
		if m.Bounds(x, y) && m.Graph[y][x] == nil {
			ret = append(ret, dir)
		}
	}
	return ret
}

func (m *Maze) Next(node *Node) *Node {
	choices := m.Neighbors(node)
	if len(choices) == 0 {
		return nil
	}
	r := randByte()
	dir := choices[r%len(choices)]
	x, y := node.Dir(dir)
	next := &Node{X: x, Y: y}
	node.Attach(dir, next)
	m.Graph[y][x] = next
	return next
}

func (m *Maze) Generate(x, y int) {
	if !m.Bounds(x, y) {
		log.Fatal("starting coordinate out of bounds")
	}
	var stack []*Node
	cur := &Node{X: x, Y: y}
	m.Start = cur
	m.Graph[y][x] = cur
	stack = append(stack, cur)
	for len(stack) > 0 {
		for len(m.Neighbors(cur)) > 0 {
			cur = m.Next(cur)
			stack = append(stack, cur)
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
}

func (m *Maze) Print() {
	fmt.Println()
	for y := 0; y < m.Height; y++ {
		fmt.Println(m.Graph[y])
	}
	fmt.Println()
}

func (n *Node) Draw(img *image.RGBA, spacing int) chan image.Rectangle {
	ret := make(chan image.Rectangle)
	go func() {
		x, y := n.X*spacing, n.Y*spacing
		if n.S == nil {
			ret <- image.Rect(x, y+spacing, x+spacing, y+spacing+1)
		}
		if n.N == nil {
			ret <- image.Rect(x, y, x+spacing, y+1)
		}
		if n.W == nil {
			ret <- image.Rect(x, y, x+1, y+spacing)
		}
		if n.E == nil {
			ret <- image.Rect(x+spacing, y, x+spacing+1, y+spacing)
		}
		close(ret)
	}()
	return ret
}

func (m *Maze) Draw(name string, spacing int) {
	height := m.Height * spacing
	width := m.Width * spacing
	img := image.NewRGBA(image.Rect(0, 0, width+1, height+1))
	c := &color.RGBA{255, 255, 255, 255}
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			node := m.Graph[y][x]
			if node != nil {
				for r := range node.Draw(img, spacing) {
					draw.Draw(img, r, &image.Uniform{c}, image.ZP, draw.Src)
				}
			}
		}
	}
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(file, img)
}

func NewMaze(height, width int) *Maze {
	var graph [][]*Node
	for i := 0; i < height; i++ {
		graph = append(graph, make([]*Node, width))
	}
	return &Maze{graph, nil, height, width}
}
