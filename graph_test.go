package graph_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	graph "github.com/oligoden/math-graph"

	"github.com/pkg/profile"
)

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
		t.Fatal("expected cyclic error")
	}
}

func TestLinkError(t *testing.T) {
	g := graph.New()
	g.Add("a")

	err := g.Link("b", "a")
	if err == nil {
		t.Fatal("expected error")
	}
	exp := "the node b does not exist"
	got := err.Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	err = g.Link("a", "b")
	if err == nil {
		t.Fatal("expected error")
	}
	exp = "the node b does not exist"
	got = err.Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestRuns(t *testing.T) {
	var got string

	g := graph.New()
	g.Add("a")
	g.Add("b")
	g.Add("c")
	g.Add("d")
	g.Add("e")
	g.Add("f")

	nodes := g.Nodes()
	if len(nodes) != 6 {
		t.Error("expected 6 nodes")
	}

	err := g.Evaluate()
	if err != nil {
		t.Error(err)
	}

	nodes = g.StartNodes()
	if len(nodes) != 6 {
		t.Error("expected 6 start nodes")
	}

	g.Link("a", "c")
	g.Link("b", "c")
	g.Link("b", "d")
	g.Link("c", "d")
	g.Link("d", "e")
	err = g.Evaluate()
	if err != nil {
		t.Error(err)
	}

	nodes = g.StartNodes()
	if len(nodes) != 3 {
		t.Error("expected 3 start nodes")
	}

	g.CompileRun(func(s string) error { got += s; return nil })
	if len(got) != 6 {
		t.Errorf(`expected %d characters, got "%s"`, 6, got)
	}
	if !strings.Contains(got[0:3], "a") {
		t.Errorf(`expected "a"`)
	}
	if !strings.Contains(got[0:3], "b") {
		t.Errorf(`expected "b"`)
	}
	if !strings.Contains(got[0:3], "f") {
		t.Errorf(`expected "f"`)
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

	var testString string
	f = func(name string) error {
		testString += name
		return nil
	}

	err = g.ReverseRun(f, "z")
	if err == nil {
		t.Error("expected and error")
	}
	testString = ""

	err = g.ReverseRun(f, "c")
	if err != nil {
		t.Error(err)
	}
	if testString != "cab" && testString != "cba" {
		t.Error(`expected "cab" or "cba", got`, testString)
	}
}

func TestSmallRun1(t *testing.T) {
	g := graph.New()

	g.Add("a")
	g.Add("b")
	g.Add("c")

	g.Link("a", "b")
	g.Link("a", "c")
	g.Link("b", "c")

	err := g.Evaluate()
	if err != nil {
		t.Error(err)
	}

	testRun := []string{}
	f := func(name string) error {
		testRun = append(testRun, name)
		return nil
	}
	err = g.CompileRun(f)
	if err != nil {
		t.Error(err)
	}

	if len(testRun) != 3 {
		t.Errorf(`%+v`, g)
		t.Fatal("expected 3 nodes, got", len(testRun))
	}

	exp := "c"
	got := testRun[2]
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "a"
	got = testRun[0]
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
		t.Errorf(`%+v`, g)
		t.Errorf(`%+v`, g.Nodes()["a"])
		t.Errorf(`%+v`, g.Nodes()["b"])
		t.Errorf(`%+v`, g.Nodes()["c"])
	}

	testRun = []string{}
	err = g.ReverseRun(f, "c")
	if err != nil {
		t.Error(err)
	}

	if len(testRun) != 4 {
		t.Errorf(`%+v`, g)
		t.Fatal("expected 4 nodes, got", len(testRun))
	}

	exp = "a"
	got = testRun[3]
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "c"
	got = testRun[0]
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
		t.Errorf(`%+v`, g)
		t.Errorf(`%+v`, g.Nodes()["a"])
		t.Errorf(`%+v`, g.Nodes()["b"])
		t.Errorf(`%+v`, g.Nodes()["c"])
	}

	exp = `digraph {
	"a" -> "b";
	"a" -> "c";
	"b" -> "c";
}`
	got = g.Output().String()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestSmallRun2(t *testing.T) {
	g := graph.New()

	g.Add("a")
	g.Add("d")
	g.Add("b")
	g.Add("c")

	g.Link("a", "b")
	g.Link("a", "d")
	g.Link("b", "c")
	g.Link("c", "d")

	err := g.Evaluate()
	if err != nil {
		t.Error(err)
	}

	testRun := []string{}
	f := func(name string) error {
		testRun = append(testRun, name)
		return nil
	}
	err = g.CompileRun(f)
	if err != nil {
		t.Error(err)
	}

	if len(testRun) != 4 {
		t.Errorf(`%+v`, g)
		t.Fatal("expected 4 nodes, got", len(testRun))
	}

	exp := "d"
	got := testRun[3]
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "a"
	got = testRun[0]
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
		t.Errorf(`%+v`, g)
		t.Errorf(`%+v`, g.Nodes()["a"])
		t.Errorf(`%+v`, g.Nodes()["b"])
		t.Errorf(`%+v`, g.Nodes()["c"])
		t.Errorf(`%+v`, g.Nodes()["d"])
	}
}

