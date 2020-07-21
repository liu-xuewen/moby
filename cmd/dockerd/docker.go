package main

import (
	"fmt"
	"os"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/dockerversion"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/docker/rootless"
	"github.com/moby/buildkit/util/apicaps"
	"github.com/moby/term"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	honorXDG bool
)

func newDaemonCommand() (*cobra.Command, error) {
	opts := newDaemonOptions(config.New())

	cmd := &cobra.Command{
		Use:           "dockerd [OPTIONS]",
		Short:         "A self-sufficient runtime for containers.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.flags = cmd.Flags()
			return runDaemon(opts)
		},
		DisableFlagsInUseLine: true,
		Version:               fmt.Sprintf("%s, build %s", dockerversion.Version, dockerversion.GitCommit),
	}
	cli.SetupRootCommand(cmd)

	flags := cmd.Flags()
	flags.BoolP("version", "v", false, "Print version information and quit")
	defaultDaemonConfigFile, err := getDefaultDaemonConfigFile()
	if err != nil {
		return nil, err
	}
	flags.StringVar(&opts.configFile, "config-file", defaultDaemonConfigFile, "Daemon configuration file")
	opts.InstallFlags(flags)
	if err := installConfigFlags(opts.daemonConfig, flags); err != nil {
		return nil, err
	}
	installServiceFlags(flags)

	return cmd, nil
}

func init() {
	if dockerversion.ProductName != "" {
		apicaps.ExportedProduct = dockerversion.ProductName
	}
	// When running with RootlessKit, $XDG_RUNTIME_DIR, $XDG_DATA_HOME, and $XDG_CONFIG_HOME needs to be
	// honored as the default dirs, because we are unlikely to have permissions to access the system-wide
	// directories.
	//
	// Note that even running with --rootless, when not running with RootlessKit, honorXDG needs to be kept false,
	// because the system-wide directories in the current mount namespace are expected to be accessible.
	// ("rootful" dockerd in rootless dockerd, #38702)
	/*
	使用Rootless Kit运行时，需要将$XDG_RUNTIME_DIR、$XDG_DATA_HOME和$XDG_CONFIG_HOME作为默认目录，因为我们不太可能拥有访问系统范围目录的权限。
	请注意，即使使用-rootless运行，当不使用Rootless Kit运行时，honorXDG也需要保持为false，因为当前挂载名称空间中的系统范围目录应该是可访问的。
	(“Rootful”dockerd in unroot dockerd，#38702)
	*/
	honorXDG = rootless.RunningWithRootlessKit()
}

func main() {
	if reexec.Init() {
		return
	}

	// initial log formatting; this setting is updated after the daemon configuration is loaded.
	/*
	初始日志格式化；
	此设置在加载守护程序配置后更新。
	*/
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: jsonmessage.RFC3339NanoFixed,
		FullTimestamp:   true,
	})

	// Set terminal emulation based on platform as required.
	/*
	根据需要设置基于平台的终端仿真。
	*/
	_, stdout, stderr := term.StdStreams()

	initLogging(stdout, stderr)

	onError := func(err error) {
		fmt.Fprintf(stderr, "%s\n", err)
		os.Exit(1)
	}

	cmd, err := newDaemonCommand()
	if err != nil {
		onError(err)
	}
	cmd.SetOut(stdout)
	if err := cmd.Execute(); err != nil {
		onError(err)
	}
}
