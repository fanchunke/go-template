package main

import (
	"fmt"
	"go-template/internal/config"
	"go-template/internal/log"
	"go-template/internal/server"
	"go-template/internal/signals"
	"go-template/internal/version"
	"os"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	// flags definition
	fs := pflag.NewFlagSet("go-template", pflag.ContinueOnError)
	configFile := fs.String("conf", "/home/works/program/conf/online.conf", "configiration file")

	versionFlag := fs.BoolP("version", "v", false, "get version number")

	// parse flags
	err := fs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.PrintDefaults()
		os.Exit(2)
	case *versionFlag:
		fmt.Println(version.VERSION)
		os.Exit(0)
	}

	// load config
	cfg := config.New(*configFile)

	// configure logging
	logger, _ := log.New(cfg.Logger.Level, cfg.Logger.OutputPaths, cfg.Logger.ErrorOutputPaths)
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// log version and port
	logger.Info("Starting server",
		zap.Any("config", cfg),
	)

	// start HTTP server
	srv, _ := server.NewServer(cfg)
	stopCh := signals.SetupSignalHandler()
	srv.Run(stopCh)
}
