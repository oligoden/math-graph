package graph_test

import (
	"fmt"

	graph "github.com/oligoden/math-graph"
)

func ExampleUsage() {
	var exp, got string

	g := graph.New()

	g.Add("a")
	g.Add("b")
	g.Add("c")
	g.Add("d")
	g.Add("e")

	g.Link("a", "b")
	g.Link("b", "c")
	g.Link("d", "c")

	err := g.Evaluate()
	if err != nil {
		fmt.Println(err)
		return
	}

	exp = "a"
	_, fnd := g.StartNodes()[exp]
	if !fnd {
		fmt.Printf(`expected to find "%s"`, exp)
		return
	}

	g.SetRun(func(s string) error { got += s; return nil }, exp)
	exp = "abc"
	if exp != got {
		fmt.Printf(`expected "%s", got "%s"`, exp, got)
		return
	}
	fmt.Println(exp)

	//Output:
	//abc
}
