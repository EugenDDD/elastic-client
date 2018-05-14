package cmd

import (
	"elastic-search/global"
	"elastic-search/util"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var configFile string

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set elastic config",
	Long:  `Set elastic config details`,
	// Args: func(cmd *cobra.Command, args []string) error {

	// 	flagSet := cmd.Flags()
	// 	var i int
	// 	for _, flag := range flagSet {
	// 		i++
	// 	}

	// 	if len(cmd.Flags) < 1 {
	// 		return errors.New("requires at least one arg")
	// 	}
	// 	return nil
	// },
	Run: func(cmd *cobra.Command, args []string) {
		// No args for command
		checkExtraArgs(args)

		// Read config
		lock.Lock()
		defer lock.Unlock()

		yamlFile, err := ioutil.ReadFile(global.ConfigFile)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			panic(err)
		}

		if config == nil {
			content, err := yaml.Marshal(&configParams)
			if err != nil {
				panic(err)
			}

			err = ioutil.WriteFile(global.ConfigFile, content, 0644)
			if err != nil {
				panic(err)
			}
		}

		var isChanged bool
		var resultMessage string

		if cmd.PersistentFlags().Lookup("server").Changed {
			if !strings.EqualFold(config.Server, configParams.Server) {
				config.Server = configParams.Server
				isChanged = true
				resultMessage += "Server changed!\n"
			}
		}

		if cmd.PersistentFlags().Lookup("port").Changed {
			if config.Port != configParams.Port {
				config.Port = configParams.Port
				isChanged = true
				resultMessage += "Port changed!\n"
			}
		}

		if cmd.PersistentFlags().Lookup("user").Changed {
			if !strings.EqualFold(config.User, configParams.User) {
				config.User = configParams.User
				isChanged = true
				resultMessage += "User changed!\n"
			}
		}

		if cmd.PersistentFlags().Lookup("pass").Changed {
			fmt.Println("Please type the password!")

			pass, err := gopass.GetPasswd()
			if err != nil {
				fmt.Println("Error reading password")
			}

			key := "elastic the the best way to use logs"

			decryptPass, err := util.DecryptString(config.Pass, key)
			if err != nil {
				panic("Error when dec string")
			}

			if !strings.EqualFold(string(pass), decryptPass) {
				encryptPass, err := util.EncryptString(string(pass), key)
				if err != nil {
					panic("Error when dec string")
				}
				config.Pass = encryptPass
				isChanged = true
				resultMessage += "Password changed!\n"
			}
		}

		if isChanged {
			output, err := yaml.Marshal(&config)
			if err != nil {
				panic(err)
			}

			f, err := os.Create(global.ConfigFile)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			_, err = f.WriteString(string(output))
			if err != nil {
				panic(err)
			}

			fmt.Println(resultMessage)
		} else {
			fmt.Println("No changes found!")
		}

	},
}

func init() {

	setCmd.PersistentFlags().StringVarP(&configParams.Server, "server", "s", "http://10.96.212.210", "Set kibana server")
	setCmd.PersistentFlags().IntVarP(&configParams.Port, "port", "P", 9200, "Set kibana server port")
	setCmd.PersistentFlags().StringVarP(&configParams.User, "user", "u", "eugen.lupului@metrosystems.net", "Set kibana username")
	// setCmd.PersistentFlags().StringVarP(&configParams.Pass, "pass", "p", "", "Set kibana username password")
	setCmd.PersistentFlags().BoolVarP(&changePass, "pass", "p", false, "Set kibana username password")

}

// func updateConfig(first, second *c.Config) bool {
// 	isUpdated := false

// 	if changePass {
// 		fmt.Println("Please type the password!")

// 		pass, err := gopass.GetPasswd()
// 		if err != nil {
// 			fmt.Println("Error reading password")
// 		}

// 		key := "elastic the the best way to use logs"
// 		encryptPass, err := util.EncryptString(string(pass), key)
// 		if err != nil {
// 			panic("Error when dec string")
// 		}
// 		second.Pass = string(encryptPass)
// 	}

// 	if !strings.EqualFold(first.Server, second.Server) {
// 		first.Server = second.Server
// 		isUpdated = true

// 	}

// 	if first.Port != second.Port {
// 		first.Port = second.Port
// 		isUpdated = true
// 		fmt.Println("Port changed!")
// 	}

// 	if !strings.EqualFold(first.User, second.User) {
// 		first.User = second.User
// 		isUpdated = true
// 		fmt.Println("User changed!")

// 	}

// 	if !isPasswordEqual(first.Pass, second.Pass) && changePass {
// 		first.Pass = second.Pass
// 		isUpdated = true
// 		fmt.Println("Password changed!")
// 	}

// 	return isUpdated
// }

// func isPasswordEqual(password1, password2 string) bool {
// 	key := "elastic the the best way to use logs"

// 	pass1, err := util.DecryptString(password1, key)
// 	if err != nil {
// 		panic("Error when dec string")
// 	}

// 	pass2 := "RMwLmbuqMbEXzOCPXV6iaXAGJ7tofJCw9HM="

// 	if strings.EqualFold(password2, "") {
// 		password2 = "123456789012345678901234567890"
// 		fmt.Println(pass1)

// 		pass2, err = util.DecryptString(password2, key)

// 		if err != nil {
// 			panic("Error when dec string")
// 		}
// 	}

// 	return strings.EqualFold(pass1, pass2)
// }

func checkExtraArgs(args []string) {
	if len(args) == 1 {
		fmt.Println(fmt.Sprintf("Please remove the extra argument %s", args[0]))
	} else if len(args) > 1 {
		fmt.Println(fmt.Sprintf("Please remove the extra arguments %v", args[0:]))
	}
}
