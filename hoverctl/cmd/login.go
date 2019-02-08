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

			hostFlag, err := cmd.Flags().GetString("host")
			handleIfError(err)
			adminPortFlag, err := cmd.Flags().GetInt("admin-port")
			handleIfError(err)
			proxyPortFlag, err := cmd.Flags().GetInt("proxy-port")
			handleIfError(err)

			target = configuration.NewTarget(newTargetFlag, hostFlag, adminPortFlag, proxyPortFlag)
		}

		if username == "" {
			username = askForInput("Username", false)
		}
		if password == "" {
			password = askForInput("Password", true)
		}

		if username == "" || password == "" {
			handleIfError(errors.New("missing username or password"))
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
	loginCmd.Flags().Int("admin-port", 0, "A port number for the Hoverfly API/GUI. Overrides the default Hoverfly admin port (8888)")
	loginCmd.Flags().Int("proxy-port", 0, "A port number for the Hoverfly proxy. Overrides the default Hoverfly proxy port (8500)")
	loginCmd.Flags().String("host", "", "A host on which a Hoverfly instance is running. Overrides the default Hoverfly host (localhost). HTTP protocol is assumed if scheme is not specified.")
	loginCmd.Flags().StringVar(&username, "username", "", "Username to authenticate against Hoverfly with")
	loginCmd.Flags().StringVar(&password, "password", "", "Password to authenticate against Hoverfly with")
}
