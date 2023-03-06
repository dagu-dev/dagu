package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yohamta/dagu/internal/agent"
	"github.com/yohamta/dagu/internal/dag"
)

func newDryCommand() *cli.Command {
	return &cli.Command{
		Name:  "dry",
		Usage: "dagu dry [--params=\"<params>\"] <config>",
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
			d, err := loadDAG(c, c.Args().Get(0), c.String("params"))
			if err != nil {
				return err
			}
			return dryRun(d)
		},
	}
}

func dryRun(d *dag.DAG) error {
	a := &agent.Agent{AgentConfig: &agent.AgentConfig{
		DAG: d,
		Dry: true,
	}}
	listenSignals(func(sig os.Signal) {
		a.Signal(sig)
	})
	return a.Run()
}
