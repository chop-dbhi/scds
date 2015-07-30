package main

import (
	"fmt"
	"os"
)

var defaultUsage = `scds <cmd> [options...]

SCD (slowly-changing dimensions) Store (scds) provides an interface for storing
the state of an object and seeing how it changed from the last state. The read and
write interface is key-value based where the value is a valid JSON document.

Commands:

	help		Prints the usage information.
	config		Prints all the configuration options.
	put			Puts an object in the store.
	get			Gets the latest state of an object from the store.
	keys		Returns a list of keys in the store.
	log			Returns an ordered set of diffs for an object.
	http		Runs an HTTP service with a comparable set of commands.
	subscribe	Subscribes one or more emails to receive notifications.
	unsubscribe	Unsubscribes one or more emails from receiving notifications.

Global Options:

	-debug	Turn on debug output.

	-mongo.uri <uri>	Specify one or more MongoDB hosts [default: localhost/scds].

	-smtp.host <host>		Host of the SMTP server [default: localhost].
	-smtp.port <port>		Port of the SMTP server [default: 25].
	-smtp.user <user>		User to authenticate with the SMTP server.
	-smtp.password <pass>	Password to authenticate with the SMTP server.
	-smtp.from <from>		From email address.

Run 'sdcs help <cmd>' to get help about a specific command.
`

var putUsage = `scds put <key> <object>

Puts an object into the store. If the object does not exist, it will create
it, otherwise it will compare it with the existing state.
`

var getUsage = `scds get <key>

Gets the current state of an object if it exists.

Options:

	-version <int>	Gets the state at a specific version.
	-time <time>	Gets the state at the specified time (Unix timestamp).
`

var keysUsage = `scds keys

Gets a list of keys in the store.
`

var logUsage = `scds log <key>

Returns an ordered set of diffs for the object making up the log.
`

var httpUsage = `scds http [--host=<host>] [--port=<port>]

Runs an HTTP server that defines endpoints corresponding to the command-line
interface (CLI).

Endpoints:

	GET /keys						Returns a list keys in the store.

	PUT /objects/:key				Puts an object in the store.
	GET /objects/:key				Gets the latest state of an object from the store.
	GET /objects/:key/v/:version	Gets the state of an object at the specified version.
	GET /objects/:key/t/:time		Gets the state of an object at the specified time.

	GET /log/:key					Returns an ordered set of diffs for an object.

Options:

	-host <host>	The host to bind the HTTP server to [default: localhost].
	-port <port>	The port to bind the HTTP server to [default: 5000].

`

var configUsage = `scds config

Prints the configuration options defined across the configuration file, environment
variables, and command-line flags.
`

var subscribeUsage = `scds subscribe email [emails...]

Subscribes one or more email addresses to receive notifications. Email
addresses that already subscribed will not be subscribed again.
`

var unsubscribeUsage = `scds unsubscribe email [emails...]

Unsubscribes one or more email addresses from receiving notifications. Emails
that are not subscribes will be ignored.
`

func PrintUsage(cmd string) {
	var usage string

	switch cmd {
	case "get":
		usage = getUsage

	case "keys":
		usage = keysUsage

	case "put":
		usage = putUsage

	case "log":
		usage = logUsage

	case "http":
		usage = httpUsage

	case "config":
		usage = configUsage

	case "subscribe":
		usage = subscribeUsage

	case "unsubscribe":
		usage = unsubscribeUsage

	default:
		usage = defaultUsage
	}

	fmt.Println(usage)
	os.Exit(1)
}
