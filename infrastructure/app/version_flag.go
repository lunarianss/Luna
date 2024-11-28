// Copyright 2024 Benjamin Lee <cyan0908@163.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package app

// Package verflag defines utility functions to handle command line flags
// related to version of IAM.

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	flag "github.com/spf13/pflag"
)

type versionValue int

var (
	// GitVersion is semantic version.
	GitVersion = "v0.0.0-master+$Format:%h$"
	// BuildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	BuildDate = "1970-01-01T00:00:00Z"
	// GitCommit sha1 from git, output of $(git rev-parse HEAD).
	GitCommit = "$Format:%H$"
	// GitTreeState state of git tree, either "clean" or "dirty".
	GitTreeState = "clean"
)

// Define some const.
const (
	VersionFalse versionValue = 0
	VersionTrue  versionValue = 1
	VersionRaw   versionValue = 2
)

const strRawVersion string = "raw"

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw

		return nil
	}

	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}

	return err
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}

	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// The type of the flag as required by the pflag.Value interface.
func (v *versionValue) Type() string {
	return "version"
}

// VersionVar defines a flag with the specified name and usage string.
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	flag.VarP(p, name, "v", usage)
	// "--version" will be treated as "--version=true"
	flag.Lookup(name).NoOptDefVal = "true"
}

// Version wraps the VersionVar function.
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)

	return p
}

const versionFlagName = "version"

var versionFlag = Version(versionFlagName, VersionFalse, "Print version information and quit.")

// AddFlags registers this package's flags on arbitrary FlagSets, such that they point to the
// same value as the global flags.
func AddVersionFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(versionFlagName))
}

// Info contains versioning information.
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	if s, err := info.Text(); err == nil {
		return string(s)
	}

	return info.GitVersion
}

// nolint: errchkjson
// ToJSON returns the JSON string of version information.
func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)

	return string(s)
}

func (info Info) tableColorKeyStr(tag string) string {
	return color.BlueString(tag)
}

func (info Info) tableColorValueStr(value string) string {
	return fmt.Sprintf(": %s", value)
}

// Text encodes the version information into UTF-8-encoded text and
// returns the result.
func (info Info) Text() ([]byte, error) {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow(info.tableColorKeyStr("gitVersion"), info.tableColorValueStr(info.GitVersion))
	table.AddRow(info.tableColorKeyStr("gitCommit"), info.tableColorValueStr(info.GitCommit))
	table.AddRow(info.tableColorKeyStr("gitTreeState"), info.tableColorValueStr(info.GitTreeState))
	table.AddRow(info.tableColorKeyStr("buildDate"), info.tableColorValueStr(info.BuildDate))
	table.AddRow(info.tableColorKeyStr("goVersion"), info.tableColorValueStr(info.GoVersion))
	table.AddRow(info.tableColorKeyStr("compiler"), info.tableColorValueStr(info.Compiler))
	table.AddRow(info.tableColorKeyStr("platform"), info.tableColorValueStr(info.Platform))

	return table.Bytes(), nil
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func GetVersionInfo() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in pkg/version/base.go
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
