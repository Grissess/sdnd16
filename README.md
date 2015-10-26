# sdnd16
This repo is dedicated to Team 5s project for CS350 (software design & development) at Clarkson University.
The objective of this project is to create a system which determines shortest paths, and backup paths, between nodes on a graph (network topology).
Topologies are given as flat text files which are then read into an internal graph structure which is stored in a database.
Once the topology is stored internally and in a database a shortest paths algorithm is run across the graph and the resulting paths are stored in a database as well.
This project is implemented in goland using the redis NoSQL database as per client specification.
