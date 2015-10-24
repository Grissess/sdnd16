package main

import "fmt"
import "github.com/garyburd/redigo/redis"

func main() {
//INIT OMIT
c, err := redis.Dial("tcp", ":6379")
if err != nil {
panic(err)
}
defer c.Close()

//set
c.Do("SET", "message1", "hello!")
c.Do("hset", "Hello", "node1page1", "3-5-1:9")

//get
world, err := redis.String(c.Do("GET", "message1"))

if err != nil {
fmt.Println("key not found")
}

fmt.Println(world)
//ENDINIT OMIT
}
