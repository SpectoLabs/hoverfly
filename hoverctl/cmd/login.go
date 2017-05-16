package cmd

import (
	"errors"
	"fmt"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	"github.com/spf13/cobra"
)

var username, password string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Hoverfly",
	Long: `
Will authenticate against the /api/token-auth API endpoint 
target Hoverfly instance using the provided username and 
password.

The generated authentication token is then stored on the
target in the hoverctl configuration file.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		newTargetFlag, _ := cmd.Flags().GetString("new-target")

		if newTargetFlag != "" {
			if config.GetTarget(newTargetFlag) != nil {
				handleIfError(fmt.Errorf("Target %s already exists\n\nUse a different target name or run `hoverctl targets update %[1]s`", newTargetFlag))
			}

			// If the host is set to a remote instance, the default HTTPS port
			// is used instead of 8888 which is set in wrapper.NewTarget()
			if adminPortFlag == 0 && (hostFlag != "" && !wrapper.IsLocal(hostFlag)) {
				adminPortFlag = 443
			}

			target = configuration.NewTarget(newTargetFlag, hostFlag, adminPortFlag, proxyPortFlag)
		}

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
			handleIfError(err)
		}

		target.AuthToken = token

		config.NewTarget(*target)
		handleIfError(config.WriteToFile(hoverflyDirectory))

		fmt.Println("Login successful")
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)

	loginCmd.Flags().String("new-target", "", "A name for a new target that hoverctl will create and associate the Hoverfly instance to")

	loginCmd.Flags().StringVar(&username, "username", "", "Username to authenticate against Hoverfly with")
	loginCmd.Flags().StringVar(&password, "password", "", "Password to autenticate against Hoverfly with")
}
