![](https://github.com/oligoden/math-graph/workflows/Build/badge.svg)

# Graph Theory Maths Library

This Go library allows you to create graphs, evaluate them and
run functions on the nodes.

## Install

```bash
go get github.com/oligoden/math-graph
```

## Usage

First initialize a graph, add nodes add connections(links) and evaluate the graph.

```golang
g := graph.New()

g.Add("a")
g.Add("b")
g.Add("c")
g.Add("d")
g.Add("e")

g.Link("a", "c")
g.Link("b", "c")
g.Link("b", "d")
g.Link("c", "d")
g.Link("d", "e")

g.Evaluate()
```

Now you can get the start nodes with

```golang
g.StartNodes()
```

or run your own function with
```golang
g.Run(func(node string){
    //code
})
```

on each node down the graph.