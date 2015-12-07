package main

import (
	"os/exec"
	"errors"
	"io"
	"math"
	"fmt"
	"strings"
	"bytes"
	"strconv"
    "net/http"
    "html/template"
	"github.com/Grissess/sdnd16/database"
	"github.com/Grissess/sdnd16/utils"
	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
	"github.com/gonum/graph/encoding/dot"
)

const (
	root = "webintf/"
	db_network = "tcp"
)

var (
	t_search = template.Must(template.ParseFiles(root + "search.gtpl"))
	t_error = template.Must(template.ParseFiles(root + "error.gtpl"))
	t_path = template.Must(template.ParseFiles(root + "path.gtpl"))
	t_db = template.Must(template.ParseFiles(root + "db.gtpl"))
	t_node = template.Must(template.ParseFiles(root + "node.gtpl"))
	db_address = "128.153.144.171:6379"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/db")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func view_search(w http.ResponseWriter, r *http.Request) {
	dbs, err := database.GetAllTopologies(db_network, db_address)
	if err != nil {
		t_error.Execute(w, err)
		return
	}
	err2 := t_search.Execute(w, dbs)
	if err2 != nil {
		t_error.Execute(w, err2)
	}
}


type tinPath struct {
	Dbname string
	Rawpath string
	Path []string
	Netpath string
	Fullpath string
	Cost int
}

func view_path(w http.ResponseWriter, r *http.Request, dbname, srcnode, dstnode string) {
	db, err1 := database.ConnectToDatabase(dbname, db_network, db_address)
	if err1 != nil {
		t_error.Execute(w, err1)
		return
	}
	labels, err6 := db.GetLabels()
	if err6 != nil {
		t_error.Execute(w, err6)
		return
	}
	srcid, err4 := strconv.Atoi(srcnode)
	if err4 != nil {
		t_error.Execute(w, err4)
		return
	}
	dstid, err5 := strconv.Atoi(dstnode)
	if err5 != nil {
		t_error.Execute(w, err5)
		return
	}
	path, err2 := db.GetPath(labels[srcid], labels[dstid])
	if err2 != nil {
		t_error.Execute(w, err2)
		return
	}
	parts := strings.Split(path, " ")
	if len(parts) < 3 {
		t_error.Execute(w, "No apparent path exists")
		return
	}
	pathpart := parts[:len(parts) - 2]
	revlabels := utils.GetRevLabels(labels)
	idpath := make([]int, len(pathpart))
	for idx, part := range(pathpart) {
		idpath[idx] = revlabels[part]
	}
	sidpath := make([]string, len(idpath))
	for idx, id := range(idpath) {
		sidpath[idx] = fmt.Sprint(id)
	}
	cost, _ := strconv.Atoi(parts[len(parts) - 1])
	err3 := t_path.Execute(w, tinPath{Dbname: dbname, Rawpath: path, Path: pathpart, Netpath: strings.Join(sidpath, "/"), Fullpath: strings.Join(pathpart, " -> "), Cost: cost})
	if err3 != nil {
		t_error.Execute(w, err3)
	}
}

type tinNode struct {
	Dbname string
	Database database.RoutingDatabase
	Node string
}

func view_node(w http.ResponseWriter, r *http.Request, dbname, src string) {
	db, err1 := database.ConnectToDatabase(dbname, db_network, db_address)
	if err1 != nil {
		t_error.Execute(w, err1)
		return
	}
	err2 := t_node.Execute(w, tinNode{Dbname: dbname, Database: db, Node: src})
	if err2 != nil {
		t_error.Execute(w, err2)
	}
}

type tinDb struct {
	Database database.RoutingDatabase
	Dbname string
}

func view_db(w http.ResponseWriter, r *http.Request, dbname string) {
	db, err1 := database.ConnectToDatabase(dbname, db_network, db_address)
	if err1 != nil {
		t_error.Execute(w, err1)
		return
	}
	err2 := t_db.Execute(w, tinDb{Dbname: dbname, Database: db})
	if err2 != nil {
		t_error.Execute(w, err2)
	}
}

