package main

import (
	"fmt"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"

	svc "github.com/vshiva/goreactapp"
)

var (
	// errorExitCode returns a urfave decorated error which indicates a exit
	// code 1. To be returned from a urfave action.
	errorExitCode = cli.NewExitError("", 1)
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s\n Version:  %s\n Git Commit:  %s\n Go Version:  %s\n OS/Arch:  %s/%s\n Built:  %s\n",
			c.App.Name, c.App.Version, svc.GitCommit,
			runtime.Version(), runtime.GOOS, runtime.GOARCH, c.App.Compiled.String())
	}

	app := cli.NewApp()

	app.Name = "cool crazy app"
	app.Copyright = "(c) 2018 Copyright"
	app.Usage = "cool crazy app description"

	app.Version = svc.Version()
	app.Compiled = svc.CompiledAt()
	app.Before = setupLogging
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
		},
	}
	app.Commands = []cli.Command{
		serverCommand,
	}

	app.Run(os.Args)
}

//SetupLogging helps with setting up logger
func setupLogging(c *cli.Context) error {
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	// Dynamically return false or true based on the logger output's
	// file descriptor referring to a terminal or not.
	if os.Getenv("TERM") == "dumb" || !isLogrusTerminal() {
		log.SetFormatter(log.Formatter(&log.JSONFormatter{}))
	}
	return nil
}

// isLogrusTerminal checks if the standard logger of Logrus is a terminal.
func isLogrusTerminal() bool {
	w := log.StandardLogger().Out
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
