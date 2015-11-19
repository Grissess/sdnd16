#!/bin/bash

clear

echo "Starting install"

echo "Feching redis go integration package"

go get github.com/garyburd/redigo/redis

echo "Feching graph package"

go get github.com/gyuho/goraph/graph

echo "Feching program package"

go get github.com/Grissess/sdnd16/

echo "Installing client"

go install src/github.com/Grississ/sdnd16/client/client.go

echo "Install Finished"
echo 