func db_view(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/");
	if len(parts) >= 3 && len(parts[2]) > 0 {
		dbname := parts[2]
		if len(parts) >= 4 && len(parts[3]) > 0 {
			srcnode := parts[3]
			if len(parts) >= 5 && len(parts[4]) > 0 {
				dstnode := parts[4]
				view_path(w, r, dbname, srcnode, dstnode)
			}
			view_node(w, r, dbname, srcnode)
		}
		view_db(w, r, dbname)
	}
	view_search(w, r)
}

func get_graph(w http.ResponseWriter, r *http.Request) (graph.Graph, string, *database.RoutingDatabase, error) {
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 4 {
		t_error.Execute(w, "Not enough components in path")
		return nil, "", nil, errors.New("Not enough components in path")
	}
	dbname := path[3]
	db, err1 := database.ConnectToDatabase(dbname, db_network, db_address)
	if err1 != nil {
		t_error.Execute(w, err1)
		return nil, "", nil, err1
	}
	topo, err2 := db.GetTopology()
	if err2 != nil {
		t_error.Execute(w, err2)
		return nil, "", &db, err2
	}
	fmt.Println("topo:")
	fmt.Println(topo)
	g := utils.GraphFromNeighborMap(topo)
	fmt.Println("graph:")
	fmt.Println(g)
	return g, dbname, &db, nil
}

type dbPathGraph struct {
	 *simple.UndirectedGraph
	 dbName string
	 labels map[int]string
	 pathnodes map[int]bool
	 pathslice []int
	 first int
	 last int
}

type dbPathEdge struct {
	simple.Edge
	pathnodes map[int]bool
	pathslice []int
}

type dbPathNode struct {
	simple.Node
	dbName string
	label string
	color string
	first int
}

func (self dbPathGraph) Edge(u, v graph.Node) graph.Edge {
	e := self.UndirectedGraph.Edge(u, v)
	return dbPathEdge{Edge: e.(simple.Edge), pathnodes: self.pathnodes, pathslice: self.pathslice}
}

func (self dbPathGraph) Nodes() []graph.Node {
	nodes := self.UndirectedGraph.Nodes()
	res := make([]graph.Node, len(nodes))
	for idx, elem := range(nodes) {
		color := "#cccccc"
		if elem.ID() == self.first {
			color = "#ffcccc"
		} else if elem.ID() == self.last {
			color = "#ccffcc"
		} else if self.pathnodes[elem.ID()] {
			color = "#ccccff"
		}
		res[idx] = dbPathNode{Node: elem.(simple.Node), dbName: self.dbName, label: self.labels[elem.ID()], color: color, first: self.first}
	}
	return res
}

func (self dbPathEdge) DOTAttributes() []dot.Attribute {
	res := []dot.Attribute{
		dot.Attribute{Key: "label", Value: fmt.Sprintf("%d", int(self.Weight()))},
	}
	if self.pathnodes[self.From().ID()] && self.pathnodes[self.To().ID()] {
		var fidx, tidx int
		for i := 0; i < len(self.pathslice); i++ {
			if self.From().ID() == self.pathslice[i] {
				fidx = i
			}
			if self.To().ID() == self.pathslice[i] {
				tidx = i
			}
		}
		if fidx - tidx == 1 || fidx - tidx == -1 {
			res = append(res, dot.Attribute{Key: "color", Value: "\"#0077ff\""})
		}
	}
	return res
}

func (self dbPathNode) DOTAttributes() []dot.Attribute {
	return []dot.Attribute{
		dot.Attribute{Key: "href", Value: fmt.Sprintf("\"http://localhost:8080/db/%s/%d/%d\"", self.dbName, self.first, self.ID())},
		dot.Attribute{Key: "style", Value: "filled"},
		dot.Attribute{Key: "fillcolor", Value: fmt.Sprintf("\"%s\"", self.color)},
		dot.Attribute{Key: "label", Value: fmt.Sprintf("\"%s\"", self.label)},
	}
}

func (self dbPathNode) ID() int {
	return self.Node.ID()
}

