package version

import "fmt"

// These variables get set by the build script via the LDFLAGS.
var (
	VersionPrerelease string
	VersionMetadata   string
	VersionCore       string = "0.1.0"
)

// GetVersionInfo reports the semantic version (SemVer) of the GitVersion.
// This the internal version information for the GitVersion tool.
func GetVersionInfo() string {

	version := VersionCore

	if VersionPrerelease != "" {
		version = fmt.Sprintf("%s-%s", version, VersionPrerelease)
	}

	if VersionMetadata != "" {
		version = fmt.Sprintf("%s+%s", version, VersionMetadata)
	}

	return version

}
