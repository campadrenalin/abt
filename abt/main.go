package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/campadrenalin/abt"
)

type Action struct {
	verbose bool
	help    bool
}

func doAction(action *Action, path *string) (err error) {

	abtfile, err := abt.OpenABTFile(path)
	if err != nil {
		return
	}

	fmt.Printf("%#v\n", *abtfile)
	return
}

func main() {
	var action Action

	flagset := flag.NewFlagSet("abt", flag.PanicOnError)
	flagset.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\tabt [ FILE... ]\n")
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flagset.PrintDefaults()
	}

	flagset.BoolVar(&action.verbose, "v", false, "Verbose output")
	flagset.BoolVar(&action.help, "help", false, "Show this message")

	flagset.Parse(os.Args[1:])

	if action.help {
		flagset.Usage()
		return
	}
	for _, value := range flagset.Args() {
		err := doAction(&action, &value)
        if err != nil {
            fmt.Printf(
                "Could not process file '%s'\n%s\n",
                value,
                err.Error(),
            )
        }
	}
}