func rend_path(w http.ResponseWriter, r *http.Request) {
	g, n, db, err1 := get_graph(w, r)
	if err1 != nil {
		return  // Already rendered
	}
	labels, err3 := db.GetLabels()
	if err3 != nil {
		t_error.Execute(w, err3)
		return
	}
	path := strings.Split(r.URL.Path, "/")
	nodes := path[4:]
	pathnodes := make(map[int]bool)
	pathslice := make([]int, len(nodes))
	for idx, s := range(nodes) {
		nid, err4 := strconv.Atoi(s)
		if err4 != nil {
			t_error.Execute(w, err4)
			return
		}
		pathnodes[nid] = true
		pathslice[idx] = nid
	}
	first := pathslice[0]
	last := pathslice[len(pathslice) - 1]
	ng := dbPathGraph{UndirectedGraph: g.(*simple.UndirectedGraph), dbName: n, labels: labels, pathnodes: pathnodes, pathslice: pathslice, first: first, last: last}
	data, err2 := dot.Marshal(ng, "", "", "", false)
	fmt.Println("dot:")
	fmt.Println(string(data))
	if err2 != nil {
		t_error.Execute(w, err2)
		return
	}
	buf := bytes.NewBuffer(data)
	dotter := exec.Command("dot", "-Tsvg")
	dotin, _ := dotter.StdinPipe()
	dotout, _ := dotter.StdoutPipe()
	w.Header().Set("Content-type", "image/svg+xml")
	dotter.Start()
	io.Copy(dotin, buf)
	dotin.Close()
	io.Copy(w, dotout)
//	dotter := exec.Command("dot", "-Tsvg")
//	dotin, _ := dotter.StdinPipe()
//	dotout, _ := dotter.StdoutPipe()
//	w.Header().Set("Content-type", "image/svg+xml")
//	dotter.Start()
//	io.WriteString(dotin, "digraph {\n")
//	for i := 0; i < len(nodes)-1; i++ {
//		io.WriteString(dotin, fmt.Sprintf("%s -> %s\n", nodes[i], nodes[i+1]))
//	}
//	io.WriteString(dotin, "}")
//	dotin.Close()
//	io.Copy(w, dotout)
}

type dbNodeGraph struct {
	*simple.UndirectedGraph
	dbName string
	labels map[int]string
}

type dbNodeEdge struct {
	simple.Edge
}

type dbNodeNode struct {
	simple.Node
	dbName string
	label string
}

func (self dbNodeEdge) DOTAttributes() []dot.Attribute {
	return []dot.Attribute{
		dot.Attribute{Key: "label", Value: fmt.Sprintf("%d", int(self.Weight()))},
	}
}

func (self dbNodeNode) DOTAttributes() []dot.Attribute {
	return []dot.Attribute{
		dot.Attribute{Key: "href", Value: fmt.Sprintf("\"http://localhost:8080/db/%s/%d\"", self.dbName, self.ID())},
		dot.Attribute{Key: "style", Value: "filled"},
		dot.Attribute{Key: "fillcolor", Value: "\"#cccccc\""},
		dot.Attribute{Key: "label", Value: fmt.Sprintf("\"%s\"", self.label)},
	}
}

func (self dbNodeNode) ID() int {
	return self.Node.ID()
}

