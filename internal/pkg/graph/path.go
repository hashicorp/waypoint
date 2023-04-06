// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package graph

// TopoShortestPath returns the shortest path information given the
// topological sort of the graph L. L can be retrieved using any topological
// sort algorithm such as KahnSort.
//
// The return value are two maps with the distance to and edge to information,
// respectively. distTo maps the total distance from source to the given
// vertex. edgeTo maps the previous edge to get to a vertex from source.
func (g *Graph) TopoShortestPath(L TopoOrder) (distTo map[interface{}]int, edgeTo map[interface{}]Vertex) {
	/*
	   Set the distance to the source to 0;
	   Set the distances to all other vertices to infinity;
	   For each vertex u in L
	      - Walk through all neighbors v of u;
	      - If dist(v) > dist(u) + w(u, v)
	         - Set dist(v) <- dist(u) + w(u, v);
	*/

	// Set the distance to the source to 0;
	// Set the distances to all other vertices to infinity;
	// We don't actually set anything to "infinity" here since we can simulate
	// it by checking for existance in the map.
	distTo = map[interface{}]int{}
	edgeTo = map[interface{}]Vertex{}

	// For each vertex u in L
	for _, u := range L {
		uh := hashcode(u)

		// Walk through all neighbors v of u;
		for vh, weight := range g.adjacencyOut[uh] {
			// x = dist(u) + w(u, v)
			x := distTo[uh] + weight

			// If dist(v) > dist(u) + w(u, v)
			if _, ok := distTo[vh]; !ok || distTo[vh] > x {
				distTo[vh] = x
				edgeTo[vh] = u
			}
		}
	}

	return distTo, edgeTo
}

// EdgeToPath turns an "edge to" mapping into a vertex slice of the path
// from some target.
func (g *Graph) EdgeToPath(target Vertex, edgeTo map[interface{}]Vertex) []Vertex {
	var result []Vertex

	current := target
	for current != nil {
		result = append(result, current)
		current = edgeTo[hashcode(current)]
	}

	// Reverse it, since this puts the path in reverse order
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}

	return result
}
