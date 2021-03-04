package version

import (
	"fmt"
	"time"
)

// APPNAME is the app-global application name string, which should be substituted with a real value during build.
var APPNAME = "UNKNOWN"

// VERSION is the app-global version string, which should be substituted with a real value during build.
var VERSION = "UNKNOWN"

// GOVERSION is the Golang version used to generate the binary.
var GOVERSION = "UNKNOWN"

// BUILDTIME is the timestamp at which the binary was created.
var BUILDTIME = "UNKNOWN"

// COMMITHASH is the git commit hash that was used to generate the binary.
var COMMITHASH = "UNKNOWN"

// Version of the application
func Version() string {
	return fmt.Sprintf("%s %s (commit %s built with go %s the %s)", APPNAME, VERSION, COMMITHASH, GOVERSION, BUILDTIME)
}

// Compiled swill transform the BUILDTIME constant into a valid time.Time (defaults to time.Now())
func Compiled() time.Time {
	if BUILDTIME == "UNKNOWN" {
		return time.Now()
	}

	t, err := time.Parse("2006-02-15T15:04:05Z+00:00", BUILDTIME)
	if err != nil {
		return time.Now()
	}

	return t
}
