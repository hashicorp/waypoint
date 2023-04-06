// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package graph

// Cycles returns all the detected cycles. This may not be fully exhaustive
// since we use Tarjan's algoritm for strongly connected components to detect
// cycles and this isn't guaranteed to find all cycles.
func (g *Graph) Cycles() [][]Vertex {
	var cycles [][]Vertex
	for _, cycle := range g.StronglyConnected() {
		if len(cycle) > 1 {
			cycles = append(cycles, cycle)
		}
	}

	return cycles
}

// StronglyConnected returns the list of strongly connected components
// within the Graph g.
func (g *Graph) StronglyConnected() [][]Vertex {
	vs := g.Vertices()
	acct := sccAcct{
		NextIndex:   1,
		VertexIndex: make(map[Vertex]int, len(vs)),
	}
	for _, v := range vs {
		// Recurse on any non-visited nodes
		if acct.VertexIndex[v] == 0 {
			stronglyConnected(&acct, g, v)
		}
	}
	return acct.SCC
}

func stronglyConnected(acct *sccAcct, g *Graph, v Vertex) int {
	// Initial vertex visit
	index := acct.visit(v)
	minIdx := index

	for _, target := range g.OutEdges(v) {
		targetIdx := acct.VertexIndex[target]

		// Recurse on successor if not yet visited
		if targetIdx == 0 {
			minIdx = min(minIdx, stronglyConnected(acct, g, target))
		} else if acct.inStack(target) {
			// Check if the vertex is in the stack
			minIdx = min(minIdx, targetIdx)
		}
	}

	// Pop the strongly connected components off the stack if
	// this is a root vertex
	if index == minIdx {
		var scc []Vertex
		for {
			v2 := acct.pop()
			scc = append(scc, v2)
			if v2 == v {
				break
			}
		}

		acct.SCC = append(acct.SCC, scc)
	}

	return minIdx
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// sccAcct is used ot pass around accounting information for
// the StronglyConnectedComponents algorithm
type sccAcct struct {
	NextIndex   int
	VertexIndex map[Vertex]int
	Stack       []Vertex
	SCC         [][]Vertex
}

// visit assigns an index and pushes a vertex onto the stack
func (s *sccAcct) visit(v Vertex) int {
	idx := s.NextIndex
	s.VertexIndex[v] = idx
	s.NextIndex++
	s.push(v)
	return idx
}

// push adds a vertex to the stack
func (s *sccAcct) push(n Vertex) {
	s.Stack = append(s.Stack, n)
}

// pop removes a vertex from the stack
func (s *sccAcct) pop() Vertex {
	n := len(s.Stack)
	if n == 0 {
		return nil
	}
	vertex := s.Stack[n-1]
	s.Stack = s.Stack[:n-1]
	return vertex
}

// inStack checks if a vertex is in the stack
func (s *sccAcct) inStack(needle Vertex) bool {
	for _, n := range s.Stack {
		if n == needle {
			return true
		}
	}
	return false
}
