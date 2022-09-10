package faker

import (
	"strings"
)

// Directory is a faker struct for Directory
type Directory struct {
	Faker      *Faker
	OSResolver OSResolver
}

// Directory returns a fake directory path (the directory path style is dependent OS dependent)
func (d Directory) Directory(levels int) string {
	switch d.OSResolver.OS() {
	case "windows":
		return d.WindowsDirectory(levels)
	default:
		return d.UnixDirectory(levels)
	}
}

// UnixDirectory returns a fake Unix directory path, regardless of the host OS
func (d Directory) UnixDirectory(levels int) string {
	return "/" + strings.Join(d.Faker.Lorem().Words(levels), "/")
}

// WindowsDirectory returns a fake Windows directory path, regardless of the host OS
func (d Directory) WindowsDirectory(levels int) string {
	return d.DriveLetter() + strings.Join(d.Faker.Lorem().Words(levels), "\\")
}

// DriveLetter returns a fake Win32 drive letter
func (d Directory) DriveLetter() string {
	return d.Faker.RandomLetter() + ":\\"
}
