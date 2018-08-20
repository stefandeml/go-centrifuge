package version

import (
	"github.com/Masterminds/semver"
	"fmt"
)

// Git SHA1 commit hash of the release (set via linker flags - ldflags)
var gitCommit = "master"

// CentrifugeNodeVersion is the current version of the app
const centrifugeNodeVersion = "0.0.1-pre-alpha"

func GetVersion() *semver.Version {

	v, err := semver.NewVersion(fmt.Sprintf("%s+%s", centrifugeNodeVersion, gitCommit))
	if err != nil {
		log.Panicf("Invalid CentrifugeNodeVersion specified: %s", centrifugeNodeVersion)
	}

	return v

}