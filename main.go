package main

import (
	"elastic-search/cmd"
	"elastic-search/global"
)

var (
	major = "0"
	minor = "0"
	patch = "0"
)

func main() {
	cmd.Execute()
}

func init() {
	global.CurrentVersion.Major = major
	global.CurrentVersion.Minor = minor
	global.CurrentVersion.Patch = patch
}
