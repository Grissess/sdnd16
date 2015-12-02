//The test client for this project, allows for storage and retrieval of information
package main

import (
	"fmt"
	"github.com/Grissess/sdnd16/database"
	"github.com/Grissess/sdnd16/utils"
        "github.com/Grissess/sdnd16/algorithms"
        "github.com/gonum/graph/path"
        "github.com/fatih/color"
)

func main() {

        var ipAddress, recordName, filename, input string
        var db database.RoutingDatabase
        var DBerr error

        color.Set(color.FgHiGreen, color.Bold, color.BlinkSlow)
        fmt.Println("\nWelcome to the NRA System, written in golang\n")
        color.Unset()


        fmt.Println("Enter IP address and port of the database server you'd like to use")
        fmt.Println("If none is entered, the default server will be used (testing only)\n")
        color.Unset()

        color.Set(color.FgWhite)
        fmt.Print("IP address and port >> ")
        color.Unset()

        fmt.Scanln(&ipAddress)

        if ipAddress == "" {
                color.Set(color.FgHiBlue)
                fmt.Println("No ip address selected, using default database")
                color.Unset()
                ipAddress = "128.153.144.171:6379"
        }

        color.Set(color.FgWhite)
        fmt.Print("Enter the name of the database record you wish to use >> ")
        color.Unset()
        fmt.Scanln(&recordName)

        exists, DBerr :=  database.DatabaseExists(recordName, "tcp", ipAddress)

        if DBerr != nil {
                panic(DBerr)
        }

        if exists {
                color.Set(color.FgHiBlue)
                fmt.Println("This record exists!\n")
                color.Unset()

                db, DBerr = database.ConnectToDatabase(recordName, "tcp", ipAddress)

                if DBerr != nil {
                        panic(DBerr)
                }
                queryPaths(db)

       } else {
                color.Set(color.FgHiBlue)
                fmt.Println("This record does not exist.  Please check your spelling")
                color.Set(color.FgWhite)
                fmt.Print("Would you like to create a record by this name? [Y/N] >> ")
                color.Unset()
                fmt.Scanln(&input)

                if input ==  "Y" || input == "y" {
                        color.Set(color.FgWhite)
                        fmt.Print("Give me the name of a topology file >> ")
                        color.Unset() 

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

                        color.Set(color.FgHiGreen, color.Bold)
                        fmt.Println("Database created correctly!")
                        color.Unset()

                        queryPaths(db)

                } else {
                        color.Set(color.FgHiGreen, color.Bold)
                        fmt.Println("Thank you for using golang NRA!")
                        color.Unset()
                }
        }
}

func queryPaths(db database.RoutingDatabase) {
        var input, src, dst string

        for {
                color.Set(color.FgWhite)
                fmt.Print("Get node paths? [Y/N] >> ")
                color.Unset()
                fmt.Scanln(&input)
                if input == "y" || input == "Y" {
                        color.Set(color.FgWhite)
                        fmt.Print("Enter a first node >> ")
                        color.Unset()
                        fmt.Scanln(&src)

                        color.Set(color.FgWhite)
                        fmt.Print("Enter a second node >>")
                        color.Unset()
                        fmt.Scanln(&dst)

                        path, DBerr := db.GetPath(src, dst)

                        if DBerr != nil {
                                color.Set(color.FgHiRed, color.Bold)
                                fmt.Println(DBerr, "\n")
                                color.Unset()
                        } else {
                        color.Set(color.FgGreen, color.Bold)
                        fmt.Println("The shortest path is: ", path, "\n")
                        color.Unset()

                        }
                 }else {
                                color.Set(color.FgHiGreen, color.Bold)
                                fmt.Println("Thank you for using golang NRA!")
                                color.Unset()
                                break
                 }
        }
}
