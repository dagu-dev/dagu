package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yohamta/dagman/internal/agent"
	"github.com/yohamta/dagman/internal/config"
)

func newDryCommand() *cli.Command {
	cl := config.NewConfigLoader()
	return &cli.Command{
		Name:  "dry",
		Usage: "dagman dry [--params=\"<params>\"] <config>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "params",
				Usage:    "parameters",
				Value:    "",
				Required: false,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("config file must be specified.")
			}
			if c.NArg() != 1 {
				return errors.New("too many parameters.")
			}
			config_file_path := c.Args().Get(0)
			cfg, err := cl.Load(config_file_path, c.String("params"))
			if err != nil {
				return err
			}
			return dryRun(cfg)
		},
	}
}

func dryRun(cfg *config.Config) error {
	a := &agent.Agent{Config: &agent.Config{
		DAG: cfg,
		Dry: true,
	}}
	listenSignals(func(sig os.Signal) {
		a.Signal(sig)
	})

	err := a.Run()
	if err != nil {
		log.Printf("[DRY] failed %v", err)
	}
	return nil
}
