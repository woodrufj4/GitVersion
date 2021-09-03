package version

import (
	"testing"

	"github.com/Masterminds/semver"
)

func TestGetVersionInfo(t *testing.T) {

	VersionCore = "0.1.0"
	VersionPrerelease = "dev.1"
	VersionMetadata = "build.1234"

	version := GetVersionInfo()

	_, err := semver.NewVersion(version)

	if err != nil {
		t.Fatalf("expected gitversion's version to be semantically formatted: %v", err)
	}

}
