package main

import (
	"errors"
	"log"

	"github.com/yohamta/dagman/internal/config"
	"github.com/yohamta/dagman/internal/controller"
	"github.com/yohamta/dagman/internal/models"

	"github.com/urfave/cli/v2"
)

func newStatusCommand() *cli.Command {
	cl := config.NewConfigLoader()
	return &cli.Command{
		Name:  "status",
		Usage: "dagman status <config>",
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("config file must be specified.")
			}
			config_file_path := c.Args().Get(0)
			cfg, err := cl.Load(config_file_path, "")
			if err != nil {
				return err
			}
			return queryStatus(cfg)
		},
	}
}

func queryStatus(cfg *config.Config) error {
	status, err := controller.New(cfg).GetStatus()
	if err != nil {
		return err
	}
	res := &models.StatusResponse{
		Status: status,
	}
	log.Printf("Pid=%d Status=%s", res.Status.Pid, res.Status.Status)
	return nil
}
