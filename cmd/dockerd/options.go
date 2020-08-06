package main

import (
	"os"
	"path/filepath"

	cliconfig "github.com/docker/docker/cli/config"
	"github.com/docker/docker/daemon/config"
	"github.com/docker/docker/opts"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/spf13/pflag"
)

const (
	// DefaultCaFile is the default filename for the CA pem file
	DefaultCaFile = "ca.pem"
	// DefaultKeyFile is the default filename for the key pem file
	DefaultKeyFile = "key.pem"
	// DefaultCertFile is the default filename for the cert pem file
	DefaultCertFile = "cert.pem"
	// FlagTLSVerify is the flag name for the TLS verification option
	FlagTLSVerify = "tlsverify"
)

var (
	dockerCertPath  = os.Getenv("DOCKER_CERT_PATH")
	dockerTLSVerify = os.Getenv("DOCKER_TLS_VERIFY") != ""
)

type daemonOptions struct {
	configFile   string
	daemonConfig *config.Config
	flags        *pflag.FlagSet
	Debug        bool
	Hosts        []string
	LogLevel     string
	TLS          bool
	TLSVerify    bool
	TLSOptions   *tlsconfig.Options
}

// newDaemonOptions returns a new daemonFlags
func newDaemonOptions(config *config.Config) *daemonOptions {
	return &daemonOptions{
		daemonConfig: config,
	}
}

// InstallFlags adds flags for the common options on the FlagSet
func (o *daemonOptions) InstallFlags(flags *pflag.FlagSet) {
	if dockerCertPath == "" {
		// cliconfig.Dir returns $DOCKER_CONFIG or ~/.docker.
		// cliconfig.Dir does not look up $XDG_CONFIG_HOME
		dockerCertPath = cliconfig.Dir()
	}

	flags.BoolVarP(&o.Debug, "debug", "D", false, "Enable debug mode")
	flags.StringVarP(&o.LogLevel, "log-level", "l", "info", `Set the logging level ("debug"|"info"|"warn"|"error"|"fatal")`)
	flags.BoolVar(&o.TLS, "tls", false, "Use TLS; implied by --tlsverify")
	flags.BoolVar(&o.TLSVerify, FlagTLSVerify, dockerTLSVerify, "Use TLS and verify the remote")

	// TODO use flag flags.String("identity"}, "i", "", "Path to libtrust key file")

	o.TLSOptions = &tlsconfig.Options{
		CAFile:   filepath.Join(dockerCertPath, DefaultCaFile),
		CertFile: filepath.Join(dockerCertPath, DefaultCertFile),
		KeyFile:  filepath.Join(dockerCertPath, DefaultKeyFile),
	}
	tlsOptions := o.TLSOptions
	flags.Var(opts.NewQuotedString(&tlsOptions.CAFile), "tlscacert", "Trust certs signed only by this CA")
	flags.Var(opts.NewQuotedString(&tlsOptions.CertFile), "tlscert", "Path to TLS certificate file")
	flags.Var(opts.NewQuotedString(&tlsOptions.KeyFile), "tlskey", "Path to TLS key file")

	hostOpt := opts.NewNamedListOptsRef("hosts", &o.Hosts, opts.ValidateHost)
	flags.VarP(hostOpt, "host", "H", "Daemon socket(s) to connect to")
}

// SetDefaultOptions sets default values for options after flag parsing is
// complete
// SetDefaultOptions在标志解析完成后设置选项的默认值
func (o *daemonOptions) SetDefaultOptions(flags *pflag.FlagSet) {
	// Regardless of whether the user sets it to true or false, if they
	// specify --tlsverify at all then we need to turn on TLS
	// TLSVerify can be true even if not set due to DOCKER_TLS_VERIFY env var, so we need
	// to check that here as well
	//
	// 无论用户将其设置为TRUE还是FALSE，如果他们指定了--tlsVerify，那么我们需要打开TLS
	// TLS即使由于DOKER_TLS_VERIFY环境变量而未设置Verify Can为TRUE，因此我们也需要在此处进行检查
	if flags.Changed(FlagTLSVerify) || o.TLSVerify {
		o.TLS = true
	}

	if !o.TLS {
		o.TLSOptions = nil
	} else {
		tlsOptions := o.TLSOptions
		tlsOptions.InsecureSkipVerify = !o.TLSVerify

		// Reset CertFile and KeyFile to empty string if the user did not specify
		// the respective flags and the respective default files were not found.
		if !flags.Changed("tlscert") {
			if _, err := os.Stat(tlsOptions.CertFile); os.IsNotExist(err) {
				tlsOptions.CertFile = ""
			}
		}
		if !flags.Changed("tlskey") {
			if _, err := os.Stat(tlsOptions.KeyFile); os.IsNotExist(err) {
				tlsOptions.KeyFile = ""
			}
		}
	}
}
