package main

import(
    "fmt"
    "net/http"
    "html/template"
)

func index(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        
    }
}

func view(w http.ResponseWriter, r *http.Request) {
}

func main() {
    http.HandleFunc("/", index)
    http.HandleFunc("/view", view)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}
