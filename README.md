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
Benchmark10N32E-8     	     898	   1187336 ns/op	   52580 B/op	     251 allocs/op
Benchmark20N162E-8    	      56	  18251628 ns/op	  483484 B/op	    1329 allocs/op
Benchmark30N392E-8    	      14	  95632523 ns/op	 2282834 B/op	    3922 allocs/op
Benchmark40N722E-8    	       4	 296888227 ns/op	 4202094 B/op	    7464 allocs/op
Benchmark50N1152E-8   	       2	 684910964 ns/op	 6437896 B/op	   12330 allocs/op
Benchmark60N1682E-8   	       1	1813735357 ns/op	23594432 B/op	   28260 allocs/op
```