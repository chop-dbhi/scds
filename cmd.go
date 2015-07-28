package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func putCmd(cfg *Config, args []string) {
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

	defer cfg.Mongo.Close()
	o, err := Put(cfg, args[0], val)

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

func getCmd(cfg *Config, args []string) {
	var (
		v  int
		ts string
	)

	fs := flag.NewFlagSet("get", flag.ExitOnError)

	fs.IntVar(&v, "version", 0, "Specific revision to get.")
	fs.StringVar(&ts, "time", "", "Returns the object as of the specified time.")

	fs.Parse(args)

	args = fs.Args()

	if len(args) != 1 {
		PrintUsage("get")
	}

	t, err := ParseTimeString(ts)

	if v > 0 && t > 0 {
		fmt.Println("error: version and time are mutually exclusive\n")
		PrintUsage("get")
	}

	defer cfg.Mongo.Close()

	o, err := get(cfg, args[0], true)

	if err != nil {
		log.Fatal(err)
	}

	if o == nil {
		return
	}

	if v > 0 {
		o = o.AtVersion(v)
	} else if t > 0 {
		o = o.AtTime(t)
	}

	if o == nil {
		return
	}

	o.History = nil

	b, err := json.MarshalIndent(o, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", b)
}

func logCmd(cfg *Config, args []string) {
	if len(args) != 1 {
		PrintUsage("log")
	}

	defer cfg.Mongo.Close()
	l, err := Log(cfg, args[0])

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

func httpCmd(cfg *Config, args []string) {
	fs := flag.NewFlagSet("http", flag.ExitOnError)

	fs.StringVar(&cfg.HTTP.Host, "host", "localhost", "Host to bind to.")
	fs.IntVar(&cfg.HTTP.Port, "port", 5000, "Port to bind to.")

	fs.Parse(args)

	defer cfg.Mongo.Close()

	runHTTP(cfg)
}
