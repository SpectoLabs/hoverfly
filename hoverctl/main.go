package main

import "github.com/SpectoLabs/hoverfly/v2/hoverctl/cmd"

var (
	hoverctlVersion = "0.0.1"
)

func main() {
	cmd.Execute(hoverctlVersion)
}
