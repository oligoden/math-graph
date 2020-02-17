package graph

import (
	"fmt"
)

type Graph struct {
	nodes map[string]*Node
	adj   map[string]map[string]uint
}

func New() *Graph {
	g := &Graph{}
	g.nodes = make(map[string]*Node)
	g.adj = make(map[string]map[string]uint)
	return g
}

func (g *Graph) Nodes() map[string]*Node {
	r := make(map[string]*Node)
	for name, node := range g.nodes {
		r[name] = node
	}
	return r
}

func (g *Graph) StartNodes() map[string]*Node {
	r := make(map[string]*Node)
	for name, node := range g.nodes {
		if node.set && node.level == 0 {
			r[name] = node
		}
	}
	return r
}

func (g *Graph) Add(name string) error {
	g.nodes[name] = &Node{}
	g.nodes[name].children = make(map[string]*Node)
	g.adj[name] = make(map[string]uint)
	for k := range g.adj {
		g.adj[k][name] = 0
		g.adj[name][k] = 0
	}
	return nil
}

func (g *Graph) Link(from, to string) error {
	g.adj[from][to] = 1
	return nil
}

func (g *Graph) Evaluate() error {
	err := g.evaluate(g.adj, 1)
	if err != nil {
		return err
	}
	return nil
}

func (g *Graph) evaluate(pre map[string]map[string]uint, depth uint) error {
	depth++
	if depth == 100 {
		return fmt.Errorf("max walk depht reached")
	}
	done := true
	prd := make(map[string]map[string]uint)
	for row := range pre {
		prd[row] = make(map[string]uint)
		for col := range pre[row] {
			for elm := range pre {
				prd[row][col] += pre[row][elm] * g.adj[elm][col]
			}
			if prd[row][col] > 0 {
				if row == col {
					return fmt.Errorf("cyclical")
				}
				done = false
			}
		}
	}

	if done {
		return nil
	}

	err := g.evaluate(prd, depth)
	if err != nil {
		return err
	}
	fmt.Println(depth)
	fmt.Println(prd)
	for row := range prd {
		if g.nodes[row].set {
			continue
		}
		for col := range prd[row] {
			if prd[row][col] > 0 {
				fmt.Println(row, col)
				g.nodes[row].set = true
				g.set(row, 0)
			}
		}
	}
	return nil
}

func (g *Graph) set(node string, level uint) {
	g.nodes[node].set = true
	g.nodes[node].level = level
	for child := range g.adj[node] {
		if g.adj[node][child] > 0 {
			g.nodes[node].children[child] = g.nodes[child]
			g.set(child, level+1)
		}
	}
}

func (g *Graph) Run(f func(node string)) {
	for name, node := range g.StartNodes() {
		g.run(f, name, node)
	}
}

func (g *Graph) run(f func(node string), name string, node *Node) {
	f(name)

	for name, node := range node.children {
		g.run(f, name, node)
	}
}

type Node struct {
	level    uint
	set      bool
	children map[string]*Node
}
