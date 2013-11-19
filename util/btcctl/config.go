package main

import (
	"fmt"
	"github.com/conformal/btcutil"
	"github.com/conformal/go-flags"
	"os"
	"path/filepath"
)

var (
	btcdHomeDir        = btcutil.AppDataDir("btcd", false)
	btcctlHomeDir      = btcutil.AppDataDir("btcctl", false)
	defaultConfigFile  = filepath.Join(btcctlHomeDir, "btcctl.conf")
	defaultRPCCertFile = filepath.Join(btcdHomeDir, "rpc.cert")
)

// config defines the configuration options for btcctl.
//
// See loadConfig for details on the configuration load process.
type config struct {
	ConfigFile    string `short:"C" long:"configfile" description:"Path to configuration file"`
	RpcUser       string `short:"u" long:"rpcuser" description:"RPC username"`
	RpcPassword   string `short:"P" long:"rpcpass" description:"RPC password"`
	RpcServer     string `short:"s" long:"rpcserver" description:"RPC server to connect to"`
	RpcCert       string `short:"c" long:"rpccert" description:"RPC server certificate chain for validation"`
	TlsSkipVerify bool   `long:"skipverify" description:"Do not verify tls certificates (not recommended!)"`
}

// loadConfig initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
// 	1) Start with a default config with sane settings
// 	2) Pre-parse the command line to check for an alternative config file
// 	3) Load configuration file overwriting defaults with any specified options
// 	4) Parse CLI options and overwrite/add any specified options
//
// The above results in functioning properly without any config settings
// while still allowing the user to override settings with config files and
// command line options.  Command line options always take precedence.
func loadConfig() (*flags.Parser, *config, []string, error) {
	// Default config.
	cfg := config{
		ConfigFile: defaultConfigFile,
		RpcServer:  "localhost:8334",
		RpcCert:    defaultRPCCertFile,
	}

	// Pre-parse the command line options to see if an alternative config
	// file was specified.  Any errors can be ignored here since they will
	// be caught be the final parse below.
	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.None)
	preParser.Parse()

	// Load additional config from file.
	parser := flags.NewParser(&cfg, flags.PassDoubleDash|flags.HelpFlag)
	err := parser.ParseIniFile(preCfg.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			fmt.Fprintln(os.Stderr, err)
			return parser, nil, nil, err
		}
	}

	// Parse command line options again to ensure they take precedence.
	remainingArgs, err := parser.Parse()
	if err != nil {
		return parser, nil, nil, err
	}

	return parser, &cfg, remainingArgs, nil
}
