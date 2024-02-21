// Copyright 2020 by Red Hat, Inc. All rights reserved.
// Use of this source is goverend by the Apache License
// that can be found in the LICENSE file.

/*
Package root is the top of the subcommand parser

root handles setup of the command line flags and initialization
of the weldr API client's configuration values. It also holds
commandline flags that can be accessed by subcommands.
*/
package root

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/osbuild/weldr-client/v2/cloud"
	"github.com/osbuild/weldr-client/v2/weldr"
)

const helpTemplate = `{{.Short}}{{if .Long}}

Description:
  {{.Long}}{{end}}

{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

// jsonResponse contains the responses from the API server
type jsonResponse struct {
	Method string                 `json:"method"`
	Path   string                 `json:"path"`
	Status int                    `json:"status"`
	Body   map[string]interface{} `json:"body"`
}

var (
	rootCmd = &cobra.Command{
		Use:   path.Base(os.Args[0]),
		Short: "composer commandline tool",
		Long:  "commandline tool for osbuild-composer",
	}
	docCmd = &cobra.Command{
		Use:   "doc DIRECTORY",
		Short: "Generate manpage files",
		Long:  "Generate manpage files for all the commands, one per file",
		RunE: func(cmd *cobra.Command, args []string) error {
			header := &doc.GenManHeader{
				Title:   "COMPOSER-CLI",
				Section: "1",
				Source:  "composer-cli",
			}
			return doc.GenManTree(rootCmd, header, args[0])
		},
		Args: cobra.ExactArgs(1),
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display the version and exit",
		Long:  "Display the version and git hash used for the build",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("composer-cli v%s\n", Version)
			return nil
		},
	}
	apiVersion  int
	httpTimeout int
	// JSONOutput is the state of --json cmdline flag
	JSONOutput      bool
	logPath         string
	weldrSocketPath string
	cloudSocketPath string
	testMode        int

	// Version is set by the build
	Version = "DEVEL"

	// Client is the weldr.Client used to communicate with the server
	Client weldr.Client
	// Cloud is the cloud.Client used to communicate with the cloudapi service
	Cloud cloud.Client

	// Original Stdout
	oldStdout *os.File

	// jsonResponses holds the responses from the API server
	jsonResponses []jsonResponse
)

func init() {
	rootCmd.PersistentFlags().IntVarP(&apiVersion, "api", "a", 1, "WELDR Server API Version to use")
	rootCmd.PersistentFlags().BoolVarP(&JSONOutput, "json", "j", false, "Output the raw JSON response instead of the normal output")
	rootCmd.PersistentFlags().StringVar(&logPath, "log", "", "Path to optional logfile")
	rootCmd.PersistentFlags().StringVarP(&weldrSocketPath, "socket", "s", "/run/weldr/api.socket", "Path to the WELDR API server's socket file")
	rootCmd.PersistentFlags().StringVarP(&cloudSocketPath, "cloudsocket", "", "/run/cloudapi/api.socket", "Path to the cloudapi server's socket file")
	rootCmd.PersistentFlags().IntVar(&testMode, "test", 0, "Pass test mode to compose. 1=Mock compose with fail. 2=Mock compose with finished.")
	rootCmd.PersistentFlags().IntVar(&httpTimeout, "timeout", 240, "Timeout to use for server communication. Set to 0 for no timeout")

}

// Init sets up Cobra and adds the doc command to the root cmdline parser
func Init() {
	cobra.OnInitialize(initConfig)

	// Command to generate manpage documentation
	AddRootCommand(docCmd)

	// Display the version
	AddRootCommand(versionCmd)
}

func initConfig() {
	initWeldrClient()
	initCloudClient()
	setupJSONOutput()
}

func initWeldrClient() {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx = context.Background()

	if httpTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(httpTimeout)*time.Second)
		defer cancel()
	}

	Client = weldr.InitClientUnixSocket(ctx, apiVersion, weldrSocketPath)
}

func initCloudClient() {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx = context.Background()

	if httpTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(httpTimeout)*time.Second)
		defer cancel()
	}

	Cloud = cloud.InitClientUnixSocket(ctx, cloudSocketPath)
}

// setupJSONOutput configures the callback function and disables Stdout
func setupJSONOutput() {
	if JSONOutput {
		// Disable Stdout output so that only json is output
		oldStdout = os.Stdout
		os.Stdout = nil

		Client.SetRawCallback(func(method string, path string, status int, data []byte) {
			// Convert the data to a generic data structure, then pretty-print it
			var r jsonResponse
			r.Method = method
			r.Path = path
			r.Status = status
			err := json.Unmarshal(data, &r.Body)
			if err == nil {
				jsonResponses = append(jsonResponses, r)
			}
		})
		Cloud.SetRawCallback(func(method string, path string, status int, data []byte) {
			// Convert the data to a generic data structure, then pretty-print it
			var r jsonResponse
			r.Method = method
			r.Path = path
			r.Status = status
			err := json.Unmarshal(data, &r.Body)
			if err == nil {
				jsonResponses = append(jsonResponses, r)
			}
		})
	} else {
		if oldStdout != nil {
			os.Stdout = oldStdout
			oldStdout = nil
		}
		Client.SetRawCallback(func(string, string, int, []byte) {})
		Cloud.SetRawCallback(func(string, string, int, []byte) {})
	}
}

// Execute runs the commands on the commandline
func Execute() error {
	err := rootCmd.Execute()
	if JSONOutput {
		s, jerr := json.MarshalIndent(jsonResponses, "", "    ")
		if jerr == nil {
			fmt.Fprintln(oldStdout, string(s))
		}
	}
	return err
}

// AddRootCommand adds a cobra command to the list of root commands
func AddRootCommand(cmd *cobra.Command) {
	cmd.SetHelpTemplate(helpTemplate)
	rootCmd.AddCommand(cmd)
}

// ExecutionError prints an error to stderr, sets silent flags on the cobra command, and
// returns an error to the caller suitable for assignment to error
func ExecutionError(cmd *cobra.Command, format string, a ...interface{}) error {
	s := fmt.Sprintf(format, a...)
	if len(s) > 0 {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", s)
	}
	cmd.SilenceErrors = true // cobra will not print errors returned from commands after this
	cmd.SilenceUsage = true  // cobra will not print usage on errors after this
	return fmt.Errorf(s)
}

// ExecutionErrors prints a list of errors to stderr, then calls ExecutionError
func ExecutionErrors(cmd *cobra.Command, errors []weldr.APIErrorMsg) error {
	// When JSON output is enabled the errors are in the JSON so skip printing them
	if !JSONOutput {
		for _, s := range errors {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", s)
		}
	}
	return ExecutionError(cmd, "")
}

// GetCommaArgs returns a list of the arguments, split by commas and spaces
// They can be grouped or separated, the return list should be the same for all variations
// empty fields, eg. ,, are ignored by collapsing repeated , and spaces into one.
func GetCommaArgs(args []string) []string {
	var result []string
	// Gather up all the arguments, with or without commas
	f := func(c rune) bool {
		return c == ',' || c == ' '
	}
	for _, arg := range args {
		result = append(result, strings.FieldsFunc(arg, f)...)
	}
	return result
}
