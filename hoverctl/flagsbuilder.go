package main

type FlagsBuilder struct {
	Webserver   string
	Certificate string
	Key         string
}

func (this FlagsBuilder) BuildFlags() []string {
	flags := []string{}
	if this.Webserver == "webserver" {
		flags = append(flags, "-webserver")
	}
	if this.Certificate != "" {
		flags = append(flags, "-cert="+this.Certificate)
	}
	if this.Key != "" {
		flags = append(flags, "-key="+this.Key)
	}

	return flags
}
