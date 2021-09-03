package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/woodrufj4/GitVersion/command"
	"github.com/woodrufj4/GitVersion/version"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {

	commands := command.Commands()

	cli := &cli.CLI{
		Name:     "gitversion",
		Version:  version.GetVersionInfo(),
		Args:     args,
		Commands: commands,
	}

	exitCode, err := cli.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	}

	return exitCode

}
