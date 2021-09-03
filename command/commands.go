package command

import (
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
)

func Commands() map[string]cli.CommandFactory {

	logger := hclog.Default()

	commands := map[string]cli.CommandFactory{

		"derive": func() (cli.Command, error) {
			return &DeriveCommand{
				Logger: logger.Named("derive"),
			}, nil
		},
	}

	return commands
}
