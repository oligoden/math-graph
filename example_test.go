package graph_test

import (
	"fmt"

	graph "github.com/oligoden/math-graph"
)

func ExampleSetRun() {
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

	var refs string
	g.SetRun(func(s string) error { refs += s; return nil }, "a")
	fmt.Println(refs)

	//Output:
	//abc
}
