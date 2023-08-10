// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package graph

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransitiveReduction(t *testing.T) {
	require := require.New(t)

	var g Graph

	g.Add("A")
	g.Add("B")
	g.Add("C")
	g.Add("D")
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")
	g.AddEdge("C", "D")
	g.AddEdge("A", "D") // we expect this edge to be removed

	g.TransitiveReduction()

	require.Equal([]Vertex{"B"}, g.OutEdges("A"))
	require.Equal([]Vertex{"C"}, g.OutEdges("B"))
	require.Equal([]Vertex{"D"}, g.OutEdges("C"))
}
