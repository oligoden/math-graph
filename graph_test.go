package graph_test

import (
	"testing"

	graph "github.com/oligoden/math-graph"
)

func Test(t *testing.T) {
	g := graph.New()
	g.Add("a")
	g.Add("b")
	g.Add("c")
	g.Add("d")
	g.Add("e")

	nodes := g.Nodes()
	if len(nodes) != 5 {
		t.Error("expected 4 nodes")
	}

	g.Link("a", "c")
	g.Link("b", "c")
	g.Link("b", "d")
	g.Link("c", "d")
	g.Link("d", "e")
	err := g.Evaluate()
	if err != nil {
		t.Error(err)
	}

	nodes = g.StartNodes()
	if len(nodes) != 2 {
		t.Error("expected 2 start nodes")
	}

	testFunc := make(map[string]bool)
	f := func(name string) {
		testFunc[name] = true
	}
	g.Run(f)
	if !testFunc["a"] {
		t.Error("expected test function a to be run")
	}
}
