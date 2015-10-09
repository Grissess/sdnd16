// Interfaces for use in network functions

package network;

type Label string;

type ErrNullPtr string;

type ErrNoNode struct {
	label Label
	graph *DsGraph
}

type ErrAlreadyNode struct {
	node *DsNode
	graph *DsGraph
}

type ErrNoEdge struct {
	src, dst *DsNode
	graph *DsGraph
}

type ErrAlreadyEdge struct {
	edge *DsEdge
	graph *DsGraph
}

type ErrNotOwned struct {
	graph *DsGraph
	member GraphMember
}

type Attributed interface {
	GetAttr(string) interface{}
	SetAttr(string, interface{})
}

type GraphMember interface {
	GetGraph() *DsGraph
}

type Edge interface {
	GetNodes() [2]*DsNode
	SetNodes([2]*DsNode)
}

type DirectedEdge interface {
	Edge
	GetSrc() *DsNode
	SetSrc(*DsNode)
	GetDst() *DsNode
	SetDst(*DsNode)
}

type Node interface {
	GetEdges() []*DsEdge
	GetNeighbors() []*DsNode
	GetEdgeTo(*DsNode) (*DsEdge, error)
	GetEdgeFrom(*DsNode) (*DsEdge, error)
}

type DirectedNode interface {
	Node
	GetOutgoing() []*DsEdge
	GetIncoming() []*DsEdge
}

type Graph interface {
	GetAllEdges() []*DsEdge
	GetAllNodes() []*DsNode
	GetNode(Label) (*DsNode, error)
	NewNode(Label) (*DsNode, error)
	GetOrCreateNode(Label) *DsNode
	NewEdge(*DsNode, *DsNode) *DsEdge
}
