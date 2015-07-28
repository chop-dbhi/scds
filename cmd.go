package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

func putCmd(s *mgo.Session, args []string) {
	if len(args) != 1 {
		PrintUsage("put")
	}

	var (
		err error
		val map[string]interface{}
	)

	// Decoded provided argument other read from stdin.
	if len(args) == 2 {
		err = json.Unmarshal([]byte(args[1]), &val)
	} else {
		err = json.NewDecoder(os.Stdin).Decode(&val)
	}

	o, err := Put(s, args[0], val)

	if err != nil {
		log.Fatal(err)
	}

	if o == nil {
		return
	}

	b, err := json.MarshalIndent(o, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", b)
}

func getCmd(s *mgo.Session, args []string) {
	if len(args) != 1 {
		PrintUsage("get")
	}

	o, err := Get(s, args[0])

	if err != nil {
		log.Fatal(err)
	}

	if o == nil {
		return
	}

	b, err := json.MarshalIndent(o, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", b)
}

func logCmd(s *mgo.Session, args []string) {
	if len(args) != 1 {
		PrintUsage("log")
	}

	l, err := Log(s, args[0])

	if err != nil {
		log.Fatal(err)
	}

	if l == nil {
		return
	}

	b, err := json.MarshalIndent(l, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", b)
}

func httpCmd(s *mgo.Session, args []string) {
	runHTTP(s)
}
