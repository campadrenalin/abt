package main

import (
	"fmt"
	"os"

	"github.com/campadrenalin/abt"
)

type Action func(config *Config, args []string) (err error)

type ArgLengthError struct {
	SubCommand string
	Expected   int
	Given      int
}

func (err ArgLengthError) Error() string {
	return fmt.Sprintf(
		"'abt %s' takes %d argument(s), got %d",
		err.SubCommand,
		err.Expected,
		err.Given,
	)
}

func doList(config *Config, args []string) (err error) {

	if len(args) != 1 {
		err = ArgLengthError{"list", 1, len(args)}
		return
	}
	path := args[0]

	abtfile, err := abt.OpenABTFile(&path)
	if err != nil {
		return
	}

	fmt.Printf("File contains %d sections", len(abtfile.Sections))
	if len(abtfile.Sections) > 0 {
		fmt.Printf(":\n")
	} else {
		fmt.Printf(".\n")
	}

	for _, section := range abtfile.Sections {
		fmt.Printf(
			"  %d\t(%d-%d)\t%s",
			section.Filesize,
			section.Start,
			section.Start+section.Size,
			section.Path,
		)
		if section.Origin != "" {
			fmt.Printf(" <- %s", section.Origin)
		}
		fmt.Printf("\n")
	}

	return
}

func doCreate(config *Config, args []string) (err error) {

	if len(args) != 1 {
		err = ArgLengthError{"create", 1, len(args)}
		return
	}
	path := args[0]

	file, err := os.Create(path)
	if err != nil {
		return
	}

	var abtfile abt.ABTFile
	err = abtfile.Write(file)
	return
}

func doInsert(config *Config, args []string) (err error) {

	if len(args) != 2 {
		err = ArgLengthError{"insert", 2, len(args)}
		return
	}
	path, source := args[0], args[1]

	abtfile, err := abt.OpenABTFile(&path)
	if err != nil {
		return
	}

	// Insert metadata
	var sdstruct abt.SectionData
	sdstruct.Path = source
	sdstruct.Origin = source
	abtfile.Sections = append(abtfile.Sections, sdstruct)

	file, err := os.Create(path)
	if err != nil {
		return
	}

	err = abtfile.Write(file)
	return
}
