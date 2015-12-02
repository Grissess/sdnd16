//The test client for this project, allows for storage and retrieval of information
package main

import (
	"fmt"
	"github.com/Grissess/sdnd16/database"
	"github.com/Grissess/sdnd16/utils"
        "github.com/Grissess/sdnd16/algorithms"
        "github.com/gonum/graph/path"
)

func main() {

        var ipAddress, recordName, filename, input string
        var db database.RoutingDatabase
        var DBerr error

        fmt.Println("Welcome to the NRA System, written in golang")
        fmt.Println("Enter IP address and port of the database server you'd like to use")
        fmt.Println("If none is entered, the default server will be used (testing only)\n")
        fmt.Print("IP address and port >> ")
        fmt.Scanln(&ipAddress)

        if ipAddress == "" {
                fmt.Println("No ip address selected, using default database")
                ipAddress = "128.153.144.171:6379"
        }

        fmt.Print("Enter the name of the database record you wish to use >> ")
        fmt.Scanln(&recordName)

        exists, DBerr :=  database.DatabaseExists(recordName, "tcp", ipAddress)

        if DBerr != nil {
                panic(DBerr)
        }

        if exists {
                fmt.Println("This record exists!")
                db, DBerr = database.ConnectToDatabase(recordName, "tcp", ipAddress)

                if DBerr != nil {
                        panic(DBerr)
                }
                queryPaths(db)

       } else {
                fmt.Println("This record does not exist.  Please check your spelling")
                fmt.Print("Would you like to create a record by this name? [Y/N] >> ")
                fmt.Scanln(&input)

                if input ==  "Y" || input == "y" {
                        fmt.Print("Give me the name of a topology file >> ")
                        fmt.Scanln(&filename)
                        g, labels, Uerror := utils.ReadFileToGraph(filename)

                        if Uerror != nil {
                                panic(Uerror)
                        }

                        revLabels := utils.GetRevLabels(labels)
                        paths := path.DijkstraAllPaths(g)
                        pathMap := algorithms.ConvertAllPaths(g, paths)

                        realPathMap := make(map[int]map[int]string)

                        for k  := range pathMap{
                               smallMap := pathMap[k]
                               newMap := make(map[int]string)
                               for k2 := range  smallMap{
                                        stringForPath := smallMap[k2].PathString(labels)
                                        newMap[k2] = stringForPath
                                }
                                realPathMap[k] = newMap
                        }

                        db, DBerr = database.NewRoutingDatabase(recordName, "tcp", ipAddress, revLabels, realPathMap)

                        if DBerr != nil {
                                panic(DBerr)
                        }

                        queryPaths(db)

                } else {
                        fmt.Println("Very well.")
                        fmt.Println("Thank you for using golang NRA!")
                }
        }
}

func queryPaths(db database.RoutingDatabase) {
        var input, src, dst string

        for {
                fmt.Print("Get node paths? [Y/N] >> ")
                fmt.Scanln(&input)
                if input == "y" || input == "Y" {
                        fmt.Print("Enter a first node >> ")
                        fmt.Scanln(&src)
                        fmt.Print("Enter a second node >>")
                        fmt.Scanln(&dst)
                        path, DBerr := db.GetPath(src, dst)

                        if DBerr != nil {
                                fmt.Println(DBerr)
                        } else {
                        fmt.Println("The shortest path is: ", path)
                        }
                 }else {
                                fmt.Println("Thank you for using golang NRA!")
                                break
                 }
        }
}
