package helper

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
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
