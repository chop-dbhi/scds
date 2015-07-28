package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var (
	host         string
	port         int
	mongoURI     string
	mongoSession *mgo.Session
)

func main() {
	flag.StringVar(&host, "host", "", "Host")
	flag.IntVar(&port, "port", 5000, "Port")
	flag.StringVar(&mongoURI, "mongo", "localhost/scds", "URI of the MongoDB host or cluster.")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		PrintUsage("help")
	}

	// Route command.
	switch args[0] {
	case "put":
		session := initDB()
		defer session.Close()

		putCmd(session, args[1:])

	case "get":
		session := initDB()
		defer session.Close()

		getCmd(session, args[1:])

	case "log":
		session := initDB()
		defer session.Close()

		logCmd(session, args[1:])

	case "http":
		session := initDB()
		defer session.Close()

		httpCmd(session, args[1:])

	default:
		// Print usage of speific command.
		if len(args) == 2 {
			PrintUsage(args[1])
		}

		PrintUsage("help")
	}
}
