package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/tanksuzuki/jig/config"
)

const version = "v0.0.1"

type flag struct {
	Config  string `short:"c" long:"config" description:"Config file to load" env:"JIG_CONFIG" value-name:"path" default:"~/jig.toml"`
	Help    bool   `short:"h" long:"help" description:"Help for jig or Print your script usage"`
	Version bool   `long:"version" description:"Print version information and quit"`
}

func main() {
	var f flag

	parser := flags.NewParser(&f, flags.PassAfterNonOption)
	parser.Usage = "[<option>] <script_name> [<script_arg>...]"
	args, err := parser.Parse()
	if err != nil {
		log.Fatalf("failed to parse arguments: %s", err)
	}

	if f.Version {
		fmt.Printf("jig version %s\n", version)
		return
	}

	if f.Help && len(args) == 0 {
		parser.WriteHelp(os.Stdout)
		return
	}

	c, err := config.Read(f.Config, "")
	if err != nil {
		log.Fatalf("failed to load config file: %s", err)
	}

	switch {
	case len(args) == 0:
		c.List()
		return
	case f.Help:
		if err := c.Help(args[0]); err != nil {
			log.Fatalf("failed to print usage: %s", err)
		}
		return
	default:
		code, err := c.Run(args[0], args[1:])
		if err != nil {
			log.Fatalf("failed to run script: %s", err)
		}
		os.Exit(code)
	}
}
