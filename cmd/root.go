package cmd

import (
	"elastic-search/global"
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "elastic-search",
	// Short: "An example of elastic",
	// 	Long: `This application shows how to create modern CLI
	// applications in go using Cobra CLI library`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(getLogCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if global.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(global.ConfigFile)
	} else {
		// Find home directory.
		home, err := homedir.Expand(global.ConfigFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra-example" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
		// fmt.Println()
	}
}
