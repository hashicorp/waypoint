// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package graph

import "fmt"

// KahnSort will return the topological sort of the graph using Kahn's algorithm.
// The graph must not have any cycles or this will panic.
func (g *Graph) KahnSort() TopoOrder {
	/*
	   L ← Empty list that will contain the sorted elements
	   S ← Set of all nodes with no incoming edge

	   while S is non-empty do
	       remove a node n from S
	       add n to tail of L
	       for each node m with an edge e from n to m do
	           remove edge e from the graph
	           if m has no other incoming edges then
	               insert m into S

	   if graph has edges then
	       return error   (graph has at least one cycle)
	   else
	       return L   (a topologically sorted order)
	*/

	// Copy the graph
	g = g.Copy()

	// L ← Empty list that will contain the sorted elements
	L := make([]Vertex, 0, len(g.adjacencyOut))

	// S ← Set of all nodes with no incoming edge
	S := []interface{}{}
	for v, list := range g.adjacencyIn {
		if len(list) == 0 {
			S = append(S, v)
		}
	}

	// while S is non-empty do
	for len(S) > 0 {
		// remove a node n from S
		n := S[len(S)-1]
		S = S[:len(S)-1]

		// add n to tail of L
		L = append(L, g.hash[n])

		// for each node m with an edge e from n to m do
		for m := range g.adjacencyOut[n] {
			// remove edge e from the graph
			g.RemoveEdge(n, m)

			// if m has no other incoming edges then
			if len(g.adjacencyIn[m]) == 0 {
				// insert m into S
				S = append(S, m)
			}
		}
	}

	// if graph has edges then
	//   return error   (graph has at least one cycle)
	for _, out := range g.adjacencyOut {
		if len(out) > 0 {
			// We have cycles, so let's do cycle detection to give a better
			// error message.
			cycles := g.Cycles()
			panic(fmt.Sprintf("graph has cycles: %v", cycles))
		}
	}

	return L
}

// TopoOrder is a topological ordering.
type TopoOrder []Vertex

// At returns a new TopoOrder that starts at the given vertex.
// This returns a slice into the ordering so it is not safe to modify.
func (t TopoOrder) At(v Vertex) TopoOrder {
	for i := 0; i < len(t); i++ {
		if t[i] == v {
			return t[i:]
		}
	}
	return nil
}

// Until returns a new TopoOrder that ends at (and includes) the given
// vertex. This returns a slice into the ordering so it is not safe to modify.
func (t TopoOrder) Until(v Vertex) TopoOrder {
	for i := 0; i < len(t); i++ {
		if t[i] == v {
			return t[:i]
		}
	}
	return nil
}
