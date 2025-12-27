package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/markwharris77/atlas/internal/commands"
	"github.com/markwharris77/atlas/internal/config"
	"github.com/markwharris77/atlas/internal/log"
	"github.com/markwharris77/atlas/internal/tools"
)

func main() {

	configPath := flag.String("file", "atlas.yml", "path to config file")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			`atlas - simple deploy tool
Usage:
atlas [flags] <command>

Commands:
describe  Print loaded config
deploy    Deploy (stub for now)

Flags:`)
		flag.PrintDefaults()
	}

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose output")

	flag.Parse()

	log.SetVerbose(verbose)

	cfg, err := config.Load(*configPath)

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// Command (first non-flag arg)
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("misisng command")
		flag.Usage()
		os.Exit(3)
	}

	cmd := args[0]

	switch cmd {
	case "describe":
		if err := config.Print(cfg); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(2)
		}
	case "deploy":
		err = commands.RunDeploy(commands.DeployOptions{
			Name:      cfg.App.Name,
			LocalPath: cfg.Deploy.LocalDir,
			Target: tools.Target{
				User: cfg.Deploy.User,
				Host: cfg.Deploy.Host,
				Port: cfg.Deploy.Port,
			},
			RunUser:    cfg.Run.User,
			RunCommand: cfg.Run.Command,
			RunEnv:     cfg.Run.Env,
		})

		if err != nil {
			log.Error("failed to deploy: %s", err)
		}

	default:
		fmt.Fprintln(os.Stderr, "unknown command:", cmd)
		flag.Usage()
		os.Exit(4)
	}

}
