#!/bin/bash

echo "Starting install"

echo "Feching redis go integration package"

go get github.com/garyburd/redigo/redis

echo "Feching graph package"

go get github.com/gonum/graph
go get github.com/gonum/matrix/mat64

echo "Fetching various utilities"

go get github.com/eapache/queue
go get github.com/fatih/color

echo "Fetchin most recent cleint"

go get github.com/Grissess/sdnd16

echo "Install Finished"
