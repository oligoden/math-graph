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
g.SetRun(func(node string){
    //code
}, "a")
```

on each node down the graph.

## Benchmarks

```none
Benchmark10N8E-8      	   38554	     29688 ns/op	    9321 B/op	     103 allocs/op
Benchmark20N98E-8     	    3609	    335295 ns/op	   51123 B/op	     335 allocs/op
Benchmark30N288E-8    	     772	   1599971 ns/op	  163803 B/op	     714 allocs/op
Benchmark40N578E-8    	     240	   4923217 ns/op	  372283 B/op	    1233 allocs/op
Benchmark50N1152E-8   	     100	  12172180 ns/op	  729566 B/op	    1877 allocs/op
Benchmark60N1682E-8   	      44	  25401861 ns/op	 1240786 B/op	    2718 allocs/op
```