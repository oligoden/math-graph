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
Benchmark10N32E-8     	   28482	     42399 ns/op	   14461 B/op	     151 allocs/op
Benchmark20N162E-8    	    3002	    365577 ns/op	   70443 B/op	     442 allocs/op
Benchmark30N392E-8    	     740	   1517831 ns/op	  207188 B/op	     879 allocs/op
Benchmark40N722E-8    	     262	   5220934 ns/op	  453054 B/op	    1473 allocs/op
Benchmark50N1152E-8   	     108	  10728966 ns/op	  862323 B/op	    2181 allocs/op
Benchmark60N1682E-8   	      61	  21751799 ns/op	 1404854 B/op	    3048 allocs/op
```