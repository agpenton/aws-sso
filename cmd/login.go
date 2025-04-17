/*
Copyright © 2022 Asdrubal Gonzalez Penton <agpenton@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
	"github.com/julienroland/usg"
	"github.com/spf13/cobra"
)

type Profile struct {
	aws_access_key_id     string
	aws_secret_access_key string
	aws_session_token     string
}

type block struct {
	Profile Profile
}

type Config struct {
	StartUrl    string    `json:"startUrl"`
	Region      string    `json:"region"`
	AccessToken string    `json:"accessToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login command for the AWS SSO.",
	Long:  `A Login command for the AWS SSO, will`,
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
		s.Start()
		if len(os.Args) <= 3 {
			fmt.Println(usg.Get.ExclamationMark, "Please provide a aws profile")
			return
		}
		if timeValidator().Before(time.Now().Local()) {
			s.Suffix = "   Processing data...  \n"
			log.Println(usg.Get.ExclamationMark, "The credentials are Expired")
			if _, err := os.Stat(ssoCacheDir); os.IsNotExist(err) {
				// Directory does not exist, create it
				err := os.MkdirAll(ssoCacheDir, 0755)
				if err != nil {
					fmt.Println("Error creating directory:", err)
					return
				}
				fmt.Println("Directory created successfully!")
			}
			s.Suffix = "  Login to Profile...  \n"
			ssoLogin(profile)
			log.Println(usg.Get.Tick, "Done")
		} else {
			s.Suffix = "  Validate the time of credentials ...  \n"
			timeValidator().Before(time.Now().Local())
			log.Println(usg.Get.Tick, "Success")
		}
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")
	loginCmd.Flags().StringVarP(&profile, "profile", "p", "", "profile to login (required)")
	err := loginCmd.MarkFlagRequired("profile")
	if err != nil {
		return
	}

	rootCmd.AddCommand(loginCmd)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ssoLogin(profile string) string {
	app := "aws"

	arg0 := "sso"
	arg1 := "login"
	arg2 := "--profile"
	arg3 := profile

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	// Print the output
	//log.Println(string(stdout))

	return string(stdout)
}

func timeValidator() time.Time {
	var expirationDate time.Time
	defer reportPanic()

	f, err := os.Open(ssoCacheDir)
	if err != nil {
		fmt.Println(err)
	}
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range files {
		if v.Name() != "botocore-client-id-eu-central-1.json" {
			// Open our jsonFile
			jsonFile, err := os.Open(ssoCacheDir + "/" + v.Name())

			// if we os.Open returns an error then handle it
			if err != nil {
				fmt.Println(err)
			}

			defer jsonFile.Close()
			//log.Println("Successfully Opened ", v.Name())

			byteValue, _ := io.ReadAll(jsonFile)
			var config Config

			err = json.Unmarshal(byteValue, &config)
			check(err)

			expirationDate = config.ExpiresAt
			os.Setenv("AWS_REGION", config.Region)
		}

	}

	return expirationDate.Local()
}
