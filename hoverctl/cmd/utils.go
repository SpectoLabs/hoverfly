package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
)

func handleIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func checkArgAndExit(args []string, message, command string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, message)
		fmt.Fprintln(os.Stderr, "\nTry hoverctl "+command+" --help for more information")
		os.Exit(1)
	}
}

func checkTargetAndExit(target *configuration.Target) {
	if target == nil {
		handleIfError(fmt.Errorf("%[1]s is not a target\n\nRun `hoverctl targets create %[1]s`", targetNameFlag))
	}
}

func askForConfirmation(message string) bool {
	if force {
		return true
	}

	for {
		response := askForInput(message+" [y/n]", false)

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func askForInput(value string, sensitive bool) string {
	if force {
		return ""
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s: ", value)
		if sensitive {
			responseBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			handleIfError(err)
			fmt.Println("")

			return strings.TrimSpace(string(responseBytes))
		} else {
			response, err := reader.ReadString('\n')
			handleIfError(err)

			return strings.TrimSpace(response)
		}
	}
}

func drawTable(data [][]string, header bool) {
	table := tablewriter.NewWriter(os.Stdout)
	if header {
		table.SetHeader(data[0])
		data = data[1:]
	}

	for _, v := range data {
		table.Append(v)
	}
	fmt.Print("\n")
	table.Render()
}

func getContentShorthand(content string) string {
	if content != "" {
		scriptArr := strings.Split(content, "\n")
		if verbose || len(scriptArr) < 5 {
			return content
		} else {
			return fmt.Sprintln(scriptArr[0], scriptArr[1], scriptArr[2], scriptArr[3], scriptArr[4], "...")
		}
	}
	return ""
}
