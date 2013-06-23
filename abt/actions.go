package main

import (
	"fmt"
	"os"

	"github.com/campadrenalin/abt"
)

type Action func (config *Config, path *string) (err error)

func doList(config *Config, path *string) (err error) {

	abtfile, err := abt.OpenABTFile(path)
	if err != nil {
		return
	}

	fmt.Printf("%#v\n", *abtfile)

    return
}

func doCreate(config *Config, path *string) (err error) {

    file, err := os.Create(*path)
    if err != nil {
        return
    }

    var data abt.ABTFile
    err = data.Write(file)
    return
}
