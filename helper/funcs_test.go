package helper

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/Masterminds/semver"
)

var (
	testSHA = "0b85e68c13650b5ba793c681defe1dc12f03fb50"
)

type execFunc func(name string, arg ...string) *exec.Cmd

func resetExecCommand() {
	execCommand = exec.Command
}

func isHelperProcess() bool {
	return os.Getenv("IS_HELPER_PROCESS") == "1"
}

// Helper function to mock test with exec.Command
//
// Credit: https://npf.io/2015/06/testing-exec-command/
func shellHelper(helperName string) execFunc {
	return func(name string, arg ...string) *exec.Cmd {
		cs := []string{fmt.Sprintf("-test.run=%s", helperName), "--"}
		cs = append(cs, arg...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"IS_HELPER_PROCESS=1"}
		return cmd
	}
}

func TestCurrentCommitSHAHelperProcess(t *testing.T) {
	if !isHelperProcess() {
		return
	}

	fmt.Fprint(os.Stdout, testSHA)
	os.Exit(0)
}

func TestCurrentCommitSHA(t *testing.T) {
	execCommand = shellHelper("TestCurrentCommitSHAHelperProcess")
	defer resetExecCommand()

	sha, err := CurrentCommitSHA()

	if err != nil {
		t.Fatalf("expected not to get an error: %v", err)
	}

	if sha != testSHA {
		t.Fatalf("expected sha to be %s, but was %s", testSHA, sha)
	}
}

func TestTagExistsHelperProcess(t *testing.T) {
	if !isHelperProcess() {
		return
	}

	tags := "testing\n0.1.0-rc.1\nlatest_version"
	fmt.Fprint(os.Stdout, tags)
	os.Exit(0)
}

func TestTagExists(t *testing.T) {

	execCommand = shellHelper("TestTagExistsHelperProcess")
	defer resetExecCommand()

	tagName := "0.1.0-rc.1"

	exists, err := TagExists(tagName)

	if err != nil {
		t.Fatalf("expected not to incounter an error: %v", err)
	}

	if !exists {
		t.Fatalf("expected tag %s to exists, but it does not", tagName)
	}

}

func TestTagNotExists(t *testing.T) {

	execCommand = shellHelper("TestTagExistsHelperProcess")
	defer resetExecCommand()

	tagName := "GitVersion"

	exists, err := TagExists(tagName)

	if err != nil {
		t.Fatalf("expected not to encounter an error: %v", err)
	}

	if exists {
		t.Fatalf("expected tag %s not to exists, but it does", tagName)
	}
}

func TestCurrentCommitTagsHelperProcess(t *testing.T) {
	if !isHelperProcess() {
		return
	}

	tags := "tag: testing, tag: 0.1.0-rc.1, tag: latest_version, HEAD -> develop"
	fmt.Fprint(os.Stdout, tags)
	os.Exit(0)
}

func TestCurrentCommitTags(t *testing.T) {
	execCommand = shellHelper("TestCurrentCommitTagsHelperProcess")
	defer resetExecCommand()

	tags, err := CurrentCommitTags()

	if err != nil {
		t.Fatalf("did not expect to receive an error: %v", err)
	}

	for _, tag := range tags {
		switch tag {
		case "testing", "0.1.0-rc.1", "latest_version":
			continue

		default:
			t.Fatalf("expected tag to be one of the known tags, but was %s", tag)
		}
	}
}

func TestTestHighestSemverTagHelperProcess(t *testing.T) {
	if !isHelperProcess() {
		return
	}

	tags := "testing\n0.1.0-rc.1\nlatest_version\n1.0.0-rc.1\n1.0.1\n1.0.13\nanother_test"
	fmt.Fprint(os.Stdout, tags)
	os.Exit(0)
}

func TestHighestSemverTag(t *testing.T) {
	execCommand = shellHelper("TestTestHighestSemverTagHelperProcess")
	defer resetExecCommand()

	benchmark, err := semver.NewVersion("1.0.13")

	if err != nil {
		t.Fatalf("should not error for benchmark version: %v", err)
	}

	highest, err := HighestSemverTag()

	if err != nil {
		t.Fatalf("encountered an error when retrieving highest semver tag: %v", err)
	}

	if !benchmark.Equal(highest) {
		t.Fatalf("expected the highest semver tag to be %s, but was %s", benchmark.String(), highest.String())
	}

}

func TestCurrentBranchNameHelperProcess(t *testing.T) {
	if !isHelperProcess() {
		return
	}

	tags := "develop"
	fmt.Fprint(os.Stdout, tags)
	os.Exit(0)
}

func TestCurrentBranchName(t *testing.T) {

	execCommand = shellHelper("TestCurrentBranchNameHelperProcess")
	defer resetExecCommand()

	branchName, err := CurrentBranchName()

	if err != nil {
		t.Fatalf("did not expect to encounter an error: %v", err)
	}

	expectedBranchName := "develop"

	if branchName != expectedBranchName {
		t.Fatalf("expected branch name to be %s, but was %s", expectedBranchName, branchName)
	}
}

func TestCurrentCommitMessageHelperProcess(t *testing.T) {
	if !isHelperProcess() {
		return
	}

	msg := "Initial commit\n\nThere could be multiple lines for a given commit message."
	fmt.Fprint(os.Stdout, msg)
	os.Exit(0)
}

func TestCurrentCommitMessage(t *testing.T) {
	execCommand = shellHelper("TestCurrentCommitMessageHelperProcess")
	defer resetExecCommand()

	message, err := CurrentCommitMessage()

	if err != nil {
		t.Fatalf("should not encounter an error: %v", err)
	}

	msg := "Initial commit\n\nThere could be multiple lines for a given commit message."

	if message != msg {
		t.Fatalf("expected commit message to be:\n\n%s\n\nbut was:\n\n%s", msg, message)
	}
}
