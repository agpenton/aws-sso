/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	_ "bytes"
	"fmt"
	"io/ioutil"
	//"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"

	"github.com/spf13/cobra"
)

// Function to check if command exist
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Function to get the current directory.
//func currentDir() string {
//	path, err := os.Getwd()
//	check(err)
//	return path
//}

// Loading the data from .envrc file.
func loadEnvrc() {
	err := godotenv.Load(".envrc")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	awsProfile := os.Getenv("AWS_PROFILE")

	log.Printf("The Profile is: %v", awsProfile)
}

func envrcVars() string {
	var envrc = []string{
		fmt.Sprintf("export AWS_PROFILE=\"%v\"", os.Getenv("AWS_PROFILE")),
		fmt.Sprintf("export AWS_ACCESS_KEY_ID=\"%v\"", os.Getenv("AWS_ACCESS_KEY_ID")),
		fmt.Sprintf("export AWS_SECRET_ACCESS_KEY=\"%v\"", os.Getenv("AWS_SECRET_ACCESS_KEY")),
		fmt.Sprintf("export AWS_SESSION_TOKEN=\"%v\"", os.Getenv("AWS_SESSION_TOKEN")),
		fmt.Sprintf("export AWS_REGION=\"%v\"", os.Getenv("AWS_REGION")),
	}
	output := strings.Join(envrc, "\n")

	return output
}

// Writing the data inside the .envrc file
func envFile() {
	filename := ".envrc"
	var _, err = os.Stat(filename)

	output := envrcVars()
	if os.IsNotExist(err) {
		log.Printf("Creating the file %v", filename)
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		checkFatal(err)
		if _, err := f.Write([]byte(output)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("The file %s already exists!\n", filename)
		loadEnvrc()
		modifyEnvrc()
		Shellout("direnv allow")
		return
	}

	log.Println("File created successfully", pwdDir+"/"+filename)
}

// Modify the file if exist.
func modifyEnvrc() {
	filename := ".envrc"
	file := pwdDir + "/" + filename
	output := envrcVars()
	err := ioutil.WriteFile(file, []byte(output), 0644)
	checkFatal(err)
}

const ShellToUse = "bash"

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

// envrcCmd represents the envrc command
var envrcCmd = &cobra.Command{
	Use:   "envrc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//accessTempKey, secretTempkey, tempToken, nil := CreateSession(profile)
		fmt.Println("envrc called")

		if commandExists("direnv") == true {
			fmt.Println("you are here", currentDir())
			fmt.Println("Do you want to create a .envrc? y/n")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				return
			}

			if input == "y" {
				envFile()
				log.Println("The temporary credentials were added to the .envrc file")
			} else {
				fmt.Println("----------------------------------")
				fmt.Println("AWS Profile", profile)
				fmt.Println("Access Key: ", os.Getenv("AWS_ACCESS_KEY_ID"))
				fmt.Println("Secret Key: ", os.Getenv("AWS_SECRET_ACCESS_KEY"))
				fmt.Println("Session Token: ", os.Getenv("AWS_SESSION_TOKEN"))
				fmt.Println("----------------------------------")
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(envrcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envrcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envrcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
