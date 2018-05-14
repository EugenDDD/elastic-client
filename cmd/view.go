package cmd

import (
	"elastic-search/global"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var configSerial []byte
var err error
var all, server, port, user, pass, viewConf bool
var els string

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View elastic config",
	Long:  `View elastic config details`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {

		noFlag := true

		configSerial, err = ioutil.ReadFile(global.ConfigFile)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(configSerial, &config)
		if err != nil {
			panic(err)
		}

		if cmd.Flags().Lookup("all").Changed {
			fmt.Println(string(configSerial))
			return
		}

		if cmd.Flags().Lookup("server").Changed {
			fmt.Println("server:", config.Server)
			noFlag = false
		}

		if cmd.Flags().Lookup("port").Changed {
			fmt.Println("port:", config.Port)
			noFlag = false
		}

		if cmd.Flags().Lookup("user").Changed {
			fmt.Println("user:", config.User)
			noFlag = false
		}

		if cmd.Flags().Lookup("pass").Changed {
			fmt.Println("pass:", config.Pass)
			noFlag = false
		}

		if noFlag {
			fmt.Printf("Please provide at least one flag!\n\n")
			cmd.Usage()
		}
		return

	},
}

func init() {
	viewCmd.Flags().BoolVarP(&all, "all", "a", false, "View all configs")
	viewCmd.Flags().BoolVarP(&server, "server", "s", false, "Set kibana server")
	viewCmd.Flags().BoolVarP(&port, "port", "P", false, "Set kibana server port")
	viewCmd.Flags().BoolVarP(&user, "user", "u", false, "Set kibana server username")
	viewCmd.Flags().BoolVarP(&pass, "pass", "p", false, "Set kibana username password")
}
