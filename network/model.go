// Actual networking model

package network;

import (
	"fmt"
)

// Types

type DsEdge struct {
	graph *DsGraph
	src, dst *DsNode
	attrs map[string]interface{}
}

type DsNode struct {
	graph *DsGraph
	label Label
	outgoing map[*DsNode]*DsEdge
	incoming map[*DsNode]*DsEdge
	attrs map[string]interface{}
}

type DsGraph struct {
	edges []DsEdge
	nodes []DsNode
	labels map[Label]*DsNode
	attrs map[string]interface{}
}

// Implementations

// Edge--get/set

func (self *DsEdge) GetGraph() *DsGraph {
	return self.graph;
}

func (self *DsEdge) GetAttr(attr string) interface{} {
	return self.attrs[attr];
}

func (self *DsEdge) SetAttr(attr string, val interface{}) {
	self.attrs[attr] = val;
}

func (self *DsEdge) GetSrc() *DsNode {
	return self.src;
}

func (self *DsEdge) SetSrc(src *DsNode) {
	self.src = src;
}

func (self *DsEdge) GetDst() *DsNode {
	return self.dst;
}

func (self *DsEdge) SetDst(dst *DsNode) {
	self.dst = dst;
}

// Node--get/set

func (self *DsNode) GetGraph() *DsGraph {
	return self.graph;
}

func (self *DsNode) GetLabel() Label {
        return self.label;
}

func (self *DsNode) GetAttr(attr string) interface{} {
	return self.attrs[attr];
}

func (self *DsNode) SetAttr(attr string, val interface{}) {
	self.attrs[attr] = val;
}

func (self *DsNode) GetEdges() []*DsEdge {
	ret := make([]*DsEdge, len(self.outgoing)+len(self.incoming));
	idx := 0;
	for _, edge := range(self.outgoing) {
		ret[idx] = edge;
		idx++;
	}
	for _, edge := range(self.incoming) {
		ret[idx] = edge;
		idx++;
	}
	return ret;
}

func (self *DsNode) GetOutgoing() []*DsEdge {
	ret := make([]*DsEdge, len(self.outgoing));
	idx := 0;
	for _, edge := range(self.outgoing) {
		ret[idx] = edge;
		idx++;
	}
	return ret;
}

func (self *DsNode) GetIncoming() []*DsEdge {
	ret := make([]*DsEdge, len(self.incoming));
	idx := 0;
	for _, edge := range(self.incoming) {
		ret[idx] = edge;
		idx++;
	}
	return ret;
}

func (self *DsNode) GetNeighbors() []*DsNode {
	ret := make([]*DsNode, len(self.outgoing));
	idx := 0;
	for neighbor, _ := range(self.outgoing) {
		ret[idx] = neighbor;
		idx++;
	}
	return ret;
}

func (self *DsNode) GetEdgeTo(node *DsNode) (*DsEdge, error) {
	edge, exists := self.outgoing[node];
	if !exists {
		return nil, ErrNoEdge{src: self, dst: node, graph: self.graph};
	}
	return edge, nil;
}

func (self *DsNode) GetEdgeFrom(node *DsNode) (*DsEdge, error) {
	edge, exists := self.incoming[node];
	if !exists {
		return nil, ErrNoEdge{src: self, dst: node, graph: self.graph};
	}
	return edge, nil;
}

// Graph--get/set

func (self *DsGraph) GetAttr(attr string) interface{} {
	return self.attrs[attr];
}

func (self *DsGraph) SetAttr(attr string, val interface{}) {
	self.attrs[attr] = val;
}

func (self *DsGraph) GetAllNodes() []*DsNode {
	ret := make([]*DsNode, len(self.nodes));
	for idx := range(self.nodes) {
		ret[idx] = &(self.nodes[idx]);
	}
	return ret;
}

func (self *DsGraph) GetAllEdges() []*DsEdge {
	ret := make([]*DsEdge, len(self.edges));
	for idx, edge := range(self.edges) {
		ret[idx] = &edge;
	}
	return ret;
}

