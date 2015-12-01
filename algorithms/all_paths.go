// All-paths implementation adapter

package algorithms

import (
	"fmt"
	"strings"
	"github.com/gonum/graph"
	"github.com/gonum/graph/path"
)

// Structure to represent a single path on the graph from paths[0] to paths[len(paths)-1]
type Path struct {
	path []graph.Node
	weight float64
}

// Returns the canonical representation of a path given a map from node IDs to string labels
func (p Path) PathString(labels map[int]string) string {
	pathlen := len(p.path);
	names := make([]string, pathlen + 2);
	for idx, node := range(p.path) {
		names[idx] = labels[node.ID()];
	}
	names[pathlen] = "|";
	names[pathlen + 1] = fmt.Sprint(p.weight);
	return strings.Join(names, " ");
}

// Converts a path.AllShortest over a graph.Graph to a map of node pairs, such that
// every pair of graph.Node in the graph maps to a Path.
func ConvertAllPaths(g graph.Graph, ap path.AllShortest) map[int]map[int]Path {
	nodes := g.Nodes();
	ret := make(map[int]map[int]Path);

	for _, start := range(nodes) {
		pathmap :=  make(map[int]Path);
		ret[start.ID()] = pathmap;
		for _, end := range(nodes) {
			if start == end {
				continue;
			}
			path, weight, _ := ap.Between(start, end);
			pathmap[end.ID()] = Path{path: path, weight: weight};
		}
	}

	return ret;
}
