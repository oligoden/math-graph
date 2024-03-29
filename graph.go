package graph

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

type Graph struct {
	nodes   map[string]*Node
	indexes map[int]string
	levels  []map[int]*Node
	adj     [][]uint
}

type Node struct {
	level       uint
	set         bool
	childrenSet bool
	children    map[string]*Node
	adjIndex    int
}

func New() *Graph {
	g := &Graph{}
	g.nodes = make(map[string]*Node)
	g.indexes = make(map[int]string)
	g.adj = [][]uint{}
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
	g.nodes[name].adjIndex = len(g.adj)
	g.indexes[g.nodes[name].adjIndex] = name
	for i := range g.adj {
		g.adj[i] = append(g.adj[i], 0)
	}
	g.adj = append(g.adj, make([]uint, len(g.adj)+1))

	return nil
}

func (g *Graph) Rename(oldName, newName string) error {
	g.nodes[newName] = g.nodes[oldName]
	delete(g.nodes, oldName)
	g.indexes[g.nodes[newName].adjIndex] = newName
	return nil
}

func (g *Graph) Remove(name string) error {
	if _, fnd := g.nodes[name]; !fnd {
		return errors.New("the node " + name + " does not exist")
	}

	i := g.nodes[name].adjIndex
	l := g.nodes[name].level

	delete(g.nodes, name)
	delete(g.levels[l], i)
	for j := i; j < len(g.adj)-1; j++ {
		g.indexes[j] = g.indexes[j+1]
		g.nodes[g.indexes[j+1]].adjIndex = j
	}
	delete(g.indexes, len(g.adj))

	// delete row from adj matrix
	if i < len(g.adj)-1 {
		copy(g.adj[i:], g.adj[i+1:])
	}
	g.adj[len(g.adj)-1] = nil
	g.adj = g.adj[:len(g.adj)-1]

	// delete col from adj matrix
	for j := range g.adj {
		if i < len(g.adj[j])-1 {
			copy(g.adj[j][i:], g.adj[j][i+1:])
		}
		g.adj[j][len(g.adj[j])-1] = 0
		g.adj[j] = g.adj[j][:len(g.adj[j])-1]
	}

	return nil
}

func (g *Graph) Link(from, to string) error {
	if _, fnd := g.nodes[from]; !fnd {
		return errors.New("the node " + from + " does not exist")
	}

	if _, fnd := g.nodes[to]; !fnd {
		return errors.New("the node " + to + " does not exist")
	}

	g.adj[g.nodes[from].adjIndex][g.nodes[to].adjIndex] = 1
	return nil
}

func (g *Graph) Unlink(from, to string) error {
	if _, fnd := g.nodes[from]; !fnd {
		return errors.New("the node " + from + " does not exist")
	}

	if _, fnd := g.nodes[to]; !fnd {
		return errors.New("the node " + to + " does not exist")
	}

	g.adj[g.nodes[from].adjIndex][g.nodes[to].adjIndex] = 0
	return nil
}

func (g *Graph) Evaluate() error {
	for key := range g.nodes {
		g.nodes[key].set = false
		g.nodes[key].childrenSet = false
		g.nodes[key].children = make(map[string]*Node)
	}

	err := g.evaluate(g.adj, 1)
	if err != nil {
		return err
	}

	for _, node := range g.nodes {
		row := node.adjIndex
		for col := range g.adj[row] {
			if g.adj[row][col] > 0 {
				g.setNode(g.indexes[col], 1)
			}
		}
	}

	// Setting starting nodes (nodes that is not set yet) and children
	for key, node := range g.nodes {
		if !g.nodes[key].set {
			g.nodes[key].set = true
			g.nodes[key].level = 0

			if g.levels[0] == nil {
				g.levels[0] = make(map[int]*Node)
			}
			g.levels[0][node.adjIndex] = g.nodes[key]

			if !g.nodes[key].childrenSet {
				g.setChildren(key)
			}
		}
	}
	return nil
}

func (g *Graph) setChildren(node string) {
	row := g.nodes[node].adjIndex
	for col := range g.adj[row] {
		if g.adj[row][col] > 0 {
			child := g.indexes[col]
			g.nodes[node].children[child] = g.nodes[child]
			if !g.nodes[child].childrenSet {
				g.setChildren(child)
				g.nodes[child].childrenSet = true
			}
		}
	}
}

func (g *Graph) evaluate(pre [][]uint, depth uint) error {
	depth++

	if depth == 100 {
		return fmt.Errorf("max walk depth reached")
	}

	// do matrix multiplication prd = pre * adj
	// done will stay true if there are no walks of length depth
	length := len(pre)
	done := true
	prd := make([][]uint, length)
	for row := range pre {
		prd[row] = make([]uint, length)
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
		g.levels = make([]map[int]*Node, depth)
		return nil
	}

	err := g.evaluate(prd, depth)
	if err != nil {
		return err
	}

	for n := range g.nodes {
		row := g.nodes[n].adjIndex
		for col := range prd[row] {
			if prd[row][col] > 0 {
				g.setNode(g.indexes[col], depth)
			}
		}
	}
	return nil
}

func (g *Graph) setNode(key string, level uint) {
	if g.nodes[key].set {
		return
	}

	g.nodes[key].set = true
	g.nodes[key].level = level

	if g.levels[level] == nil {
		g.levels[level] = make(map[int]*Node)
	}
	g.levels[level][g.nodes[key].adjIndex] = g.nodes[key]
}

// CompileRun executes a runfunc on every node
// in topological order, starting at the 0 level nodes
// through to the highest level nodes. This guarantees
// that all posible parent nodes are executed before child nodes.
// However, walks are not followed. The nodes are simply scanned over.
// CompileRun is usfull for a naive bulk execution in topological order.
func (g *Graph) CompileRun(f func(string) error) error {
	for _, keys := range g.levels {
		for key := range keys {
			err := f(g.indexes[key])
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

	node, ok := g.nodes[name]
	if !ok {
		return errors.New("node " + name + " does not exist")
	}

	for n := range node.children {
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

	n, ok := g.nodes[name]
	if !ok {
		return errors.New("node " + name + " does not exist")
	}
	col := n.adjIndex

	for row := range g.adj {
		if g.adj[row][col] > 0 {
			g.ReverseRun(f, g.indexes[row])
		}
	}

	return nil
}

func (g *Graph) Output(as ...string) *bytes.Buffer {
	buf := bytes.NewBufferString("digraph {\n")

	if len(as)%2 == 1 {
		as = as[0 : len(as)-1]
	}

	for name := range g.nodes {
		fmt.Fprintf(buf, "\t\"%s\"", name)
		for i := 0; i < len(as)/2; i++ {
			match, _ := regexp.MatchString(as[2*i], name)
			if match {
				fmt.Fprintf(buf, " %s", as[2*i+1])
				break
			}
		}
		fmt.Fprintln(buf, ";")
	}

	for _, node := range g.nodes {
		row := node.adjIndex
		for col := range g.adj[row] {
			if g.adj[row][col] == 1 {
				fmt.Fprintf(buf, "\t\"%s\" -> \"%s\";\n", g.indexes[row], g.indexes[col])
			}
		}
	}

	buf.WriteString("}")
	return buf
}

func (g *Graph) Graph() ([]string, [][]string) {
	nds := []string{}
	lks := [][]string{}

	for name, node := range g.nodes {
		nds = append(nds, name)
		row := node.adjIndex
		for col := range g.adj[row] {
			if g.adj[row][col] == 1 {
				lks = append(lks, []string{g.indexes[row], g.indexes[col]})
			}
		}
	}
	return nds, lks
}
