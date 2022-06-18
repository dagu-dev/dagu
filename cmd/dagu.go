package cmd

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	"github.com/yohamta/dagu/internal/constants"
)

var version = "0.0.0"
var stdin io.ReadCloser

func main() {
	setVersion()
	err := run()
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func listenSignals(abortFunc func(sig os.Signal)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigs {
			log.Printf("\nSignal: %v", sig)
			abortFunc(sig)
		}
	}()
}

func setVersion() {
	constants.Version = version
}

func run() error {
	stdin = os.Stdin
	app := makeApp()
	return app.Run(os.Args)
}

func makeApp() *cli.App {
	return &cli.App{
		Name:      "Dagu",
		Usage:     "Self-contained, easy-to-use workflow engine for smaller use cases",
		UsageText: "dagu [options] <start|status|stop|retry|dry|server|version> [args]",
		Commands: []*cli.Command{
			newStartCommand(),
			newStatusCommand(),
			newStopCommand(),
			newRetryCommand(),
			newDryCommand(),
			newServerCommand(),
			newVersionCommand(),
		},
	}
}
