package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func putCmd(args []string) {
	if len(args) < 1 {
		PrintUsage("put")
	}

	var (
		err error
		val map[string]interface{}
	)

	// Decode argument or read from stdin.
	if len(args) == 2 {
		err = json.Unmarshal([]byte(args[1]), &val)
	} else {
		err = json.NewDecoder(os.Stdin).Decode(&val)
	}

	if err != nil {
		log.Fatal(err)
	}

	cfg := GetConfig()

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

func getCmd(args []string) {
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

	cfg := GetConfig()

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

func keysCmd(args []string) {
	cfg := GetConfig()

	defer cfg.Mongo.Close()

	keys, err := Keys(cfg)

	if err != nil {
		log.Fatal(err)
	}

	if len(keys) == 0 {
		return
	}

	fmt.Fprintln(os.Stdout, strings.Join(keys, "\n"))
}

func logCmd(args []string) {
	if len(args) != 1 {
		PrintUsage("log")
	}

	cfg := GetConfig()

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

func httpCmd(args []string) {
	fs := flag.NewFlagSet("http", flag.ExitOnError)

	fs.String("host", "localhost", "Host to bind to.")
	fs.Int("port", 5000, "Port to bind to.")

	fs.Parse(args)

	fs.Visit(func(f *flag.Flag) {
		viper.Set(fmt.Sprintf("http.%s", f.Name), f.Value.(flag.Getter).Get())
	})

	cfg := GetConfig()

	defer cfg.Mongo.Close()

	runHTTP(cfg)
}

func configCmd(args []string) {
	b, _ := yaml.Marshal(GetConfig())

	fmt.Fprintf(os.Stdout, string(b))
}

func subscribeCmd(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No emails provided.")
	}

	cfg := GetConfig()

	n, err := SubscribeEmail(cfg, args...)

	if err != nil {
		log.Fatal(err)
	}

	if n == 1 {
		fmt.Fprintln(os.Stdout, "Subscribed 1 new email")
	} else {
		fmt.Fprintf(os.Stdout, "Subscribed %d new emails\n", n)
	}
}

func unsubscribeCmd(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No emails provided.")
	}

	cfg := GetConfig()

	n, err := UnsubscribeEmail(cfg, args...)

	if err != nil {
		log.Fatal(err)
	}

	if n == 1 {
		fmt.Fprintln(os.Stdout, "Unsubscribed 1 email")
	} else {
		fmt.Fprintf(os.Stdout, "Unsubscribed %d emails\n", n)
	}
}
