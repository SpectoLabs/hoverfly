package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var username, password string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Hoverfly",
	Long: `
Login to Hoverfly"
	`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		if username == "" {
			username = askForInput("Username", false)
		}
		if password == "" {
			password = askForInput("Password", true)
		}

		if username == "" || password == "" {
			handleIfError(errors.New("Missing username or password"))
		}

		token, err := wrapper.Login(*target, username, password)
		if err != nil {
			if verbose {
				fmt.Fprintln(os.Stderr, err.Error())
			}

			handleIfError(errors.New("Failed to login to Hoverfly"))
		}

		target.AuthToken = token

		config.NewTarget(*target)
		handleIfError(config.WriteToFile(hoverflyDirectory))

		fmt.Println("Login successful")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVar(&username, "username", "", "Username to login to Hoverfly")
	loginCmd.Flags().StringVar(&password, "password", "", "Password to login to Hoverfly")
}
