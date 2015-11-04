package main

import (
	"fmt"
	"github.com/Grissess/sdnd16/database"
	"github.com/Grissess/sdnd16/network"
	"github.com/Grissess/sdnd16/reader"
	"os"
)

func main() {

	done := 0
	var filename string
	var address string
	var name string
	var nodeLabels []string
	var topology network.DsGraph

	if len(os.Args) == 1 {

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
					rdb.SetPath(i, j, fmt.Sprintf("%d %d | %d", i, j, i+j))
				}
			}
		}
		rdb.StorePathsInDB()

	} else if len(os.Args) == 3 {
		fmt.Println("- no address specified, using default database")

		topology = reader.ReadFile(os.Args[1])
		nodeLabels = reader.LabelList(&topology)
		numberOfNodes := len(nodeLabels)

		rdb := database.NewRoutingDatabase(os.Args[2], numberOfNodes)
		err := rdb.Connect("tcp", "128.153.144.171:6379")
		if err != nil {
			panic(err)
		}

		rdb.SetTrivialPaths()
		for i := 0; i < numberOfNodes; i++ {
			for j := 0; j < numberOfNodes; j++ {
				if i != j {
					rdb.SetPath(i, j, fmt.Sprintf("%d %d | %d", i, j, i+j))
				}
			}
		}
		rdb.StorePathsInDB()

	} else if len(os.Args) == 4 {
		topology = reader.ReadFile(os.Args[1])
		nodeLabels = reader.LabelList(&topology)
		numberOfNodes := len(nodeLabels)

		rdb := database.NewRoutingDatabase(os.Args[2], numberOfNodes)
		err := rdb.Connect("tcp", os.Args[3])
		if err != nil {
			panic(err)
		}

		rdb.SetTrivialPaths()
		for i := 0; i < numberOfNodes; i++ {
			for j := 0; j < numberOfNodes; j++ {
				if i != j {
					rdb.SetPath(i, j, fmt.Sprintf("%d %d | %d", i, j, i+j))
				}
			}
		}
		rdb.StorePathsInDB()

	} else {
		fmt.Println("Invalid input program Terminated")
		done = 1
	}
	var answer string
	for !(done == 0) {
		var start int
		var end int
		fmt.Print("what is your starting node? > ")
		fmt.Scanln(&start)
		fmt.Print("what is your ending node? > ")
		fmt.Scanln(&end)
		fmt.Println("first path is")
		fmt.Println("second path is")
		fmt.Print("enter q to quit")
		fmt.Scanln(&answer)
		if answer == "q" {
			done = 0
		}
	}
}
