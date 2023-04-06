// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package graph

import (
	"container/heap"
	"math"
)

// Dijkstra implements Dijkstra's algorithm for finding single source
// shortest paths in an edge-weighted graph with non-negative edge weights.
// The graph may have cycles.
func (g *Graph) Dijkstra(src Vertex) (distTo map[interface{}]int, edgeTo map[interface{}]Vertex) {
	srchash := hashcode(src)

	/*
	   for each vertex V in G
	       distance[V] <- infinite
	       previous[V] <- NULL
	*/
	queue := make(distQueue, 0, len(g.hash))
	queueItem := map[interface{}]*distQueueItem{}
	for k := range g.hash {
		item := &distQueueItem{
			v:        k,
			distance: math.MaxInt32,
			previous: nil,
			index:    len(queue),
		}
		queueItem[k] = item

		// If V != S, add V to Priority Queue Q
		queue = append(queue, item)
	}

	// distance[S] <- 0
	queueItem[srchash].distance = 0

	// Init the heap so we can use the queue
	heap.Init(&queue)

	// while Q IS NOT EMPTY
	visited := map[interface{}]struct{}{}
	for queue.Len() > 0 {
		// U <- Extract MIN from Q
		u := heap.Pop(&queue).(*distQueueItem)
		visited[u.v] = struct{}{}

		// for each unvisited neighbour V of U
		for vhash, weight := range g.adjacencyOut[u.v] {
			if _, ok := visited[vhash]; ok {
				continue
			}

			v := queueItem[vhash]

			// tempDistance <- distance[U] + edge_weight(U, V)
			tempDistance := u.distance + int32(weight)

			// if tempDistance < distance[V]
			if tempDistance < v.distance {
				// distance[V] <- tempDistance
				// previous[V] <- U
				v.distance = tempDistance
				v.previous = u.v
				heap.Fix(&queue, v.index)
			}
		}
	}

	// Return our distance and previous map
	distTo = make(map[interface{}]int, len(queueItem))
	edgeTo = make(map[interface{}]Vertex, len(queueItem))
	for _, item := range queueItem {
		distTo[item.v] = int(item.distance)
		edgeTo[item.v] = g.hash[item.previous]
	}

	return distTo, edgeTo
}

// distQueue is a priority queue implementation on top of a heap that
// is used by Dijkstra to keep track of state. heap.Pop on this queue
// will return the item with the minimal "distance" value.
type distQueue []*distQueueItem

type distQueueItem struct {
	v        interface{} // Vertex hashcode
	distance int32
	previous interface{} // Previous vertex hashcode
	index    int
}

func (pq distQueue) Len() int { return len(pq) }

func (pq distQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq distQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *distQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*distQueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *distQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

var _ heap.Interface = (*distQueue)(nil)
