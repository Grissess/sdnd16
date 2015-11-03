package main

import (
	"github.com/Grissess/sdnd16/database"
	"github.com/Grissess/sdnd16/reader"
    "github.com/Grissess/sdnd16/network"
    "fmt"
)

func main() {

    var filename string
    var address string
    var name string
    var nodeLabels []string
    var topology network.DsGraph

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

    topology = reader.ReadFile(filename)
    nodeLabels = reader.LabelList(&topology)
    numberOfNodes := len(nodeLabels)

	rdb := database.NewRoutingDatabase(name, numberOfNodes)
    err := rdb.Connect("tcp", address)
    if err != nil {
        panic(err)
    }

    rdb.SetTrivialPaths()
	for i := 0; i < numberOfNodes; i++ {
		for j := 0; j < numberOfNodes; j++ {
			if i != j {
                rdb.SetPath(i, j, fmt.Sprintf("%d %d | %d", i, j, i + j))
            }
        }
	}
	rdb.StorePathsInDB()

}
