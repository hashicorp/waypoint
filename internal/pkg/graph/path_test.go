// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGraphTopoShortestPath(t *testing.T) {
	require := require.New(t)

	// This test is from the coursera page linked below. I chose to copy
	// something relatively well known and externally solved so I can be 100%
	// sure I got it right, plus to test cases where paths change multiple times.
	// https://www.coursera.org/lecture/algorithms-part2/edge-weighted-dags-6rxSt
	var g Graph
	g.Add(0)
	g.Add(1)
	g.Add(2)
	g.Add(3)
	g.Add(4)
	g.Add(5)
	g.Add(6)
	g.Add(7)
	g.AddEdgeWeighted(0, 1, 5)
	g.AddEdgeWeighted(0, 7, 8)
	g.AddEdgeWeighted(0, 4, 9)
	g.AddEdgeWeighted(1, 7, 4)
	g.AddEdgeWeighted(1, 3, 15)
	g.AddEdgeWeighted(1, 2, 12)
	g.AddEdgeWeighted(7, 2, 7)
	g.AddEdgeWeighted(7, 5, 6)
	g.AddEdgeWeighted(4, 7, 5)
	g.AddEdgeWeighted(4, 5, 4)
	g.AddEdgeWeighted(4, 6, 20)
	g.AddEdgeWeighted(5, 2, 1)
	g.AddEdgeWeighted(5, 6, 13)
	g.AddEdgeWeighted(2, 3, 3)
	g.AddEdgeWeighted(2, 6, 11)
	g.AddEdgeWeighted(3, 6, 9)

	sort := g.KahnSort()
	t.Logf("sorted: %#v", sort)

	distTo, edgeTo := g.TopoShortestPath(sort)
	require.Equal(0, distTo[0])
	require.Equal(5, distTo[1])
	require.Equal(14, distTo[2])
	require.Equal(17, distTo[3])
	require.Equal(9, distTo[4])
	require.Equal(13, distTo[5])
	require.Equal(25, distTo[6])
	require.Equal(8, distTo[7])

	require.Equal(nil, edgeTo[0])
	require.Equal(0, edgeTo[1])
	require.Equal(5, edgeTo[2])
	require.Equal(2, edgeTo[3])
	require.Equal(0, edgeTo[4])
	require.Equal(4, edgeTo[5])
	require.Equal(2, edgeTo[6])
	require.Equal(0, edgeTo[7])
}
