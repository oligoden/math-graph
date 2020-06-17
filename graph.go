package graph

import (
	"errors"
	"fmt"
)

type Graph struct {
	nodes  map[string]*Node
	levels []map[string]*Node
	adj    map[string]map[string]uint
}

type Node struct {
	level    uint
	set      bool
	children map[string]*Node
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
	if _, fnd := g.adj[from]; !fnd {
		return errors.New("the node " + from + " does not exist")
	}

	if _, fnd := g.adj[from][to]; !fnd {
		return errors.New("the node " + to + " does not exist")
	}

	g.adj[from][to] = 1
	return nil
}

func (g *Graph) Unlink(from, to string) error {
	if _, fnd := g.adj[from]; !fnd {
		return errors.New("the node " + from + " does not exist")
	}

	if _, fnd := g.adj[from][to]; !fnd {
		return errors.New("the node " + to + " does not exist")
	}

	g.adj[from][to] = 0
	return nil
}

func (g *Graph) Evaluate() error {
	for key := range g.nodes {
		g.nodes[key].set = false
		g.nodes[key].children = make(map[string]*Node)
	}

	err := g.evaluate(g.adj, 1)
	if err != nil {
		return err
	}
	for row := range g.adj {
		if g.nodes[row].set {
			continue
		}
		for col := range g.adj[row] {
			if g.adj[row][col] > 0 {
				g.set(row, 0)
			}
		}
		if !g.nodes[row].set {
			g.nodes[row].set = true
			g.nodes[row].level = 0

			if g.levels[0] == nil {
				g.levels[0] = make(map[string]*Node)
			}
			g.levels[0][row] = g.nodes[row]
		}
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
		g.levels = make([]map[string]*Node, depth)
		return nil
	}

	err := g.evaluate(prd, depth)
	if err != nil {
		return err
	}
	for row := range prd {
		if g.nodes[row].set {
			continue
		}
		for col := range prd[row] {
			if prd[row][col] > 0 {
				g.set(row, 0)
			}
		}
	}
	return nil
}

func (g *Graph) set(node string, level uint) {
	if g.nodes[node].set {
		return
	}

	g.nodes[node].set = true
	g.nodes[node].level = level

	if g.levels[level] == nil {
		g.levels[level] = make(map[string]*Node)
	}
	g.levels[level][node] = g.nodes[node]

	for child := range g.adj[node] {
		if g.adj[node][child] > 0 {
			g.nodes[node].children[child] = g.nodes[child]
			g.set(child, level+1)
		}
	}
}

// CompileRun executes a runfunc on every node
// in topological order, starting at the 0 level nodes
// through to the highest level nodes. This guarantees
// that all posible parent nodes are executed before child nodes.
// However, walks are not followed. The nodes are simply scanned over.
// CompileRun is usfull for a naive bulk execution in topological order.
func (g *Graph) CompileRun(f func(string) error) error {
	for _, nodes := range g.levels {
		for name := range nodes {
			err := f(name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SetRun executes only given nodes and its children.
// Note that there might be parent nodes that are not executed
// before a child node. Yet, walks are strictly followed.
// SetRun is usefull for setting properties on all child nodes
// of a specific node.
func (g *Graph) SetRun(f func(string) error, name string) error {
	err := f(name)
	if err != nil {
		return err
	}

	for n := range g.nodes[name].children {
		g.SetRun(f, n)
	}

	return nil
}

// ReverseRun executes only given nodes and its parents.
func (g *Graph) ReverseRun(f func(string) error, name string) error {
	err := f(name)
	if err != nil {
		return err
	}

	for node := range g.nodes {
		if g.adj[node][name] > 0 {
			g.ReverseRun(f, node)
		}
	}

	return nil
}
