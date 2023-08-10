// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package graph

import (
	"bytes"
	"fmt"
	"sort"
)

// Graph represents a graph structure.
//
// Unless otherwise documented, it is unsafe to call any method on Graph concurrently.
type Graph struct {
	// adjacency represents graphs using an adjaency list. Vertices are
	// represented using their hash codes for simpler equaliy checks.
	adjacencyOut map[interface{}]map[interface{}]int
	adjacencyIn  map[interface{}]map[interface{}]int

	// hash maintains the mapping of hash codes to the representative Vertex.
	// It is assumed that two identical hashcodes of v1 and v2 are semantically
	// the same Vertex even if v1 != v2 in Go.
	hash map[interface{}]Vertex
}

// Add adds a vertex to the graph. If a vertex with the same identity exists
// this will overwrite that vertex.
func (g *Graph) Add(v Vertex) Vertex {
	g.init()
	h := hashcode(v)
	if _, ok := g.adjacencyOut[h]; !ok {
		g.adjacencyOut[h] = make(map[interface{}]int)
		g.adjacencyIn[h] = make(map[interface{}]int)
		g.hash[h] = v
	}
	return v
}

// AddOverwrite is the same as Add, except that even if the vertex
// already exists with the same hashcode, the pointed to value is replaced
// with the given v. This allows two Vertex values with the same hashcode
// but different values to be replaced.
func (g *Graph) AddOverwrite(v Vertex) Vertex {
	g.init()
	h := hashcode(v)
	g.hash[h] = v
	if _, ok := g.adjacencyOut[h]; !ok {
		g.adjacencyOut[h] = make(map[interface{}]int)
		g.adjacencyIn[h] = make(map[interface{}]int)
	}
	return v
}

// Remove removes the given vertex from the graph.
func (g *Graph) Remove(v Vertex) Vertex {
	// Note we don't need to call init here because delete() operations
	// are all safe on nil maps.
	h := hashcode(v)

	// First, delete all our out-edges by deleting both the
	// main key as well as any of the other nodes that are tracking the
	// in edge
	for out := range g.adjacencyOut[h] {
		delete(g.adjacencyIn[out], h)
	}
	delete(g.adjacencyOut, h)

	// Same as above but for in edges
	for in := range g.adjacencyIn[h] {
		delete(g.adjacencyOut[in], h)
	}
	delete(g.adjacencyIn, h)

	// Forget this node completely
	delete(g.hash, h)
	return v
}

// Vertex returns the vertex by id. This can be done to get the node that
// is actually in the graph. This will return nil if the given vertex
// is not in the graph any longer.
func (g *Graph) Vertex(id interface{}) Vertex {
	g.init()
	return g.hash[id]
}

// Vertices returns the list of all the vertices in this graph.
func (g *Graph) Vertices() []Vertex {
	result := make([]Vertex, 0, len(g.hash))
	for _, v := range g.hash {
		result = append(result, v)
	}

	return result
}

// AddEdge adds a directed edge to the graph from v1 to v2. Both v1 and v2
// must already be in the Graph via Add or this will do nothing.
func (g *Graph) AddEdge(v1, v2 Vertex) {
	g.AddEdgeWeighted(v1, v2, 1)
}

// AddEdgeWeighted adds a weighted edge. This is the same as AddEdge but
// with the specified weight. This will overwrite any existing edges.
func (g *Graph) AddEdgeWeighted(v1, v2 Vertex, weight int) {
	g.init()
	h1, h2 := hashcode(v1), hashcode(v2)
	g.adjacencyOut[h1][h2] = weight
	g.adjacencyIn[h2][h1] = weight
}

func (g *Graph) RemoveEdge(v1, v2 Vertex) {
	g.init()
	h1, h2 := hashcode(v1), hashcode(v2)
	delete(g.adjacencyOut[h1], h2)
	delete(g.adjacencyIn[h2], h1)
}

func (g *Graph) OutEdges(v Vertex) []Vertex {
	edges := g.adjacencyOut[hashcode(v)]
	if len(edges) == 0 {
		return nil
	}

	result := make([]Vertex, 0, len(edges))
	for h := range edges {
		result = append(result, g.hash[h])
	}

	return result
}

func (g *Graph) InEdges(v Vertex) []Vertex {
	edges := g.adjacencyIn[hashcode(v)]
	if len(edges) == 0 {
		return nil
	}

	result := make([]Vertex, 0, len(edges))
	for h := range edges {
		result = append(result, g.hash[h])
	}

	return result
}

// Reverse reverses the graph but _does not make a copy_. Any changes to
// this graph will impact the original Graph. You must call Copy on the
// result if you want to have a copy.
func (g *Graph) Reverse() *Graph {
	return &Graph{
		adjacencyOut: g.adjacencyIn,
		adjacencyIn:  g.adjacencyOut,
		hash:         g.hash,
	}
}

// Copy copies the graph. In the copy, any added or removed edges do not
// affect the original graph. The vertices themselves are not deep copied.
func (g *Graph) Copy() *Graph {
	var g2 Graph
	g2.init()

	for k, set := range g.adjacencyOut {
		copy := make(map[interface{}]int)
		for k, v := range set {
			copy[k] = v
		}
		g2.adjacencyOut[k] = copy
	}
	for k, set := range g.adjacencyIn {
		copy := make(map[interface{}]int)
		for k, v := range set {
			copy[k] = v
		}
		g2.adjacencyIn[k] = copy
	}
	for k, v := range g.hash {
		g2.hash[k] = v
	}

	return &g2
}

// String outputs some human-friendly output for the graph structure.
func (g *Graph) String() string {
	var buf bytes.Buffer
	buf.WriteString("\n")

	// Build the list of node names and a mapping so that we can more
	// easily alphabetize the output to remain deterministic.
	names := make([]string, 0, len(g.hash))
	mapping := make(map[string]Vertex, len(g.hash))
	for _, v := range g.hash {
		name := VertexName(v)
		names = append(names, name)
		mapping[name] = v
	}
	sort.Strings(names)

	// Write each node in order...
	for _, name := range names {
		v := mapping[name]
		targets := g.adjacencyOut[hashcode(v)]

		buf.WriteString(fmt.Sprintf("%s\n", name))

		// Alphabetize dependencies
		deps := make([]string, 0, len(targets))
		for targetHash, weight := range targets {
			deps = append(deps, fmt.Sprintf(
				"%s (%d)", VertexName(g.hash[targetHash]), weight))
		}
		sort.Strings(deps)

		// Write dependencies
		for _, d := range deps {
			buf.WriteString(fmt.Sprintf("  %s\n", d))
		}
	}

	return buf.String()
}

func (g *Graph) init() {
	if g.adjacencyOut == nil {
		g.adjacencyOut = make(map[interface{}]map[interface{}]int)
	}
	if g.adjacencyIn == nil {
		g.adjacencyIn = make(map[interface{}]map[interface{}]int)
	}
	if g.hash == nil {
		g.hash = make(map[interface{}]Vertex)
	}
}
