// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package graph

type DFSFunc func(Vertex, func() error) error

func (g *Graph) DFS(start Vertex, cb DFSFunc) error {
	return g.dfs(cb, map[interface{}]struct{}{}, hashcode(start))
}

func (g *Graph) dfs(cb DFSFunc, visited map[interface{}]struct{}, v interface{}) error {
	/*
	   procedure DFS(G, v) is
	       label v as discovered
	       for all directed edges from v to w that are in G.adjacentEdges(v) do
	           if vertex w is not labeled as discovered then
	               recursively call DFS(G, w)
	*/

	// Make our map for visited
	visited[v] = struct{}{}

	// for all directed edges from v to w that are in G.adjacentEdges(v) do
	for w := range g.adjacencyOut[v] {
		// if vertex w is not labeled as discovered then
		if _, ok := visited[w]; !ok {
			// call our callback
			if err := cb(g.hash[w], func() error {
				// recursively call DFS(G, w)
				return g.dfs(cb, visited, w)
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
