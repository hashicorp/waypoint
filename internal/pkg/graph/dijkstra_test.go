// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDijkstra(t *testing.T) {
	t.Run("typical", func(t *testing.T) {
		var g Graph
		g.Add("B")
		g.Add("C")
		g.Add("D")
		g.Add("E")
		g.Add("S")
		g.Add("T")
		g.AddEdgeWeighted("S", "B", 4)
		g.AddEdgeWeighted("S", "C", 2)
		g.AddEdgeWeighted("B", "C", 1)
		g.AddEdgeWeighted("B", "D", 5)
		g.AddEdgeWeighted("C", "D", 8)
		g.AddEdgeWeighted("C", "E", 10)
		g.AddEdgeWeighted("D", "E", 2)
		g.AddEdgeWeighted("D", "T", 6)
		g.AddEdgeWeighted("E", "T", 2)

		distTo, edgeTo := g.Dijkstra("S")

		path := []interface{}{"T"}
		for next := edgeTo[path[0]]; next != nil; next = edgeTo[next] {
			path = append(path, next)
		}

		require.Equal(t, 13, distTo["T"])
		require.Equal(t, []interface{}{"T", "E", "D", "B", "S"}, path)
	})

	t.Run("cycle", func(t *testing.T) {
		var g Graph
		g.Add("A")
		g.Add("B")
		g.Add("C")
		g.Add("D")
		g.Add("E")
		g.Add("F")
		g.AddEdgeWeighted("A", "B", 10)
		g.AddEdgeWeighted("B", "C", 1)
		g.AddEdgeWeighted("C", "E", 3)
		g.AddEdgeWeighted("E", "D", 10)
		g.AddEdgeWeighted("D", "B", 10)
		g.AddEdgeWeighted("E", "F", 22)

		distTo, edgeTo := g.Dijkstra("A")

		path := []interface{}{"F"}
		for next := edgeTo[path[0]]; next != nil; next = edgeTo[next] {
			path = append(path, next)
		}

		require.Equal(t, 36, distTo["F"])
		require.Equal(t, []interface{}{"F", "E", "C", "B", "A"}, path)
	})
}
