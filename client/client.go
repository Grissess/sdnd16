package main

import (
	"fmt"
        "strings"
	"github.com/Grissess/sdnd16/database"
	"github.com/Grissess/sdnd16/reader"
	"github.com/gyuho/goraph/graph"
)

func main() {

	var filename string
	var address string
	var name string
	var nodeLabels []string
	var topology *graph.DefaultGraph

	var start string
	fmt.Print("grab or store? > ")
	fmt.Scanln(&start)

	if start == "store" {

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

		topology = reader.ReadFileToGraph(filename)
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
		fmt.Println("Paths set, storing paths")

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
		rdb.StorePathsInDB()
		fmt.Println("Paths stored in data base")
                rdb.Disconnect()

	} else if start == "grab" {
		fmt.Print("enter address and port of database > ")
		fmt.Scanln(&address)

		if address == "" {
			fmt.Println("- no address specified, using default database")
			address = "128.153.144.171:6379"
		}

/*		err := rdb.Connect("tcp", address)
		if err != nil {
			panic(err)
		}*/
		fmt.Println("Connected")

	} else {
		fmt.Println("invalid input program terminated")
	}

}
