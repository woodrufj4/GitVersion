package helper

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
)

// execCommand is a drop in replacement for testing purposes.
var execCommand = exec.Command

// CurrentCommitSHA reports the current commit SHA
func CurrentCommitSHA() (string, error) {

	cmd := execCommand("git", "log", "-n", "1", "--format=%H")
	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil

}

// TagExists checks if the provided tag exists on any commit.
func TagExists(tag string) (bool, error) {

	cmd1 := execCommand("git", "tag", "--sort=-v:refname")
	out, err := cmd1.Output()

	if err != nil {
		return false, err
	}

	tags := strings.Split(string(out), "\n")

	for _, existingTag := range tags {
		if existingTag == tag {
			return true, nil
		}
	}

	return false, nil

}

// NumCommitsSinceTag retrieves the total number of commits since
// the last known tag.
//
// If there aren't any tags, then it retrieves the total number of commits
// in history.
func NumCommitsSinceTag() (int, error) {

	cmd1 := execCommand("git", "tags", "-n", "1", "--sort=-v:refname")
	out, err := cmd1.Output()

	lastTag := "HEAD"

	if err == nil {
		lastTag = string(out)
	}

	// There are no tags, so let's get the total number of commits.
	cmd2 := execCommand("git", "rev-list", lastTag, "--count")

	countOut, err := cmd2.Output()

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(string(countOut)))

}

// CurrentCommitTag retrieves the tag names from the current / HEAD commit.
//
// This will return empty, if the current commit is not tagged.
func CurrentCommitTags() ([]string, error) {

	tags := make([]string, 0)

	cmd := execCommand("git", "log", "-n", "1", "--tags", "--format=%D")
	out, err := cmd.Output()

	if err != nil {
		return tags, err
	}

	refs := strings.Split(string(out), ",")

	for _, refName := range refs {
		trimmedName := strings.TrimSpace(refName)

		if strings.HasPrefix(trimmedName, "tag: ") {
			tags = append(tags, strings.TrimPrefix(trimmedName, "tag: "))
		}
	}

	return tags, nil
}

// HighestSemverTag retrieves the highest semantic versioned tag
func HighestSemverTag() (*semver.Version, error) {

	cmd1 := execCommand("git", "tags", "--sort=-v:refname")
	out, err := cmd1.Output()

	if err != nil {
		return nil, err
	}

	tags := strings.Split(string(out), "\n")

	highestVersion := &semver.Version{}

	for _, tag := range tags {

		tagVersion, err := semver.NewVersion(tag)

		if err != nil {
			continue
		}

		if tagVersion.GreaterThan(highestVersion) {
			highestVersion = tagVersion
		}
	}

	return highestVersion, nil

}

// CurrentBranchName reports the name of the current branch.
func CurrentBranchName() (string, error) {
	cmd := execCommand("git", "branch", "--show-current")
	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

// CurrentCommitMessage retrieves the last commit's message subject (first line).
func CurrentCommitSubject() (string, error) {
	cmd := execCommand("git", "log", "-n", "1", "--format=%s")
	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
