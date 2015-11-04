//The test client for this project, allows for storage and retrieval of information
package main

import (
        "fmt"
	"../database"
	"../reader"
	"github.com/gyuho/goraph/graph"
	"strings"
)
func main() {

	var filename string
	var address string
	var name string
	var nodeLabels []string
	var topology *graph.DefaultGraph
	var fileErr error

	var input string
	fmt.Print("grab or store? > ")
	fmt.Scanln(&input)

	if input == "store" {

		fmt.Print("enter topology filename > ")
		fmt.Scanln(&filename)
		fmt.Print("enter topology name > ")
		fmt.Scanln(&name)
		fmt.Print("enter address and port of database > ")
		fmt.Scanln(&address)

		if address == "" {
			fmt.Println("- no address specified, using default database")
			address = "128.153.144.171:6379"
		}

		topology, fileErr = reader.ReadFileToGraph(filename)

		if fileErr != nil {
			panic(fileErr)
		}
		nodeLabels = reader.GetLabelList(topology)
		numberOfNodes := len(nodeLabels)
		labelMap := reader.GetLabelMap(topology)

		rdb, err := database.NewRoutingDatabase(name, "tcp", address, labelMap)
		fmt.Println("Connecting to data base")

		if err != nil {
			panic(err)
		}
		fmt.Println("Connected")

		fmt.Println("Setting paths")

		var src, dest string

		for i := 0; i < numberOfNodes; i++ {
			for j := 0; j < numberOfNodes; j++ {
				if i != j {
					src = nodeLabels[i]
					dest = nodeLabels[j]
					paths, distance, _ := graph.Dijkstra(topology, src, dest)
					path := strings.Join(paths[1:], " ")
					rdb.SetPath(src, dest, fmt.Sprintf("%s %s | %d", src, path, int(distance[dest])))
				}
			}
		}

		fmt.Println("Paths set, storing paths")
		rdb.StorePathsInDB()
		fmt.Println("Paths stored in data base")
		rdb.Disconnect()

	} else if input == "grab" {
                fmt.Print("enter topology name > ")
		fmt.Scanln(&name)
                fmt.Print("enter address and port of database > ")
		fmt.Scanln(&address)

		if address == "" {
			fmt.Println("- no address specified, using default database")
			address = "128.153.144.171:6379"
		}
                rdb,err := database.NewRoutingDatabaseFromDB(name, "tcp", address)
                if err != nil {
			panic(err)
		}
		fmt.Println("Connected")
                        var src, dest string

                for{
                        fmt.Print("Enter a first node > ")
                        fmt.Scanln(&input)
                        src = input
                        fmt.Print("Enter a second node > ")
                        fmt.Scanln(&input)
                        dest = input

                        path, DBerr := rdb.GetPathFromDB(src, dest)
                        if DBerr != nil{
                                panic(DBerr)
                        }

                        fmt.Println("The shortest path is: ", path);
                }
                rdb.Disconnect()
	} else {
		fmt.Println("invalid input program terminated")
	}

}
