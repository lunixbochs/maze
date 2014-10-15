package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func (m *Maze) Solve(x, y, xe, ye int) []*Node {
	if !m.Bounds(x, y) || !m.Bounds(xe, ye) {
		log.Fatal("coordinate(s) out of bounds")
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
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <width> <height>\n", os.Args[0])
		return
	}
	width, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	height, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	m := NewMaze(height, width)
	m.Generate(0, 0)
	m.Draw("maze.png", 3)
	solution := m.Solve(0, 0, m.Width-1, m.Height-1)
	fmt.Println(solution)
}
