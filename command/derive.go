package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/go-hclog"
	"github.com/woodrufj4/GitVersion/helper"
	"github.com/woodrufj4/GitVersion/version"
)

type DeriveCommand struct {
	args   []string
	Logger hclog.Logger
	config *DeriveConfig
}

type DeriveConfig struct {
	PrettyPrint bool
}

func (d *DeriveCommand) Help() string {
	helpText := `
Usage: gitversion derive [options]

  Derives a semantic version (SemVer) from the current commit or branch name.
  If unable to derive a semantic version, it will default to 0.1.0 and flag the
  version as 'isDefault'.

  If the default version is provided, it will append a pre-release of 'dev' if
  the current branch is not 'main', or 'master'.

  Options:

    --pretty
      Format the output in pretty print JSON.

`
	return helpText
}

func (d *DeriveCommand) Synopsis() string {
	return "Derives a semantic version (SemVer) from Git"
}

func (d *DeriveCommand) OutputVersionDetail(input *version.VersionDetailInput) int {

	detail, err := version.NewVersionDetail(input)

	if err != nil {
		d.Logger.Error("unable to create version detail", "error", err.Error())
		return 1
	}

	var bytes []byte

	if d.config.PrettyPrint {
		bytes, err = json.MarshalIndent(detail, "", "  ")

	} else {
		bytes, err = json.Marshal(detail)
	}

	if err != nil {
		d.Logger.Error("unable to convert version detail to JSON format", "error", err.Error())
		return 1
	}

	_, err = fmt.Print(string(bytes))

	if err != nil {
		d.Logger.Error("unable to output version detail", "error", err.Error())
		return 1
	}

	return 0

}

func (d *DeriveCommand) parseFlags() error {

	fs := flag.NewFlagSet("derive command", flag.ContinueOnError)

	fs.Usage = func() {
		fmt.Println(d.Help())
	}

	config := &DeriveConfig{}

	fs.BoolVar(&config.PrettyPrint, "pretty", false, "")

	err := fs.Parse(d.args)

	if err != nil {
		return err
	}

	d.config = config

	return nil
}

// Run executes the derive command to discern the tags from the current commit or branch name.
func (d *DeriveCommand) Run(args []string) int {

	d.args = args

	err := d.parseFlags()

	if err != nil {
		d.Logger.Error("unable to parse command arguments", "error", err.Error())
		return 1
	}

	currentBranchName, err := helper.CurrentBranchName()

	if err != nil {
		d.Logger.Error("unable to retrieve the current branch name", "error", err.Error())
		return 1
	}

	currentSHA, err := helper.CurrentCommitSHA()

	if err != nil {
		d.Logger.Error("unable to retrieve the current commit SHA", "error", err.Error())
		return 1
	}

	// Check to see if the current commit has a tag and is a valid SemVer tag
	currentTags, err := helper.CurrentCommitTags()

	if err != nil {
		d.Logger.Error("unable to retrieve the current commit tags", "error", err.Error())
		return 1
	}

	versionDetailInput := &version.VersionDetailInput{
		SHA:        currentSHA,
		BranchName: currentBranchName,
	}

	for _, tag := range currentTags {
		// Check if the tag is valid Semver format
		tagVersion, err := semver.NewVersion(tag)

		if err != nil {
			continue
		}

		versionDetailInput.Semver = tagVersion
		return d.OutputVersionDetail(versionDetailInput)

	}

	// None of the current tags were valid SemVer format, let's get the tag from
	// the release branch name, if applicable
	if strings.HasPrefix(currentBranchName, "release/") || strings.HasPrefix(currentBranchName, "release-") {

		branchVersionString := ""

		if strings.HasPrefix(currentBranchName, "release/") {
			branchVersionString = strings.TrimPrefix(currentBranchName, "release/")
		} else {
			branchVersionString = strings.TrimPrefix(currentBranchName, "release-")
		}

		branchVersion, err := semver.NewVersion(branchVersionString)

		if err == nil {
			versionDetailInput.Semver = branchVersion
			return d.OutputVersionDetail(versionDetailInput)
		}

	}

	// Default to version 0.1.0
	defaultVersion, err := semver.NewVersion("0.1.0")

	if err != nil {
		d.Logger.Error("unable generate default version", "error", err.Error())
		return 1
	}

	if currentBranchName != "master" && currentBranchName != "main" {
		defaultVersion.SetPrerelease("dev")
	}

	versionDetailInput.Semver = defaultVersion
	versionDetailInput.IsDefault = true

	return d.OutputVersionDetail(versionDetailInput)

}
