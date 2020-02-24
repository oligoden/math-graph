package graph_test

import (
	"fmt"
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

	testRun := make(map[string]bool)
	f := func(name string) error {
		if name == "c" && !testRun["b"] {
			return fmt.Errorf("c before b")
		}
		testRun[name] = true
		return nil
	}
	err = g.CompileRun(f)
	if err != nil {
		t.Error(err)
	}
	if !testRun["a"] {
		t.Error("expected test function a to be run")
	}

	testRun = make(map[string]bool)
	var testFlag bool
	f = func(name string) error {
		fmt.Println(name, testFlag)
		if name == "c" {
			testFlag = true
		}
		testRun[name] = testFlag
		return nil
	}
	err = g.SetRun(f, "a")
	if err != nil {
		t.Error(err)
	}
	testFlag = false
	err = g.SetRun(f, "b")
	if err != nil {
		t.Error(err)
	}

	if testRun["a"] {
		t.Error("expected test function a not to be run")
	}
	if testRun["b"] {
		t.Error("expected test function b not to be run")
	}
	if !testRun["c"] {
		t.Error("expected test function c to be run")
	}
	if !testRun["d"] {
		t.Error("expected test function d to be run")
	}
}
