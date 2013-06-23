package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
    command string
    action  Action

	verbose bool
	help    bool
}

func main() {
	var config Config

	flagset := flag.NewFlagSet("abt", flag.PanicOnError)
	flagset.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\tabt command [ FILE... ]\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flagset.PrintDefaults()
	}

	flagset.BoolVar(&config.verbose, "v", false, "Verbose output")
	flagset.BoolVar(&config.help, "help", false, "Show this message")

	flagset.Parse(os.Args[1:])

	if config.help {
		flagset.Usage()
		return
	}

    if flagset.NArg() < 1 {
        fmt.Fprintf(os.Stderr, "abt must be run with a command.\n\n")
        flagset.Usage()
        return
    }

    config.command = flagset.Args()[0]
    actions := map[string] Action {
        "list": doList,
        "create": doCreate,
    }
    config.action = actions[config.command]
    if config.action == nil {
        fmt.Fprintf(os.Stderr, "'%s' is not a legit command.\n\n", config.command)
        flagset.Usage()
        return
    }

	for _, value := range flagset.Args()[1:] {
		err := config.action(&config, &value)
        if err != nil {
            fmt.Fprintf(
                os.Stderr,
                "Could not process file '%s'\n%s\n",
                value,
                err.Error(),
            )
        }
	}
}
