/*
Copyright © 2022 Asdrubal Gonzalez Penton agpenton@gmail.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/julienroland/usg"
	"github.com/spf13/cobra"
)

type AWSConfig struct {
	SSO struct {
		StartURL  string `toml:"sso_start_url"`
		SSORegion string `toml:"sso_region"`
		Region    string `toml:"region"`
		// StartURL  string `toml:"sso_start_url"`
		AccountID string `toml:"sso_account_id"`
		RoleName  string `toml:"sso_role_name"`
		Output    string `toml:"output"`
	} `toml:"sso"`
}

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "A command to configure to access to AWS",
	Long:  `A Command to create the configuration file to access to AWS.`,
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//	func parInput() {
//		reader := bufio.NewReader(os.Stdin)
//		fmt.Println("Enter the Profile Name: ")
//		text, _ := reader.ReadString('\n')
//		fmt.Println(text)
//	}
func main() {
	var config AWSConfig
	var profileName string
	configPath := path.Join(awsDir, "aws_config.toml")

	// Ask for input data
	fmt.Print(usg.Get.Pointer, " Enter AWS SSO region", usg.Get.ExclamationMark, " ")
	config.SSO.Region = readInput()
	fmt.Print(usg.Get.Pointer, " Enter SSO start URL", usg.Get.ExclamationMark, " ")
	config.SSO.StartURL = readInput()
	fmt.Print(usg.Get.Pointer, " Enter AWS account ID", usg.Get.ExclamationMark, " ")
	config.SSO.AccountID = readInput()
	fmt.Print(usg.Get.Pointer, " Enter AWS SSO role name", usg.Get.ExclamationMark, " ")
	config.SSO.RoleName = readInput()
	fmt.Print(usg.Get.Pointer, " Enter the Output type (json/yaml)", usg.Get.ExclamationMark, " ")
	config.SSO.Output = readInput()
	fmt.Print(usg.Get.Pointer, " Enter the Profile name", usg.Get.ExclamationMark, " ")
	profileName = readInput()

	_, err = os.Stat(configPath)
	if err != nil {
		f, err := os.Create(configPath)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	} else {
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			fmt.Println(usg.Get.Info, "Creating new config file")
		}

	}

	// Load existing config file if it exists
	// if _, err := toml.DecodeFile(configPath, &config); err != nil {
	// 	fmt.Println(usg.Get.Info, "Creating new config file")
	//
	// }

	// Check if input data is different from existing data
	if config.SSO.Region != "" && config.SSO.Region != readInput() {
		// Write new data to file
		// file, err := os.OpenFile("aws_config.toml", os.O_APPEND|os.O_WRONLY, 0644)
		file, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(usg.Get.Cross, "Error writing to config file:", err)
		}
		defer file.Close()

		// [profile staging]
		// sso_start_url = https://juniqe.awsapps.com/start
		// sso_region = eu-central-1
		// sso_account_id = 445858552116
		// sso_role_name = AdministratorAccess
		// region = eu-central-1
		// output = json

		SSORegion := config.SSO.Region

		fmt.Fprintf(file, "\n\n[profile %s]\n", profileName)
		fmt.Fprintf(file, "sso_start_url = \"%s\"\n", config.SSO.StartURL)
		fmt.Fprintf(file, "sso_region = \"%s\"\n", SSORegion)
		fmt.Fprintf(file, "sso_account_id = \"%s\"\n", config.SSO.AccountID)
		fmt.Fprintf(file, "sso_role_name = \"%s\"\n", config.SSO.RoleName)
		fmt.Fprintf(file, "region = \"%s\"\n", config.SSO.Region)
		fmt.Fprintf(file, "output = \"%s\"\n", config.SSO.Output)

		fmt.Println(usg.Get.ExclamationMark, "New data written to config file")
	} else {
		fmt.Println(usg.Get.Info, "No changes made to config file")
	}
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSuffix(input, "\n")
}