func rend_node(w http.ResponseWriter, r *http.Request) {
	g, n, db, err1 := get_graph(w, r)
	if err1 != nil {
		return  // Already rendered
	}
	parts := strings.Split(r.URL.Path, "/")
	srcid, err3 := strconv.Atoi(parts[4])
	if err3 != nil {
		t_error.Execute(w, err3)
		return
	}
	labels, err4 := db.GetLabels()
	if err4 != nil {
		t_error.Execute(w, err4)
	}
	ng := dbNodeGraph{UndirectedGraph: simple.NewUndirectedGraph(0, math.Inf(1)), dbName: n, labels: labels}
	ng.AddNode(dbNodeNode{Node: simple.Node(srcid), dbName: n, label: labels[srcid]})
	for _, dst := range(g.From(simple.Node(srcid))) {
		ng.AddNode(dbNodeNode{Node: simple.Node(dst.ID()), dbName: n, label: labels[dst.ID()]})
		ng.SetEdge(dbNodeEdge{simple.Edge{F: simple.Node(srcid), T: simple.Node(dst.ID()), W: g.Edge(simple.Node(srcid), dst).Weight()}})
	}
	data, err2 := dot.Marshal(ng, "", "", "", false)
	if err2 != nil {
		t_error.Execute(w, err2)
		return
	}
	fmt.Println(string(data))
	buf := bytes.NewBuffer(data)
	dotter := exec.Command("dot", "-Tsvg")
	dotin, _ := dotter.StdinPipe()
	dotout, _ := dotter.StdoutPipe()
	w.Header().Set("Content-type", "image/svg+xml")
	dotter.Start()
	io.Copy(dotin, buf)
	dotin.Close()
	io.Copy(w, dotout)
}

type dbViewGraph struct {
	 *simple.UndirectedGraph
	 dbName string
	 labels map[int]string
}

type dbViewEdge struct {
	simple.Edge
}

type dbViewNode struct {
	simple.Node
	dbName string
	label string
}

func (self dbViewGraph) Edge(u, v graph.Node) graph.Edge {
	e := self.UndirectedGraph.Edge(u, v)
	return dbViewEdge{e.(simple.Edge)}
}

func (self dbViewGraph) Nodes() []graph.Node {
	nodes := self.UndirectedGraph.Nodes()
	res := make([]graph.Node, len(nodes))
	for idx, elem := range(nodes) {
		res[idx] = dbViewNode{Node: elem.(simple.Node), dbName: self.dbName, label: self.labels[elem.ID()]}
	}
	return res
}

func (self dbViewEdge) DOTAttributes() []dot.Attribute {
	return []dot.Attribute{
		dot.Attribute{Key: "label", Value: fmt.Sprintf("%d", int(self.Weight()))},
	}
}

func (self dbViewNode) DOTAttributes() []dot.Attribute {
	return []dot.Attribute{
		dot.Attribute{Key: "href", Value: fmt.Sprintf("\"http://localhost:8080/db/%s/%d\"", self.dbName, self.ID())},
		dot.Attribute{Key: "style", Value: "filled"},
		dot.Attribute{Key: "fillcolor", Value: "\"#cccccc\""},
		dot.Attribute{Key: "label", Value: fmt.Sprintf("\"%s\"", self.label)},
	}
}

func (self dbViewNode) ID() int {
	return self.Node.ID()
}

func rend_db(w http.ResponseWriter, r *http.Request) {
	g, n, db, err1 := get_graph(w, r)
	if err1 != nil {
		return  // Already rendered
	}
	labels, err3 := db.GetLabels()
	if err3 != nil {
		t_error.Execute(w, err3)
		return
	}
	data, err2 := dot.Marshal(dbViewGraph{UndirectedGraph: g.(*simple.UndirectedGraph), dbName: n, labels: labels}, "", "", "", false)
	if err2 != nil {
		t_error.Execute(w, err2)
		return
	}
	fmt.Println(string(data))
	buf := bytes.NewBuffer(data)
	dotter := exec.Command("dot", "-Tsvg")
	dotin, _ := dotter.StdinPipe()
	dotout, _ := dotter.StdoutPipe()
	w.Header().Set("Content-type", "image/svg+xml")
	dotter.Start()
	io.Copy(dotin, buf)
	dotin.Close()
	io.Copy(w, dotout)
}

func main() {
	var in string
    http.HandleFunc("/", index)
    http.HandleFunc("/db/", db_view)
	http.HandleFunc("/render/path/", rend_path)
	http.HandleFunc("/render/db/", rend_db)
	http.HandleFunc("/render/node/", rend_node)
	fmt.Printf("Enter address, or leave blank for default(%s): ", db_address)
	fmt.Scanln(&in)
	if in != "" {
		db_address = in
	}
	fmt.Println("Ready.")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}