func (self *DsGraph) GetNode(label Label) (*DsNode, error) {
	node, exists := self.labels[label];
	if !exists {
		return nil, ErrNoNode{label: label, graph: self};
	}
	return node, nil;
}

func (self *DsGraph) NewNode(label Label) (*DsNode, error) {
	noderef, err := self.GetNode(label);
	if err == nil {
		return nil, ErrAlreadyNode{node: noderef, graph: self};
	}
	node := DsNode{graph: self, label: label, outgoing: make(map[*DsNode]*DsEdge), incoming: make(map[*DsNode]*DsEdge), attrs: make(map[string]interface{})};
	self.nodes = append(self.nodes, node);
	noderef = &self.nodes[len(self.nodes) - 1];
	self.labels[label] = noderef;
	return noderef, nil;
}

func (self *DsGraph) GetOrCreateNode(label Label) *DsNode {
	node, err := self.GetNode(label);
	if err == nil {
		fmt.Printf("// GOC(%p): Created %s\n", self, node);
		return node;
	}
	node, err = self.NewNode(label);
	if node == nil {
		panic(fmt.Sprintf("Node non-extant and not created; creation error: %s", err.Error()));
	}
	fmt.Printf("// GOC(%p): Returning %s\n", self, node);
	return node;
}

func (self *DsGraph) NewEdge(src *DsNode, dst *DsNode) (*DsEdge, error) {
	if src == nil {
		return nil, ErrNullPtr("NewEdge src");
	}
	if dst == nil {
		return nil, ErrNullPtr("NewEdge dst");
	}
	if src.graph != self {
		return nil, ErrNotOwned{graph: self, member: src};
	}
	if dst.graph != self {
		return nil, ErrNotOwned{graph: self, member: dst};
	}
	edgeref, existing := src.outgoing[dst];
	if existing {
		return nil, ErrAlreadyEdge{edge: edgeref, graph: self};
	}
	edge := DsEdge{graph: self, src: src, dst: dst, attrs: make(map[string]interface{})};
	self.edges = append(self.edges, edge);
	edgeref = &self.edges[len(self.edges) - 1];
	src.outgoing[dst] = edgeref;
	dst.incoming[src] = edgeref;
	return edgeref, nil;
}

func (self *DsGraph) RemoveEdge(edge *DsEdge) bool {
	if edge == nil {
		return false;
	}
	if edge.graph != self {
		fmt.Printf("// WARNING: Removing edge %v (%p) from non-owning graph %p (owned by %p)\n", edge, edge, self, edge.graph);
		return false;
	}
	if _, found := edge.src.outgoing[edge.dst]; !found {
		fmt.Printf("// SANITY: Source does not register outgoing edge %p!\n", edge);
	}
	if _, found := edge.dst.incoming[edge.src]; !found {
		fmt.Printf("// SANITY: Destination does not register incoming edge %p!\n", edge);
	}
	delete(edge.src.outgoing, edge.dst);
	delete(edge.dst.incoming, edge.src);
	for idx := range(self.edges) {  // FIXME
		if edge == &self.edges[idx] {
			self.edges = append(self.edges[:idx], self.edges[idx+1:]...);
			return true;
		}
	}
	return false;
}

func (self *DsGraph) RemoveNode(node *DsNode) {
	if node == nil {
		return;
	}
	if node.graph != self {
		fmt.Printf("// WARNING: Removing node %v (%p) from non-owning graph %p (owned by %p)\n", node, node, self, node.graph);
		return;
	}
	for _, edge := range(node.GetEdges()) {
		if !self.RemoveEdge(edge) {  // FIXME
			fmt.Printf("// SANITY: Failed to remove an edge %p returned in the set of all edges of %s!\n", edge, node);
		}
	}
	for idx := range(self.nodes) {  // FIXME
		if node == &self.nodes[idx] {
			self.nodes = append(self.nodes[:idx], self.nodes[idx+1:]...);
			return;
		}
	}
}

func NewGraph() *DsGraph {
	return &DsGraph{edges: nil, nodes: nil, labels: make(map[Label]*DsNode), attrs: make(map[string]interface{})};
}
