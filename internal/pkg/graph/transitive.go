// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package graph

// TransitiveReduction performs a transitive reduction on the graph, effectively
// removing any “shortcuts” from the graph. It performs the reduction *in-place*
// so please call Copy to avoid mutating the original.
func (g *Graph) TransitiveReduction() *Graph {
	// Recursive walk function for use later on
	var walk func(i Vertex, k Vertex, d int) (int, bool)
	walk = func(i Vertex, k Vertex, d int) (int, bool) {
		result := d
		found := false

		for _, j := range g.OutEdges(i) {
			if hashcode(j) == hashcode(k) && d >= result {
				result = d + 1
				found = true
				continue
			}
			if dd, ok := walk(j, k, d+1); ok && dd > result {
				result = dd
				found = true
				continue
			}
		}
		return result, found
	}

	// Build longest-path matrix
	depths := make(map[Vertex]map[Vertex]int)
	for _, i := range g.Vertices() {
		depths[i] = make(map[Vertex]int)

		for _, j := range g.Vertices() {
			if i == j {
				continue
			}

			if d, ok := walk(i, j, 0); ok {
				depths[i][j] = d
			}
		}
	}

	// Remove shortcuts
	for _, i := range g.Vertices() {
		for _, j := range g.OutEdges(i) {
			if depths[i][j] > 1 {
				g.RemoveEdge(i, j)
			}
		}
	}

	return g
}
