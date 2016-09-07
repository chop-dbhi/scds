package main

import (
	"flag"

	"github.com/blang/semver"
	"github.com/spf13/viper"
)

var (
	progVersion = semver.Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
		Pre: []semver.PRVersion{
			{VersionStr: "beta"},
			{VersionNum: 1, IsNum: true},
		},
	}

	buildVersion string
)

func init() {
	if buildVersion != "" {
		progVersion.Build = []string{
			buildVersion,
		}
	}
}

func main() {
	// Initialize viper and default options.
	InitConfig()

	// Setup flags.
	flag.String("config", viper.GetString("config"), "Alternate path to the config file.")

	flag.Bool("debug", viper.GetBool("debug"), "Turn on debug output.")
	flag.String("mongo.uri", viper.GetString("mongo.uri"), "URI of the MongoDB host or cluster.")

	flag.String("smtp.host", viper.GetString("smtp.host"), "Host of the SMTP server.")
	flag.Int("smtp.port", viper.GetInt("smtp.port"), "Port of the SMTP server.")
	flag.String("smtp.user", viper.GetString("smtp.user"), "SMTP user.")
	flag.String("smtp.password", viper.GetString("smtp.password"), "SMTP password.")
	flag.String("smtp.from", viper.GetString("smtp.from"), "SMTP From address.")

	flag.Parse()

	// Visit all of the seen flags to update the config.
	// All flag types in flag package support the getter interface.
	flag.Visit(func(f *flag.Flag) {
		viper.Set(f.Name, f.Value.(flag.Getter).Get())
	})

	args := flag.Args()

	if len(args) == 0 {
		PrintUsage("help")
	}

	// Route command.
	switch args[0] {
	case "version":
		versionCmd(args[1:])

	case "put":
		putCmd(args[1:])

	case "get":
		getCmd(args[1:])

	case "keys":
		keysCmd(args[1:])

	case "log":
		logCmd(args[1:])

	case "http":
		httpCmd(args[1:])

	case "config":
		configCmd(args[1:])

	case "subscribe":
		subscribeCmd(args[1:])

	case "unsubscribe":
		unsubscribeCmd(args[1:])

	default:
		// Print usage of speific command.
		if len(args) == 2 {
			PrintUsage(args[1])
		}

		PrintUsage("help")
	}
}
