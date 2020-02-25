package graph_test

import (
	"fmt"
	"math/rand"
	"testing"

	graph "github.com/oligoden/math-graph"
)

func TestBasic(t *testing.T) {
	g := graph.New()
	g.Add("a")
	g.Add("b")
	g.Link("a", "b")
	err := g.Evaluate()
	if err != nil {
		t.Error(err)
	}
	if len(g.StartNodes()) != 1 {
		t.Error("expected 1 start node, got", len(g.StartNodes()))
	}
}

func TestCyclic(t *testing.T) {
	g := graph.New()
	g.Add("a")
	g.Add("b")
	g.Add("c")
	g.Link("a", "b")
	g.Link("b", "c")
	g.Link("c", "a")
	err := g.Evaluate()
	if err == nil {
		t.Error("expected cyclic error")
	}
}

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

func Benchmark10N32E(b *testing.B) {
	ns := "abcdefghij"
	benchmarkEvaluate(ns, b)
}

func Benchmark20N162E(b *testing.B) {
	ns := "abcdefghijklmnopqrst"
	benchmarkEvaluate(ns, b)
}

func Benchmark30N392E(b *testing.B) {
	ns := "abcdefghijklmnopqrstuvwxyzABCD"
	benchmarkEvaluate(ns, b)
}

func Benchmark40N722E(b *testing.B) {
	ns := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMN"
	benchmarkEvaluate(ns, b)
}

func Benchmark50N1152E(b *testing.B) {
	ns := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWX"
	benchmarkEvaluate(ns, b)
}

func Benchmark60N1682E(b *testing.B) {
	ns := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567"
	benchmarkEvaluate(ns, b)
}

func benchmarkEvaluate(ns string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		g := graph.New()
		for i := 0; i < len(ns); i++ {
			g.Add(string(ns[i]))
		}
		for i := 0; i < len(ns)-2; i++ {
			for j := 0; j < len(ns)-i-2; j++ {
				g.Link(string(ns[i]), string(ns[i+1+rand.Intn(len(ns)-i-1)]))
			}
		}
		g.Evaluate()
	}
}