func TestLargeRun(t *testing.T) {
	g := graph.New()

	g.Add("a")
	g.Add("b")
	g.Add("c")
	g.Add("d")
	g.Add("e")
	g.Add("f")
	g.Add("g")
	g.Add("h")
	g.Add("i")
	g.Add("j")
	g.Add("k")
	g.Add("l")
	g.Add("m")
	g.Add("n")
	g.Add("o")
	g.Add("p")
	g.Add("r")
	g.Add("s")
	g.Add("t")
	g.Add("u")
	g.Add("v")
	g.Add("x")
	g.Add("y")
	g.Add("z")
	g.Add("4")
	g.Add("5")
	g.Add("6")

	g.Link("a", "b")
	g.Link("a", "c")
	g.Link("c", "d")
	g.Link("a", "e")
	g.Link("e", "f")
	g.Link("e", "g")
	g.Link("g", "b")
	g.Link("f", "b")
	g.Link("e", "h")
	g.Link("h", "i")
	g.Link("i", "b")
	g.Link("h", "j")
	g.Link("j", "k")
	g.Link("k", "b")
	g.Link("h", "l")
	g.Link("l", "m")
	g.Link("m", "b")
	g.Link("l", "n")
	g.Link("n", "o")
	g.Link("o", "b")
	g.Link("q", "p")
	g.Link("p", "r")
	g.Link("p", "s")
	g.Link("p", "t")
	g.Link("p", "u")
	g.Link("p", "v")
	g.Link("v", "b")
	g.Link("r", "b")
	g.Link("s", "b")
	g.Link("t", "b")
	g.Link("u", "b")
	g.Link("q", "x")
	g.Link("x", "y")
	g.Link("x", "z")
	g.Link("x", "4")
	g.Link("x", "5")
	g.Link("x", "6")
	g.Link("6", "b")
	g.Link("y", "b")
	g.Link("z", "b")
	g.Link("4", "b")
	g.Link("5", "b")

	err := g.Evaluate()
	if err != nil {
		t.Error(err)
	}

	testRun := []string{}
	f := func(name string) error {
		testRun = append(testRun, name)
		return nil
	}
	err = g.CompileRun(f)
	if err != nil {
		t.Error(err)
	}

	exp := "b"
	got := testRun[len(testRun)-1]
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestReRunWithAdd(t *testing.T) {
	var exp, got string

	g := graph.New()
	g.Add("a")
	g.Add("b")
	g.Add("c")
	g.Link("a", "b")
	g.Link("b", "c")
	err := g.Evaluate()
	if err != nil {
		t.Error(err)
	}

	exp = "a"
	got = ""
	_, fnd := g.StartNodes()[exp]
	if !fnd {
		t.Fatalf(`expected "%s"`, exp)
	}
	g.SetRun(func(s string) error { got += s; return nil }, exp)
	exp = "abc"
	if exp != got {
		t.Fatalf(`expected "%s", got "%s"`, exp, got)
	}

	g.Add("d")
	g.Link("d", "c")
	err = g.Evaluate()
	if err != nil {
		t.Error(err)
	}

	exp = "d"
	got = ""
	_, fnd = g.StartNodes()[exp]
	if !fnd {
		t.Fatalf(`expected "%s"`, exp)
	}
	g.SetRun(func(s string) error { got += s; return nil }, exp)
	exp = "dc"
	if exp != got {
		t.Fatalf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestReRunWithUnlink(t *testing.T) {
	var exp, got string

	g := graph.New()
	g.Add("a")
	g.Add("b")
	g.Add("c")
	g.Link("a", "b")
	g.Link("b", "c")
	err := g.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	exp = "a"
	got = ""
	_, fnd := g.StartNodes()[exp]
	if !fnd {
		t.Fatalf(`expected "%s"`, exp)
	}
	g.SetRun(func(s string) error { got += s; return nil }, exp)
	exp = "abc"
	if exp != got {
		t.Fatalf(`expected "%s", got "%s"`, exp, got)
	}

	g.Unlink("a", "b")
	g.Link("a", "c")
	err = g.Evaluate()
	if err != nil {
		t.Fatal(err)
	}

	exp = "a"
	got = ""
	_, fnd = g.StartNodes()[exp]
	if !fnd {
		t.Fatalf(`expected "%s"`, exp)
	}
	g.SetRun(func(s string) error { got += s; return nil }, exp)
	exp = "ac"
	if exp != got {
		t.Fatalf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "b"
	got = ""
	_, fnd = g.StartNodes()[exp]
	if !fnd {
		t.Fatalf(`expected "%s"`, exp)
	}
	g.SetRun(func(s string) error { got += s; return nil }, exp)
	exp = "bc"
	if exp != got {
		t.Fatalf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestProfile(t *testing.T) {
	t.Skip()
	ns := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMN"
	g := graph.New()
	for i := 0; i < len(ns); i++ {
		g.Add(string(ns[i]))
	}
	for i := 0; i < len(ns)-2; i++ {
		for j := 0; j < len(ns)-i-2; j++ {
			g.Link(string(ns[i]), string(ns[i+1+rand.Intn(len(ns)-i-1)]))
		}
	}
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	g.Evaluate()
}

func Benchmark10N8E(b *testing.B) {
	ns := "abcdefghij"
	benchmarkEvaluate(ns, b)
}

func Benchmark20N98E(b *testing.B) {
	ns := "abcdefghijklmnopqrst"
	benchmarkEvaluate(ns, b)
}

func Benchmark30N288E(b *testing.B) {
	ns := "abcdefghijklmnopqrstuvwxyzABCD"
	benchmarkEvaluate(ns, b)
}

func Benchmark40N578E(b *testing.B) {
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
		for i := 0; i < len(ns)-6; i++ {
			for j := 0; j < len(ns)-i-6; j++ {
				g.Link(string(ns[i]), string(ns[i+1+rand.Intn(len(ns)-i-1)]))
			}
		}

		g.Evaluate()

	}
}
