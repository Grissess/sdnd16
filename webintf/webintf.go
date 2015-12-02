package main

import (
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
	t_search.Execute(w, dbs);
}

func view_db(w http.ResponseWriter, r *http.Request, dbname string) {
	t_error.Execute(w, "Not implemented; db requested was "+dbname)
}

func view_node(w http.ResponseWriter, r *http.Request, dbname string, node int) {
	t_error.Execute(w, fmt.Sprintf("Not implemented; (%v %v)", dbname, node))
}

func view_path(w http.ResponseWriter, r *http.Request, dbname string, srcnode, dstnode int) {
	t_error.Execute(w, fmt.Sprintf("Not implemented; (%v %v %v)", dbname, srcnode, dstnode))
}

func db_view(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/");
	if len(parts) >= 3 && len(parts[2]) > 0 {
		dbname := parts[2]
		if len(parts) >= 4 && len(parts[3]) > 0 {
			srcnode, err1 := strconv.Atoi(parts[3])
			if err1 != nil {
				t_error.Execute(w, err1)
			}
			if len(parts) >= 5 && len(parts[4]) > 0 {
				dstnode, err2 := strconv.Atoi(parts[4])
				if err2 != nil {
					t_error.Execute(w, err2)
				}
				view_path(w, r, dbname, srcnode, dstnode)
			}
			view_node(w, r, dbname, srcnode)
		}
		view_db(w, r, dbname)
	}
	view_search(w, r)
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/db/", db_view)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}
