package main

import (
	"os/exec"
	"io"
	"fmt"
	"strings"
	"strconv"
    "net/http"
    "html/template"
	"github.com/Grissess/sdnd16/database"
)

const (
	root = "webintf/"
	db_network = "tcp"
	db_address = "128.153.144.171:6379"
)

var (
	t_search = template.Must(template.ParseFiles(root + "search.gtpl"))
	t_error = template.Must(template.ParseFiles(root + "error.gtpl"))
	t_path = template.Must(template.ParseFiles(root + "path.gtpl"))
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/db")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func view_search(w http.ResponseWriter, r *http.Request) {
	dbs, err := database.GetAllDatabases(db_network, db_address)
	if err != nil {
		t_error.Execute(w, err)
		return
	}
	err2 := t_search.Execute(w, dbs)
	if err2 != nil {
		t_error.Execute(w, err2)
	}
}

func view_db(w http.ResponseWriter, r *http.Request, dbname string) {
	t_error.Execute(w, "Not implemented; db requested was "+dbname)
}

func view_node(w http.ResponseWriter, r *http.Request, dbname, node string) {
	t_error.Execute(w, fmt.Sprintf("Not implemented; (%v %v)", dbname, node))
}

type tinPath struct {
	Rawpath string
	Path []string
	Netpath string
	Cost int
}

func view_path(w http.ResponseWriter, r *http.Request, dbname, srcnode, dstnode string) {
	db, err1 := database.ConnectToDatabase(dbname, db_network, db_address)
	if err1 != nil {
		t_error.Execute(w, err1)
		return
	}
	path, err2 := db.GetPath(srcnode, dstnode)
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
	cost, _ := strconv.Atoi(parts[len(parts) - 1])
	err3 := t_path.Execute(w, tinPath{Rawpath: path, Path: pathpart, Netpath: strings.Join(pathpart, "/"), Cost: cost})
	if err3 != nil {
		t_error.Execute(w, err3)
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

func rend_path(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	nodes := path[3:]
	dotter := exec.Command("dot", "-Tsvg")
	dotin, _ := dotter.StdinPipe()
	dotout, _ := dotter.StdoutPipe()
	w.Header().Set("Content-type", "image/svg+xml")
	dotter.Start()
	io.WriteString(dotin, "digraph {\n")
	for i := 0; i < len(nodes)-1; i++ {
		io.WriteString(dotin, fmt.Sprintf("%s -> %s\n", nodes[i], nodes[i+1]))
	}
	io.WriteString(dotin, "}")
	dotin.Close()
	io.Copy(w, dotout)
}


func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/db/", db_view)
	http.HandleFunc("/render/path/", rend_path)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}
