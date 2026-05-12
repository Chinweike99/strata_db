package main

import (
	"flag"
)


func main() {
	port := flag.String("port", "5432", "Database port")
	flag.Parse()

	db := engine.NewDatabase()
	srv := server.NewTCPServer(*port, db)
}