package main

import (
	"github.com/spf13/viper"
	"os"
	"fmt"
	"github.com/spf13/cobra"
)

const logRequest bool = true
const logResponse bool = true

var CfgFile string


// Defining the proxy command.
// The default behavior is to start the proxy
var RootCmd = &cobra.Command{
	Use:   "proxy",
	Short: "a multi tenancy docker proxy",
	Long: `a multi tenancy docker proxy`,
	Run: func(cmd *cobra.Command, args []string) {

		var dockerEnpoint string
		var localAddrString string

		dockerEnpoint = viper.GetString("docker-endpoint")
		localAddrString = viper.GetString("address")

		pConf := &ProxyConfiguration{DockerEndpoint: dockerEnpoint,
									 Address: localAddrString}
		StartProxy(pConf)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (default is $CWD/config.yaml)")
	RootCmd.PersistentFlags().StringP("address", "a", ":9000", "address to listen to")
	RootCmd.PersistentFlags().StringP("docker-endpoint", "d", "unix:///var/run/docker.sock", "docker host endpoint")

	// Bind viper to these flags so viper can read flag values along with config, env, etc.
	viper.BindPFlag("address", RootCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag("docker-endpoint", RootCmd.PersistentFlags().Lookup("docker-endpoint"))
}

// Read in config file and ENV variables if set.
func initConfig() {
	if CfgFile != "" {
		viper.SetConfigFile(CfgFile)
	}

	dir, _ := os.Getwd()
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(dir)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config.yaml file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}


func main(){

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

