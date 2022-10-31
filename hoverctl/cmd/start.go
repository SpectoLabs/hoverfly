package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	"github.com/SpectoLabs/hoverfly/hoverctl/wrapper"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Hoverfly",
	Long: `
Starts an instance of Hoverfly using the current hoverctl
target configuration.

To start an instance of Hoverfly as a webserver, add the 
argument "webserver" to the start command.

The Hoverfly process ID is stored against the target in the
hoverctl configuration file.
`,
	Run: func(cmd *cobra.Command, args []string) {
		checkTargetAndExit(target)

		if !wrapper.IsLocal(target.Host) {
			handleIfError(fmt.Errorf("Unable to start an instance of Hoverfly on a remote host (%s host: %s)\n\nRun `hoverctl start --new-target <name>`", target.Name, target.Host))
		}

		if wrapper.CheckIfRunning(*target) == nil {
			if _, err := wrapper.GetMode(*target); err == nil {
				handleIfError(fmt.Errorf("Target Hoverfly is already running \n\nRun `hoverctl stop -t %s` to stop it", target.Name))
			}
		}

		newTargetFlag, _ := cmd.Flags().GetString("new-target")

		adminPortFlag, err := cmd.Flags().GetInt("admin-port")
		handleIfError(err)
		proxyPortFlag, err := cmd.Flags().GetInt("proxy-port")
		handleIfError(err)

		if newTargetFlag != "" {
			if config.GetTarget(newTargetFlag) != nil {
				handleIfError(fmt.Errorf("Target %s already exists\n\nUse a different target name or run `hoverctl targets update %[1]s`", newTargetFlag))
			}
			hostFlag, err := cmd.Flags().GetString("host")
			handleIfError(err)
			target = configuration.NewTarget(newTargetFlag, hostFlag, adminPortFlag, proxyPortFlag)
		}

		if adminPortFlag != 0 {
			target.AdminPort = adminPortFlag
		}

		if proxyPortFlag != 0 {
			target.ProxyPort = proxyPortFlag
		}

		target.Webserver = len(args) > 0
		target.CachePath, _ = cmd.Flags().GetString("cache")
		target.DisableCache, _ = cmd.Flags().GetBool("disable-cache")
		target.ListenOnHost, _ = cmd.Flags().GetString("listen-on-host")

		target.CertificatePath, _ = cmd.Flags().GetString("certificate")
		target.KeyPath, _ = cmd.Flags().GetString("key")
		target.DisableTls, _ = cmd.Flags().GetBool("disable-tls")

		target.UpstreamProxyUrl, _ = cmd.Flags().GetString("upstream-proxy")
		target.CORS, _ = cmd.Flags().GetBool("cors")
		target.NoImportCheck, _ = cmd.Flags().GetBool("no-import-check")

		target.Simulations, _ = cmd.Flags().GetStringSlice("import")

		if pacFileLocation, _ := cmd.Flags().GetString("pac-file"); pacFileLocation != "" {

			pacFileData, err := configuration.ReadFile(pacFileLocation)
			handleIfError(err)
			target.PACFile = string(pacFileData)
		}

		if clientAuthenticationDestination, _ := cmd.Flags().GetString("client-authentication-destination"); clientAuthenticationDestination != "" {
			_, err := regexp.Compile(clientAuthenticationDestination)

			if err != nil {
				handleIfError(errors.New("Client AuthenticationDestination regex pattern does not compile"))
			}

			target.ClientAuthenticationDestination = clientAuthenticationDestination
		}
		target.ClientAuthenticationClientCert, _ = cmd.Flags().GetString("client-authentication-client-cert")
		target.ClientAuthenticationClientKey, _ = cmd.Flags().GetString("client-authentication-client-key")
		target.ClientAuthenticationCACert, _ = cmd.Flags().GetString("client-authentication-ca-cert")

		if enableAuth, _ := cmd.Flags().GetBool("auth"); enableAuth {
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")

			if username == "" {
				username = askForInput("Username", false)
			}
			if password == "" {
				password = askForInput("Password", true)
				fmt.Println("")
			}

			target.AuthEnabled = true
			target.Username = username
			target.Password = password
		}

		target.LogOutput, _ = cmd.Flags().GetStringSlice("logs-output")
		target.LogFile, _ = cmd.Flags().GetString("logs-file")

		hasLogOutputFile := false
		for _, logOutput := range target.LogOutput {
			if logOutput == "file" {
				hasLogOutputFile = true
			}
			if logOutput != "console" && logOutput != "file" {
				handleIfError(fmt.Errorf("Unknown logs-output value: " + logOutput))
			}
		}
		if !hasLogOutputFile {
			cmd.Flags().Visit(func(f *pflag.Flag) {
				if f.Name == "logs-file" {
					handleIfError(fmt.Errorf("Flag -logs-file is not allowed unless -logs-output is set to 'file'."))
				}
			})
		}

		logLevelFlag, _ := cmd.Flags().GetString("log-level")
		logLevel, err := log.ParseLevel(logLevelFlag)
		if err != nil {
			log.WithFields(log.Fields{
				"log-level": logLevelFlag,
			}).Fatal("Unknown log-level value")
		}
		target.LogLevel = logLevel.String()

		err = wrapper.Start(target)
		handleIfError(err)

		data := [][]string{
			{"admin-port", strconv.Itoa(target.AdminPort)},
		}

		if target.Webserver {
			fmt.Println("Hoverfly is now running as a webserver")
			data = append(data, []string{"webserver-port", strconv.Itoa(target.ProxyPort)})
		} else {
			fmt.Println("Hoverfly is now running")
			data = append(data, []string{"proxy-port", strconv.Itoa(target.ProxyPort)})
		}

		drawTable(data, false)

		config.NewTarget(*target)
		handleIfError(config.WriteToFile(hoverflyDirectory))
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().String("new-target", "", "A name for a new target that hoverctl will create and associate the Hoverfly instance to")

	startCmd.Flags().Int("admin-port", 0, "A port number for the Hoverfly API/GUI. Overrides the default Hoverfly admin port (8888)")
	startCmd.Flags().Int("proxy-port", 0, "A port number for the Hoverfly proxy. Overrides the default Hoverfly proxy port (8500)")
	startCmd.Flags().String("host", "", "A host on which a Hoverfly instance is running. Overrides the default Hoverfly host (localhost)")

	startCmd.Flags().String("cache", "", "A path to a BoltDB file with persisted user and token data for authentication (DEPRECATED)")
	startCmd.Flags().Bool("disable-cache", false, "Disable the request/response cache on Hoverfly (the cache that sits in front of matching)")
	startCmd.Flags().String("certificate", "", "A path to a certificate file. Overrides the default Hoverfly certificate")
	startCmd.Flags().String("key", "", "A path to a key file. Overrides the default Hoverfly TLS key")
	startCmd.Flags().Bool("disable-tls", false, "Disable TLS verification")
	startCmd.Flags().String("upstream-proxy", "", "A host for which Hoverfly will proxy its requests to")
	startCmd.Flags().String("pac-file", "", "Configure upstream proxy by PAC file")
	startCmd.Flags().String("listen-on-host", "", "Bind hoverfly listener to a host")
	startCmd.Flags().Bool("cors", false, "Enable CORS support")
	startCmd.Flags().Bool("no-import-check", false, "Skip duplicate request check when importing simulations")

	startCmd.Flags().String("client-authentication-destination", "", "Regular expression for hosts need client authentication")
	startCmd.Flags().String("client-authentication-client-cert", "", "Path to client certificate file used for authentication")
	startCmd.Flags().String("client-authentication-client-key", "", "Path to client key file used for authentication")
	startCmd.Flags().String("client-authentication-ca-cert", "", "Path to ca cert file used for authentication")

	startCmd.Flags().Bool("auth", false, "Enable authentication on Hoverfly")
	startCmd.Flags().String("username", "", "Username to authenticate Hoverfly")
	startCmd.Flags().String("password", "", "Password to authenticate Hoverfly")

	startCmd.Flags().StringSlice("import", []string{}, "Simulations to import")

	startCmd.Flags().StringSlice("logs-output", []string{}, "Locations for log output, \"console\"(default) or \"file\"")
	startCmd.Flags().String("logs-file", "", "Log file name. Use \"hoverfly-<target name>.log\" if not provided")
	startCmd.Flags().String("log-level", "info", "Set log level (panic, fatal, error, warn, info or debug)")
}
