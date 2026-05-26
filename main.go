package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/93percentai/sitemaker/internal/build"
	"github.com/93percentai/sitemaker/internal/config"
	"github.com/93percentai/sitemaker/internal/server"
)

//go:embed defaults
var defaultsFS embed.FS

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		cmdBuild(os.Args[2:])
	case "serve":
		cmdServe(os.Args[2:])
	case "init":
		cmdInit(os.Args[2:])
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func cmdBuild(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	configPath := fs.String("config", "sitemaker.toml", "path to config file")
	fs.Parse(args)

	cfg := loadConfig(*configPath)
	defaults, _ := getDefaultsFS()

	if err := build.Build(cfg, defaults); err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Build complete.")
}

func cmdServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	configPath := fs.String("config", "sitemaker.toml", "path to config file")
	port := fs.Int("port", 8080, "port to serve on")
	fs.Parse(args)

	cfg := loadConfig(*configPath)
	defaults, _ := getDefaultsFS()

	if err := server.Serve(cfg, defaults, *port); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func cmdInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	tmpl := fs.String("template", "personal", "template type: personal or org")
	fs.Parse(args)

	dir := "."
	if fs.NArg() > 0 {
		dir = fs.Arg(0)
	}

	defaults, _ := getDefaultsFS()
	if err := scaffold(dir, *tmpl, defaults); err != nil {
		fmt.Fprintf(os.Stderr, "init failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Initialized sitemaker project in %s (template: %s)\n", dir, *tmpl)
}

func loadConfig(path string) config.Config {
	cfg, err := config.Load(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config.Default()
		}
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	return cfg
}

func getDefaultsFS() (fs.FS, error) {
	return fs.Sub(defaultsFS, "defaults")
}

func printUsage() {
	fmt.Println(`sitemaker - static site generator for personal and org websites

Usage:
  sitemaker <command> [options]

Commands:
  build     Build the site to the output directory
  serve     Start a dev server with live reload
  init      Initialize a new sitemaker project
  help      Show this help message

Options for build/serve:
  -config   Path to config file (default: sitemaker.toml)

Options for serve:
  -port     Port to serve on (default: 8080)

Options for init:
  -template Template type: personal or org (default: personal)`)
}
