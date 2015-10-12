// Generic interface casting methods

package network;

import (
	"fmt"
	"strings"
	"container/heap"
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
	return fmt.Sprintf("<Node %p '%s'> { %v }", &self, self.label, self.attrs);
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

func (self DsNode) ToDot() string {
	return fmt.Sprintf("\"%s { %v }\"", self.label, self.attrs);
}

func (self DsEdge) ToDot() string {
	return fmt.Sprintf("%s -> %s [label=\"%v\"]", self.src.ToDot(), self.dst.ToDot(), self.attrs);
}

func (self DsGraph) ToDot() string {
	strs := make([]string, len(self.edges)+len(self.nodes)+2);
	idx := 0
	strs[idx] = "digraph {";
	idx++;
	for _, node := range(self.nodes) {
		strs[idx] = node.ToDot() + ";";
		idx++;
	}
	for _, edge := range(self.edges) {
		strs[idx] = edge.ToDot() + ";";
		idx++;
	}
	strs[idx] = "}";
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

const (
	ATTR_NSSEARCH = ATTR_NSROOT+"srch:"
	ATTR_NSDIST = ATTR_NSSEARCH+"dist"
)

type srchItem struct {
	node *DsNode
	dist int
}

type srchHeap []*DsNode;

func (self srchHeap) Len() int { return len(self); }
func (self srchHeap) Less(i, j int) bool { return self[i].GetAttr(ATTR_NSDIST).(int) < self[j].GetAttr(ATTR_NSDIST).(int); }
func (self srchHeap) Swap(i, j int) { self[i], self[j] = self[j], self[i]; }
func (self *srchHeap) Push(x interface{}) { *self = append(*self, x.(*DsNode)); }
func (self *srchHeap) Pop() interface{} {
	ret := (*self)[len(*self)-1];
	*self = (*self)[:len(*self)-1];
	return ret;
}

func (self *DsGraph) Search(start *DsNode, dist func(*DsEdge) int) (DsGraph, error) {
	ret := NewGraph();
	if start == nil {
		return ret, ErrNullPtr("Search start");
	}
	if start.graph != self {
		return ret, ErrNotOwned{graph: self, member: start};
	}
	_shp := srchHeap(make([]*DsNode, 0, len(self.nodes)));
	shp := &_shp;
	heap.Init(shp);
	status := make(map[*DsNode]int);
	nodes := self.GetAllNodes();
	for _, node := range(nodes) {
		status[ret.GetOrCreateNode(node.label)] = srch_UNVISITED;
	}
	newstart, _ := ret.GetNode(start.label);
	newstart.SetAttr(ATTR_NSDIST, 0);
	start.SetAttr(ATTR_NSDIST, 0);
	heap.Push(shp, start);
	for shp.Len() > 0 {
		node := heap.Pop(shp).(*DsNode);
		newnode := ret.GetOrCreateNode(node.label);
		for _, edge := range(node.GetOutgoing()) {
			dst := edge.GetDst();
			newdst := ret.GetOrCreateNode(dst.label);
			switch status[newdst] {
			case srch_UNVISITED:
				dst.SetAttr(ATTR_NSDIST, node.GetAttr(ATTR_NSDIST).(int) + dist(edge));
				newdst.SetAttr(ATTR_NSDIST, node.GetAttr(ATTR_NSDIST).(int) + dist(edge));
				ret.NewEdge(newnode, newdst);
				heap.Push(shp, dst);
				status[newdst] = srch_OPEN;

			case srch_OPEN:
				newdist := newnode.GetAttr(ATTR_NSDIST).(int) + dist(edge);
				if newdist < newdst.GetAttr(ATTR_NSDIST).(int) {
					dst.SetAttr(ATTR_NSDIST, newdist);
					newdst.SetAttr(ATTR_NSDIST, newdist);
					for _, edge := range(newdst.GetIncoming()) {
						ret.RemoveEdge(edge);
					}
					ret.NewEdge(newnode, newdst);
					heap.Init(shp);  // FIXME
				}
			}
		}
		status[newnode] = srch_CLOSED;
	}
	return ret, nil;
}
