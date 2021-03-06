package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nyaxt/otaru/facade"
	"github.com/nyaxt/otaru/logger"
	"github.com/nyaxt/otaru/version"
)

var mylog = logger.Registry().Category("otaru-server")

var (
	flagVersion   = flag.Bool("version", false, "Show version info")
	flagReadOnly  = flag.Bool("readonly", false, "Mount as read-only mode. No changes to the filesystem is allowed.")
	flagConfigDir = flag.String("configDir", facade.DefaultConfigDir(), "Config dirpath")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Print(version.DumpBuildInfo())
		os.Exit(1)
	}

	facade.BootstrapLogger()

	cfg, err := facade.NewConfig(*flagConfigDir)
	if err != nil {
		logger.Criticalf(mylog, "%v", err)
		flag.Usage()
		os.Exit(2)
	}
	if *flagReadOnly {
		cfg.ReadOnly = true
	}
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(2)
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	signal.Notify(sigC, syscall.SIGTERM)
	go func() {
		for s := range sigC {
			logger.Warningf(mylog, "Received signal: %v", s)
			cancel()
		}
	}()
	logger.Registry().AddOutput(logger.HandleCritical(func() {
		logger.Warningf(mylog, "Starting shutdown due to critical event.")
		cancel()
	}))

	if err := facade.Serve(ctx, cfg); err != nil {
		logger.Warningf(mylog, "facade.Serve end: %v", err)
	}
}
