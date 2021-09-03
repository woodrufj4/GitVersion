package version

import (
	"testing"

	"github.com/Masterminds/semver"
)

func TestNewVersionDetail(t *testing.T) {

	version, err := semver.NewVersion("0.1.0-rc.1+build.1234")

	if err != nil {
		t.Fatalf("should not fail valid semantic version: %v", err)
	}

	input := &VersionDetailInput{
		Semver:     version,
		IsDefault:  false,
		BranchName: "release/0.1.0",
		SHA:        "0b85e68c13650b5ba793c681defe1dc12f03fb50",
	}

	versionDetail, err := NewVersionDetail(input)

	if err != nil {
		t.Fatalf("failed to create a new version detail: %v", err)
	}

	if versionDetail.BranchName != input.BranchName {
		t.Fatalf("expected branch name to be %s, but was %s", input.BranchName, versionDetail.BranchName)
	}

	if versionDetail.SHA != input.SHA {
		t.Fatalf("expected commit SHA to be %s, but was %s", input.SHA, versionDetail.SHA)
	}

	if version.String() != versionDetail.SemVer {
		t.Fatalf("expected version detail semver format to be %s, but was %s", version.String(), versionDetail.SemVer)
	}

	if version.Major() != versionDetail.Major {
		t.Fatalf("expected the major version to be %d, but was %d", version.Major(), versionDetail.Major)
	}

	if version.Minor() != versionDetail.Minor {
		t.Fatalf("expected the minor version to be %d, but was %d", version.Minor(), versionDetail.Minor)
	}

	if version.Patch() != versionDetail.Patch {
		t.Fatalf("expected the patch version to be %d, but was %d", version.Patch(), versionDetail.Patch)
	}

	if version.Prerelease() != versionDetail.PreRelease {
		t.Fatalf("expected the pre-release version to be %s, but was %s", version.Prerelease(), versionDetail.PreRelease)
	}

	if version.Metadata() != versionDetail.Metadata {
		t.Fatalf("expected the metadata info to be %s, but was %s", version.Metadata(), versionDetail.Metadata)
	}

	if _, err := semver.NewVersion(versionDetail.Core); err != nil {
		t.Fatalf("expected the version detail core, to be a valid semver format: %v", err)
	}
}

func TestMissingVersionError(t *testing.T) {
	_, err := NewVersionDetail(nil)

	if err != ErrorMissingVersionInput {
		t.Fatalf("expected error to be %v, but was %v", ErrorMissingVersionInput, err)
	}
}

func TestMissingSemverError(t *testing.T) {

	input := &VersionDetailInput{
		IsDefault:  false,
		BranchName: "release/0.1.0",
		SHA:        "0b85e68c13650b5ba793c681defe1dc12f03fb50",
	}

	_, err := NewVersionDetail(input)

	if err != ErrorMissingSemverInput {
		t.Fatalf("expected error to be %v, but was %v", ErrorMissingSemverInput, err)
	}

}
