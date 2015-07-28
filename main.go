package main

import "flag"

func main() {
	cfg := new(Config)

	flag.BoolVar(&cfg.Debug, "debug", false, "Turn on debug output.")

	flag.StringVar(&cfg.Mongo.URI, "mongo.uri", "localhost/scds", "URI of the MongoDB host or cluster.")

	flag.StringVar(&cfg.SMTP.Host, "smtp.host", "localhost", "Host of the SMTP server.")
	flag.IntVar(&cfg.SMTP.Port, "smtp.port", 25, "Port of the SMTP server.")
	flag.StringVar(&cfg.SMTP.Host, "smtp.user", "", "SMTP user.")
	flag.StringVar(&cfg.SMTP.Password, "smtp.password", "", "SMTP password.")
	flag.StringVar(&cfg.SMTP.From, "smtp.from", "", "SMTP From address.")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		PrintUsage("help")
	}

	// Route command.
	switch args[0] {
	case "put":
		putCmd(cfg, args[1:])

	case "get":
		getCmd(cfg, args[1:])

	case "log":
		logCmd(cfg, args[1:])

	case "http":
		httpCmd(cfg, args[1:])

	default:
		// Print usage of speific command.
		if len(args) == 2 {
			PrintUsage(args[1])
		}

		PrintUsage("help")
	}
}
