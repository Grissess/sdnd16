// Yen's algorithm -- a way of finding the K shortest paths assuming an efficient
// algorithm exists for finding a shortest path in general.

package multipath

import (
	"sort"
	"github.com/gyuho/goraph/graph"
)

type nodePair struct {
	parent, child string
	weight float64
}

type metricPath struct {
	path []string
	metric float64
}

func (m []metricPath) Len() int {return len(m);}
func (m []metricPath) Less(i, j int) bool {return m[i].metric < m[j].metric;}
func (m []metricPath) Swap(i, j int) {
	temp := m[i];
	m[i] = m[j];
	m[j] = temp;
}


func Yen(grph graph.DefaultGraph, source, sink string, K int) [][]string {
	var err error;
	A := make([][]string, 0, K);
	A[0], _, err = graph.Dijkstra(grph, source, sink);
	if err != nil {
		return nil;
	}
	B := make([]metricPath);
	removedNodes := make([]string);
	removedEdges := make([]nodePair);
	for k := 1; k < K; k++ {
		for i := 0; i < len(A[k - 1]); i++ {
			spurNode := A[k - 1][i];
			rootPath := A[k - 1][:i];
			for _, path := range(A) {
				if path == nil {
					continue;
				}
				if rootPath == path[:i] {
					weight, _ := grph.GetWeight(path[1], path[i+1]);
					grph.DeleteEdge(path[i], path[i+1])
					removedEdges = append(removedEdges, nodePair{parent: path[i], child: path[i+1], weight: weight});
				}
			}
			for _, node := range(rootPath) {
				if node != spurNode {
					neighbors, _ := graph.GetParents(node);
					for _, neighbor := range(neighbors) {
						weight, _ := grph.GetWeight(neighbor, node);
						removedEdges = append(removedEdges, nodePair{parent: neighbor, child: node, weight: weight});
					}
					neighbors, _ = graph.GetChildren(node);
					for _, neighbor := range(neighbors) {
						weight, _ := grph.GetWeight(node, neighbor);
						removedEdges = append(removedEdges, nodePair{parent: node, child: neighbor, weight: weight});
					}
					grph.DeleteVertex(node);
				}
			}

			spurPath, spurCosts, err := graph.Dijkstra(grph, spurNode, sink);
			if err != nil {
				continue;
			}
			totalPath := append(rootPath, spurPath...);
			B = append(B, metricPath{path: totalPath, metric: spurCosts[sink]});

			for _, node := range(removedNodes) {
				grph.AddVertex(node);
			}
			for _, edge := range(removedEdges) {
				grph.AddEdge(edge.parent, edge.child, edge.weight);
			}
		}
		if len(B) == 0 {
			break;
		}
		sort.Sort(B);
		A[k] = B[0];
		B = B[1:];
	}
	return A;
}
