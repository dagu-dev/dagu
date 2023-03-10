/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yohamta/dagu/internal/config"
	"github.com/yohamta/dagu/internal/constants"
	"github.com/yohamta/dagu/internal/dag"
)

var (
	cfgFile string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "dagu",
		Short: "YAML-based DAG scheduling tool.",
		Long:  `YAML-based DAG scheduling tool.`,
	}

	version = "0.0.0"
	sigs    chan os.Signal
)

const legacyPath = ".dagu"

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	constants.Version = version
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dagu/admin.yaml)")

	regisgterCommands(rootCmd)
}

func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(path.Join(home, legacyPath))
		viper.SetConfigType("yaml")
		viper.SetConfigName("admin")
	}

	cobra.CheckErr(config.LoadConfig(home))
}

func loadDAG(dagFile, params string) (d *dag.DAG, err error) {
	dagLoader := &dag.Loader{BaseConfig: config.Get().BaseConfig}
	return dagLoader.Load(dagFile, params)
}

func listenSignals(abortFunc func(sig os.Signal)) {
	sigs = make(chan os.Signal, 100)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigs {
			abortFunc(sig)
		}
	}()
}

func getFlagString(cmd *cobra.Command, name, fallback string) string {
	if s, _ := cmd.Flags().GetString(name); s != "" {
		return s
	}
	return fallback
}
