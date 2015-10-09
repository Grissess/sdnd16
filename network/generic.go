// Generic interface casting methods

package network;

import (
	"fmt"
	"strings"
)

func (self ErrNullPtr) Error() string {
	return fmt.Sprintf("Null pointer encounted in context: %v", string(self));
}

func (self ErrNoNode) Error() string {
	return fmt.Sprintf("No node %v in %v", self.label, self.graph);
}

func (self ErrAlreadyNode) Error() string {
	return fmt.Sprintf("Node %v already exists in %v", self.node, self.graph);
}

func (self ErrNoEdge) Error() string {
	return fmt.Sprintf("No edge from %v to %v in %v", self.src, self.dst, self.graph);
}

func (self ErrAlreadyEdge) Error() string {
	return fmt.Sprintf("Edge %v already exists in %v", self.edge, self.graph);
}

func (self ErrNotOwned) Error() string {
	return fmt.Sprintf("Object %v not owned by graph %v", self.member, self.graph);
}

func (self Label) String() string {
	return string(self);
}

func (self DsNode) String() string {
	return fmt.Sprintf("<Node %p '%s'>", &self, self.label);
}

func (self DsEdge) String() string {
	return fmt.Sprintf("%v -> %v { %v }", self.src, self.dst, self.attrs);
}

func (self DsGraph) String() string {
	strs := make([]string, len(self.edges)+len(self.nodes)+2);
	idx := 0
	strs[idx] = "=== Nodes ===";
	idx++;
	for _, node := range(self.nodes) {
		strs[idx] = node.String();
		idx++;
	}
	strs[idx] = "=== Edges ===";
	idx++;
	for _, edge := range(self.edges) {
		strs[idx] = edge.String();
		idx++;
	}
	return strings.Join(strs, "\n");
}

// (self DirectedEdge)
func (self *DsEdge) GetNodes() [2]*DsNode {
	return [2]*DsNode{self.GetSrc(), self.GetDst()};
}

// (self DirectedEdge)
func (self *DsEdge) SetNodes(nodes [2]*DsNode) {
	self.SetSrc(nodes[0]);
	self.SetDst(nodes[1]);
}

/*
// (self Node)
func (self *DsNode) GetNeighbors() []*DsNode {
	edges := self.GetEdges();
	ret := make([]*DsNode, len(edges));
	for idx, edge := range edges {
		nodes := edge.GetNodes();
		if self == nodes[0] {
			ret[idx] = nodes[1];
		} else {
			ret[idx] = nodes[0];
		}
	}
	return ret;
}

// (self Node)
func (self *DsNode) GetEdgeTo(node *DsNode) (*DsEdge, error) {
	edges := self.GetEdges();
	for _, edge := range(edges) {
		nodes := edge.GetNodes();
		if self == nodes[0] {
			if node == nodes[1] {
				return edge, nil;
			}
		} else {
			if node == nodes[0] {
				return edge, nil;
			}
		}
	}
	return nil, ErrNoEdge{src: self, dst: node, graph: self.GetGraph()};
}

// (self Node)
func (self *DsNode) GetEdgeFrom(node *DsNode) (*DsEdge, error) {
	edges := self.GetEdges();
	for _, edge := range(edges) {
		nodes := edge.GetNodes();
		if self == nodes[0] {
			if node == nodes[1] {
				return edge;
			}
		} else {
			if node == nodes[0] {
				return edge;
			}
		}
	}
	return nil, ErrNoEdge{src: self, dst: node, graph: self.GetGraph()};
}
*/

const (
	srch_UNVISITED = iota
	srch_OPEN = iota
	srch_CLOSED = iota
)

/*
func (self *DsGraph) Search(start *DsNode, dist func(*DsEdge) int) (DsGraph, error) {
	if start.graph != self {
		return nil, ErrNotOwned{graph: self, member: start};
	}
	ret := NewGraph();
	status := make(map[*DsNode]int);
	nodes := self.GetAllNodes();
	for _, node := range(nodes) {
		status[node] = srch_UNVISITED;
	}
	status[start] = srch_OPEN;

	return ret,  nil
}
*/
