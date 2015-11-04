// Yen's algorithm -- a way of finding the K shortest paths assuming an efficient
// algorithm exists for finding a shortest path in general.

package multipath

import (
	"sort"
	"fmt"
	"strings"
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

type metricPathSlice []metricPath;
func (m metricPathSlice) Len() int {return len(m);}
func (m metricPathSlice) Less(i, j int) bool {return m[i].metric < m[j].metric;}
func (m metricPathSlice) Swap(i, j int) {
	temp := m[i];
	m[i] = m[j];
	m[j] = temp;
}
func (m *metricPathSlice) AppendUnique(p metricPath, solutions [][]string) {
	fmt.Printf("Appending %v to MPS %v with existing solutions (%v)\n", p, *m, solutions);
	found := false;
	for _, elem := range(*m) {
		if pathsEqual(p.path, elem.path) {
			fmt.Printf("...rejected (present in MPS)\n");
			found = true;
			break;
		}
	}
	for _, elem := range(solutions) {
		if pathsEqual(p.path, elem) {
			fmt.Printf("...rejected (present in solution)\n");
			found = true;
			break;
		}
	}
	if !found {
		fmt.Printf("...accepted (not present)\n");
		*m = append(*m, p);
	}
}

func pathsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false;
	}
	for idx, ia := range(a) {
		if ia != b[idx] {
			return false;
		}
	}
	return true;
}

func copyGraph(grph graph.Graph) graph.Graph {
	ret := graph.NewDefaultGraph();
	nodes := grph.GetVertices()
	for node, _ := range(nodes) {
		ret.AddVertex(node);
	}
	for node, _ := range(nodes) {
		children, _ := grph.GetChildren(node);
		for child, _ := range(children) {
			weight, _ := grph.GetWeight(node, child);
			ret.AddEdge(node, child, weight);
		}
	}
	return ret;
}

func ToDot(grph graph.Graph, directed bool) string {
	lines := make([]string, 0);
	var sep string;
	if directed {
		lines = append(lines, "digraph {");
		sep = "->";
	} else {
		lines = append(lines, "strict graph {");
		sep = "--";
	}
	nodes := grph.GetVertices();
	for node, _ := range(nodes) {
		lines = append(lines, fmt.Sprintf("\"%s\"", node));
		children, _ := grph.GetChildren(node);
		for child, _ := range(children) {
			weight, _ := grph.GetWeight(node, child);
			lines = append(lines, fmt.Sprintf("\"%s\" %s \"%s\" [label=\"%f\"]", node, sep, child, weight));
		}
	}
	lines = append(lines, "}");
	return strings.Join(lines, "\n");
}

func Yen(_grph graph.Graph, source, sink string, K int) [][]string {
	var err error;
	var firstPath []string;
	A := make([][]string, 0, K);
	firstPath, _, err = graph.Dijkstra(_grph, source, sink);
	A = append(A, append([]string(nil), firstPath...));
	if err != nil {
		return nil;
	}
	B := make(metricPathSlice, 0);
	removedNodes := make([]string, 0);
	removedEdges := make([]nodePair, 0);
	for k := 1; k < K; k++ {
		fmt.Printf("Current shortest paths: %s\n", A);
		for i := 0; i < len(A[k - 1]) - 1; i++ {
			grph := copyGraph(_grph);
			spurNode := A[k - 1][i];
			rootPath := A[k - 1][:i + 1];
			fmt.Printf("spurNode: %s\nrootPath: %s\n", spurNode, rootPath);

			for _, path := range(A) {
				if path == nil {
					continue;
				}
				if i + 1 < len(path) && len(path) > len(rootPath) && pathsEqual(rootPath, path[:i + 1]) {
					fmt.Printf("Cutting path part at prefix: %s == %s, disconnecting %s --> %s\n", rootPath, path[:i + 1], path[i], path[i + 1]);
					weight, _ := grph.GetWeight(path[i], path[i + 1]);
					grph.DeleteEdge(path[i], path[i + 1])
					removedEdges = append(removedEdges, nodePair{parent: path[i], child: path[i + 1], weight: weight});
				}
			}

			for _, node := range(rootPath) {
				if node != spurNode {
					fmt.Printf("Removing root node %s\n", node);
					neighbors, _ := grph.GetParents(node);
					for neighbor, _ := range(neighbors) {
						weight, _ := grph.GetWeight(neighbor, node);
						removedEdges = append(removedEdges, nodePair{parent: neighbor, child: node, weight: weight});
					}
					neighbors, _ = grph.GetChildren(node);
					for neighbor, _ := range(neighbors) {
						weight, _ := grph.GetWeight(node, neighbor);
						removedEdges = append(removedEdges, nodePair{parent: node, child: neighbor, weight: weight});
					}
					removedNodes = append(removedNodes, node);
					grph.DeleteVertex(node);
				}
			}

			spurPath, spurCosts, err := graph.Dijkstra(grph, spurNode, sink);
			fmt.Printf("Path from %s to %s: %s\n", spurNode, sink, spurPath);
			if err == nil {
				totalPath := append(rootPath[:len(rootPath) - 1], append([]string(nil), spurPath...)...);
				if totalPath[len(totalPath) - 1] != sink {
					fmt.Printf("WARNING: Rejecting (path doesn't end at sink)\n");
				} else {
					B.AppendUnique(metricPath{path: totalPath, metric: spurCosts[sink]}, A);
				}
			}

			for _, node := range(removedNodes) {
				// fmt.Printf("-- Node added: %s\n", node);
				grph.AddVertex(node);
			}
			for _, edge := range(removedEdges) {
				// fmt.Printf("-- Edge added: %s --- %f --> %s\n", edge.parent, edge.weight, edge.child);
				grph.AddEdge(edge.parent, edge.child, edge.weight);
			}
			removedNodes = nil;
			removedEdges = nil;
		}
		if len(B) == 0 {
			break;
		}
		sort.Sort(B);
		fmt.Printf("Considered shortest paths: %v\n", B);
		A = append(A, B[0].path);
		B = B[1:];
	}
	return A;
}
