// MIT License
// Copyright (c) 2020 ysicing <i@ysicing.me>

package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/cmd/flags"
	"github.com/ysicing/ergo/common"
	"github.com/ysicing/ergo/pkg/ergo/plugin"
	"github.com/ysicing/ergo/pkg/util/factory"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

const (
	cliName        = "ergo"
	cliDescription = "A simple command line client for devops"
)

var (
	globalFlags *flags.GlobalFlags
)

func Execute() {
	// create a new factory
	f := factory.DefaultFactory()
	// build the root command
	rootCmd := BuildRoot(f)
	// before hook
	// execute command
	err := rootCmd.Execute()
	// after hook
	if err != nil {
		if globalFlags.Debug {
			f.GetLog().Fatalf("%+v", err)
		} else {
			f.GetLog().Fatal(err)
		}
	}
}

// BuildRoot creates a new root command from the
func BuildRoot(f factory.Factory) *cobra.Command {
	// build the root cmd
	rootCmd := NewRootCmd(f)
	persistentFlags := rootCmd.PersistentFlags()
	globalFlags = flags.SetGlobalFlags(persistentFlags)
	// Add sub commands

	// Add main commands
	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newUpgradeCmd())
	rootCmd.AddCommand(newDebianCmd(f))
	rootCmd.AddCommand(newOPSCmd(f))
	rootCmd.AddCommand(newRepoCmd(f))
	rootCmd.AddCommand(newPluginCmd(f))
	rootCmd.AddCommand(newCodeGenCmd(f))
	rootCmd.AddCommand(newCloudCommand(f))
	// Add plugin commands

	args := os.Args
	if len(args) > 1 {

		pluginHandler := NewDefaultPluginHandler(plugin.ValidPluginFilenamePrefixes)

		cmdPathPieces := args[1:]
		if _, _, err := rootCmd.Find(cmdPathPieces); err != nil {
			var cmdName string // first "non-flag" arguments
			for _, arg := range cmdPathPieces {
				if !strings.HasPrefix(arg, "-") {
					cmdName = arg
					break
				}
			}

			switch cmdName {
			case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
				// Don't search for a plugin
			default:
				if err := HandlePluginCommand(pluginHandler, cmdPathPieces); err != nil {
					fmt.Fprintf(os.Stdout, "Error: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	return rootCmd
}

type PluginHandler interface {
	Lookup(filename string) (string, bool)
	Execute(executablePath string, cmdArgs, environment []string) error
}

func NewDefaultPluginHandler(validPrefixes []string) *DefaultPluginHandler {
	return &DefaultPluginHandler{
		ValidPrefixes: validPrefixes,
	}
}

type DefaultPluginHandler struct {
	ValidPrefixes []string
}

// Lookup implements PluginHandler
func (h *DefaultPluginHandler) Lookup(filename string) (string, bool) {
	p, _ := os.LookupEnv("PATH")
	ergobin := common.GetDefaultBinDir()
	if !strings.Contains(p, ergobin) {
		os.Setenv("PATH", fmt.Sprintf("%v:%v", p, ergobin))
	}
	for _, prefix := range h.ValidPrefixes {
		path, err := exec.LookPath(fmt.Sprintf("%s-%s", prefix, filename))
		if err != nil || len(path) == 0 {
			continue
		}
		return path, true
	}

	return "", false
}

// Execute implements PluginHandler
func (h *DefaultPluginHandler) Execute(executablePath string, cmdArgs, environment []string) error {

	// Windows does not support exec syscall.
	if runtime.GOOS == "windows" {
		cmd := exec.Command(executablePath, cmdArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Env = environment
		err := cmd.Run()
		if err == nil {
			os.Exit(0)
		}
		return err
	}

	// invoke cmd binary relaying the environment and args given
	// append executablePath to cmdArgs, as execve will make first argument the "binary name".
	return syscall.Exec(executablePath, append([]string{executablePath}, cmdArgs...), environment)
}

func HandlePluginCommand(pluginHandler PluginHandler, cmdArgs []string) error {
	var remainingArgs []string // all "non-flag" arguments
	for _, arg := range cmdArgs {
		if strings.HasPrefix(arg, "-") {
			break
		}
		remainingArgs = append(remainingArgs, strings.Replace(arg, "-", "_", -1))
	}

	if len(remainingArgs) == 0 {
		// the length of cmdArgs is at least 1
		return fmt.Errorf("flags cannot be placed before plugin name: %s", cmdArgs[0])
	}

	foundBinaryPath := ""

	// attempt to find binary, starting at longest possible name with given cmdArgs
	for len(remainingArgs) > 0 {
		path, found := pluginHandler.Lookup(strings.Join(remainingArgs, "-"))
		if !found {
			remainingArgs = remainingArgs[:len(remainingArgs)-1]
			continue
		}

		foundBinaryPath = path
		break
	}

	if len(foundBinaryPath) == 0 {
		return nil
	}

	// invoke cmd binary relaying the current environment and args given
	if err := pluginHandler.Execute(foundBinaryPath, cmdArgs[len(remainingArgs):], os.Environ()); err != nil {
		return err
	}

	return nil
}

// NewRootCmd returns a new root command
func NewRootCmd(f factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:           cliName,
		SilenceUsage:  true,
		SilenceErrors: true,
		Short:         "ergo, ergo, NB!",
		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			if cobraCmd.Annotations != nil {
				return nil
			}
			log := f.GetLog()
			if globalFlags.Silent {
				log.SetLevel(logrus.FatalLevel)
			} else if globalFlags.Debug {
				log.SetLevel(logrus.DebugLevel)
			}

			// apply extra flags TODO
			return nil
		},
		Long: cliDescription,
	}
}
