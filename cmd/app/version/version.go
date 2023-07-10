package version

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/carlmjohnson/versioninfo"
	"github.com/prometheus/common/version"
)

var (
	Version  = versioninfo.Version
	Revision = versioninfo.Revision
)

type VersionFlag string

func (v VersionFlag) Decode(ctx *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                         { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	fmt.Printf("version: %s\n", version.Version)
	fmt.Printf("revision: %s\n", version.Revision)
	app.Exit(0)
	return nil
}
