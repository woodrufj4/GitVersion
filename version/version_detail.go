package version

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver"
)

var (
	ErrorMissingVersionInput = errors.New("missing input detail")
	ErrorMissingSemverInput  = errors.New("missing semver detail")
)

type VersionDetail struct {
	Major      int64  `json:"major"`
	Minor      int64  `json:"minor"`
	Patch      int64  `json:"patch"`
	Core       string `json:"core"`
	IsDefault  bool   `json:"isDefault"`
	PreRelease string `json:"preRelease"`
	Metadata   string `json:"metadata"`
	SemVer     string `json:"semver"`
	SHA        string `json:"sha"`
	ShortSHA   string `json:"shortSha"`
	BranchName string `json:"branchName"`
}

type VersionDetailInput struct {
	Semver     *semver.Version
	IsDefault  bool
	SHA        string
	BranchName string
}

func NewVersionDetail(input *VersionDetailInput) (*VersionDetail, error) {

	if input == nil {
		return nil, ErrorMissingVersionInput
	}

	if input.Semver == nil {
		return nil, ErrorMissingSemverInput
	}

	shortSha := ""

	if len(input.SHA) >= 7 {
		shortSha = input.SHA[0:7]
	}

	return &VersionDetail{
		Major:      input.Semver.Major(),
		Minor:      input.Semver.Minor(),
		Patch:      input.Semver.Patch(),
		Core:       fmt.Sprintf("%d.%d.%d", input.Semver.Major(), input.Semver.Minor(), input.Semver.Patch()),
		PreRelease: input.Semver.Prerelease(),
		Metadata:   input.Semver.Metadata(),
		SemVer:     input.Semver.String(),
		BranchName: input.BranchName,
		SHA:        input.SHA,
		ShortSHA:   shortSha,
		IsDefault:  input.IsDefault,
	}, nil
}
