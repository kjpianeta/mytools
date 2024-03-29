package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kpianeta/mytools/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var settings = new(config.Config)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mytools",
	Short: "Various tools used daily",
	Long: `mytools is designed to be used to perfomr common tasks.

This usually involves tasks that would require multiple calls.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)
	settings.Verbose = RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Give verbose output")
	settings.OutputFile = RootCmd.PersistentFlags().StringP("file", "f", "", "Optional file to save the output to")
	settings.OutputFormat = RootCmd.PersistentFlags().StringP("output", "o", "json", "Format for the output, currently supported are csv, json, html, dot, and drawio")
	settings.AppendToOutput = RootCmd.PersistentFlags().BoolP("append", "a", false, "Add to the provided output file instead of replacing it")
	settings.NameFile = RootCmd.PersistentFlags().StringP("namefile", "n", "", "Use this file to provide names")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName(".mytools") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getName(id string) string {
	if *settings.NameFile != "" {
		nameFile, err := ioutil.ReadFile(*settings.NameFile)
		if err != nil {
			panic(err)
		}
		values := make(map[string]string)
		err = json.Unmarshal(nameFile, &values)
		if err != nil {
			panic(err)
		}
		if val, ok := values[id]; ok {
			return val
		}
	}
	return id
}
