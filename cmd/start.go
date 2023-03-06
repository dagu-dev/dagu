package main

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/yohamta/dagu/internal/agent"
	"github.com/yohamta/dagu/internal/dag"
)

func newStartCommand() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "dagu start [--params=\"<params>\"] <DAG file>",
		Flags: append(
			globalFlags,
			&cli.StringFlag{
				Name:     "params",
				Usage:    "parameters",
				Value:    "",
				Required: false,
			},
		),
		Action: func(c *cli.Context) error {
			d, err := loadDAG(c, c.Args().Get(0), strings.Trim(c.String("params"), "\""))
			if err != nil {
				return err
			}
			return start(d)
		},
	}
}

func start(d *dag.DAG) error {
	a := &agent.Agent{AgentConfig: &agent.AgentConfig{
		DAG: d,
		Dry: false,
	}}

	listenSignals(func(sig os.Signal) {
		a.Signal(sig)
	})

	return a.Run()
}
