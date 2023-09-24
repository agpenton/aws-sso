/*
Copyright © 2022 Asdrubal Gonzalez Penton agpenton@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func check(e error) {
	if e != nil {
		log.Println(e)
		//debug.PrintStack()
	}
}

func checkPanic(e error) {
	if e != nil {
		panic(e)
		//debug.PrintStack()
	}
}

func checkFatal(e error) {
	if e != nil {
		log.Fatal(e)
		//debug.PrintStack()
	}
}

func reportPanic() {
	p := recover()
	if p == nil {
		return
	}
	err, ok := p.(error)
	if ok {
		fmt.Println(err)
	} else {
		panic(p)
	}
}

func currentDir() string {
	path, err := os.Getwd()
	check(err)
	return path
}

func exitErr(e error) {
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}
}

var homeDir, _ = os.UserHomeDir()
var awsDir = homeDir + "/.aws/"
var ssoCacheDir = awsDir + "sso/cache/"
var pwdDir = currentDir()
var credentialsFile = "credentials"
var credentialsPath = awsDir + credentialsFile
var profile string
var profiles []string
var accessTempKey, secretTempkey, tempToken string
var err error

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws-sso",
	Short: "This is an app to log in with sso in aws",
	Long:  `This is an app to log in with sso in aws, and create environment variables. Other functionalities are the creation of .envrc for direnv`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	// if err != nil {
	// 	os.Exit(1)
	// }
	exitErr(err)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aws-sso.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
