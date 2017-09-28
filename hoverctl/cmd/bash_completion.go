package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
)

var completionCmd = &cobra.Command{
	Use:   "completion [destination for symbolic link (optional)]",
	Short: "Create Bash completion file for hoverctl",
	Long: `
Create a symbolic link to the Bash completion file in the specified location.
If you do not specify a location hoverctl will suggest the most probable location
based on your operating system.
`,

	Run: func(cmd *cobra.Command, args []string) {

		completionFilePath := path.Join(hoverflyDirectory.Path, "hoverctl")

		symlinkToCompletionFile, err := determineSymlinkLocation(args)
		if err != nil {
			handleIfError(err)
		}

		if !askForConfirmation(fmt.Sprintf("Are you sure you want to create the completion file in %q and symlink it to %q ?", completionFilePath, symlinkToCompletionFile)) {
			return
		}

		errCompletionFile := createCompletionFile(completionFilePath)
		if errCompletionFile != nil {
			handleIfError(errCompletionFile)
		}

		errSymlink := createSymlink(completionFilePath, symlinkToCompletionFile)
		if errSymlink != nil {
			handleIfError(errSymlink)
		}

		fmt.Println("Completion file and symbolic link created. Restart your shell to activate.")
	},
}

func determineSymlinkLocation(args []string) (string, error) {
	var symlinkToCompletionFile string
	var err error

	//optionally override the symlink location from command line
	if len(args) > 0 {
		symlinkToCompletionFile, err = expandHomeDirectory(args[0])
		if err != nil {
			return "", err
		}
	} else {
		symlinkToCompletionFile = locateProbableCompletionFileLocation()
	}
	return symlinkToCompletionFile, nil
}

func createSymlink(completionFilePath string, symlinkToCompletionFile string) error {
	err := os.Symlink(completionFilePath, symlinkToCompletionFile)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Created symbolic link [%q=>%q]\n", completionFilePath, symlinkToCompletionFile)
	}

	return nil
}

func createCompletionFile(completionFilePath string) error {
	if err := RootCmd.GenBashCompletionFile(completionFilePath); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Created completion file in %q\n", completionFilePath)
	}
	return nil
}

func locateProbableCompletionFileLocation() string {
	var completionFileLocation string
	if runtime.GOOS == "darwin" {
		completionFileLocation = path.Join(string(os.PathSeparator)+"usr", "local", "etc", "bash_completion.d", "hoverctl")
	} else {
		completionFileLocation = path.Join(string(os.PathSeparator)+"etc", "bash_completion.d", "hoverctl")
	}
	return completionFileLocation
}

func init() {
	RootCmd.AddCommand(completionCmd)
}

func expandHomeDirectory(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}
