// Simple implementions of convenient algorithms on graphs

package algorithms

import (
	"fmt"
	"math"
	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
	"github.com/eapache/queue"
)

// Builds a BFS tree (as a directed graph) from the given graph and start node.
func BFSTree(g graph.Graph, start graph.Node) *simple.DirectedGraph {
	if !g.Has(start) {
		panic(fmt.Sprintf("BFSTree: Start node %r not in graph %r", start, g));
	}

	ret := simple.NewDirectedGraph(0.0, math.Inf(1));
	seen := make(map[int]bool);
	q := queue.New();
	q.Add(start);
	ret.AddNode(simple.Node(start.ID()));

	for q.Length() > 0 {
		node := q.Peek().(graph.Node);
		q.Remove();
		for _, neighbor := range(g.From(node)) {
			if !seen[neighbor.ID()] {
				seen[neighbor.ID()] = true;
				ret.AddNode(simple.Node(neighbor.ID()))
				ret.SetEdge(simple.Edge {F: simple.Node(node.ID()), T: simple.Node(neighbor.ID()), W: g.Edge(node, neighbor).Weight()});
				q.Add(neighbor);
			}
		}
	}

	return ret;
}
