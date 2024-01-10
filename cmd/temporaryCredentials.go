/*
Copyright © 2022 Asdrubal Gonzalez Penton <agpenton@gmail.com>
*/
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/briandowns/spinner"
	"github.com/julienroland/usg"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// temporaryCredentialsCmd represents the temporaryCredentials command
var temporaryCredentialsCmd = &cobra.Command{
	Use:   "credentials",
	Short: "This is a command to create or modify the credentials created.",
	Long:  `This is a command to create or modify the credentials in the config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[11], 100*time.Millisecond) // Build our new spinner
		s.Suffix = "   Processing data...  \n"
		s.Start()
		accessTempKey, secretTempkey, tempToken, nil := CreateSession(profile)
		_, err = os.Stat(credentialsPath)
		if err != nil {

			//if os.IsNotExist(err) {
			if os.IsNotExist(err) {
				s.Suffix = "  Login to Profile ...  \n"
				ssoLogin(profile)
				log.Println(usg.Get.Tick, "Success")
				credentialsFileCreation(profile, accessTempKey, secretTempkey, tempToken)

			}
		} else {
			if timeValidator().Before(time.Now().Local()) {
				//log.Println("The credentials are Expired")
				e := os.Remove(credentialsPath)
				if e != nil {
					log.Fatal(e)
				}
				s.Suffix = "  Login to Profile.. \n"
				ssoLogin(profile)
				log.Println(usg.Get.Tick, "Success")
				s.Suffix = "  Updating Credentials.. \n"
				UpdateCredentials(profile, accessTempKey, secretTempkey, tempToken)
				log.Println(usg.Get.Tick, "Done")
			} else {
				s.Suffix = "  Updating Credentials.. \n"
				UpdateCredentials(profile, accessTempKey, secretTempkey, tempToken)
				log.Println(usg.Get.Tick, "Done")
			}
		}
	},
}

func init() {
	temporaryCredentialsCmd.Flags().StringVarP(&profile, "profile", "p", "", "profile to login (required)")
	err := temporaryCredentialsCmd.MarkFlagRequired("profile")
	if err != nil {
		return
	}

	rootCmd.AddCommand(temporaryCredentialsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// temporaryCredentialsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// temporaryCredentialsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CreateSession(profile string) (string, string, string, error) {
	//sess, err := session.NewSessionWithOptions(session.Options{
	//	SharedConfigState: session.SharedConfigEnable,
	//	Profile:           profile,
	//})
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	}))

	credentials, err := sess.Config.Credentials.Get()
	check(err)

	accessTempKey := credentials.AccessKeyID
	secretTempkey := credentials.SecretAccessKey
	tempToken := credentials.SessionToken

	os.Setenv("AWS_ACCESS_KEY_ID", accessTempKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", secretTempkey)
	os.Setenv("AWS_SESSION_TOKEN", tempToken)
	os.Setenv("AWS_PROFILE", profile)

	//checkFatal(err)

	return accessTempKey, secretTempkey, tempToken, nil
}

// UpdateRegistryConfig - updates registry settings in the config file
func UpdateCredentials(credentialsPath, accessTempKey string, secretTempkey string, tempToken string) error {
	log.Debugf(usg.Get.ExclamationMark, "UpdateCredentials hit, credential file: %s", credentialsPath)
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		return fmt.Errorf(usg.Get.Cross, "specified credentials file %s not exists locally, error: %v", credentialsPath, err)
	}

	tree, err := getTomlTree(credentialsPath)
	if err != nil {
		return fmt.Errorf(usg.Get.Cross, "failed to load the credentials file as toml tree, error: %v", err)
	}

	// auth
	tree.SetPath([]string{profile, "aws_access_key_id"}, accessTempKey)
	tree.SetPath([]string{profile, "aws_secret_access_key"}, secretTempkey)
	tree.SetPath([]string{profile, "aws_session_token"}, tempToken)

	if err := persistTomlTree(credentialsPath, tree); err != nil {
		return err
	}

	log.Debug(usg.Get.Tick, "credentials settings added successfully")
	return nil
}

// AddAnotherConfig - adds nvidia settings in config file
func AddAnotherConfig(credentialsPath, accessTempKey string, secretTempkey string, tempToken string) error {
	if _, err := os.Stat(credentialsPath); os.IsNotExist(err) {
		return fmt.Errorf("specified config file %s not exists locally, error: %v", credentialsPath, err)
	}

	tree, err := getTomlTree(credentialsPath)
	if err != nil {
		return fmt.Errorf("failed to load the config file as toml tree, error: %v", err)
	}

	// auth
	tree.SetPath([]string{profile, "aws_access_key_id"}, accessTempKey)
	tree.SetPath([]string{profile, "aws_secret_access_key"}, secretTempkey)
	tree.SetPath([]string{profile, "aws_session_token"}, tempToken)

	if err := persistTomlTree(credentialsPath, tree); err != nil {
		return err
	}

	log.Debug(usg.Get.Tick, "credentials settings added successfully")
	return nil
}

func getTomlTree(configFile string) (*toml.Tree, error) {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file, error: %v", err)
	}

	content := string(bytes)
	tree, err := toml.Load(content)
	if err != nil {
		return nil, fmt.Errorf("toml library failed to load config from the file content, error: %v", err)
	}
	return tree, nil
}

func persistTomlTree(configFile string, tree *toml.Tree) error {
	str, err := tree.ToTomlString()
	if err != nil {
		return fmt.Errorf("toml library failed to convert config to a string, error: %v", err)
	}

	data := []byte(str)
	if err := ioutil.WriteFile(configFile, data, 0666); err != nil {
		return fmt.Errorf("failed to write config to a file, error: %v", err)
	}

	log.Info(usg.Get.Tick, "updates written to file successfully")
	return nil
}

// SearchString to search for the credentials in the config file and convert them to variables.
func SearchString(profile string) Profile {
	defer reportPanic()

	// Loading the data from the toml file.
	creds, _ := toml.LoadFile(credentialsPath)

	access := creds.Get(fmt.Sprintf("%v.aws_access_key_id", profile)).(string)
	secret := creds.Get(fmt.Sprintf("%v.aws_secret_access_key", profile)).(string)
	token := creds.Get(fmt.Sprintf("%v.aws_session_token", profile)).(string)

	// Return the values from the function.
	return Profile{
		aws_access_key_id:     access,
		aws_secret_access_key: secret,
		aws_session_token:     token,
	}
}

func credentialsFileCreation(profile string, accessTempKey string, secretTempkey string, tempToken string) {
	var _, err = os.Stat(credentialsPath)

	if os.IsNotExist(err) {
		file, err := os.Create(credentialsPath)
		check(err)
		defer file.Close()

		aak := fmt.Sprintf("[%v]\naws_access_key_id = %v\n", profile, accessTempKey)
		_, err = file.WriteString(aak)
		check(err)
		err = file.Sync()
		check(err)
		w := bufio.NewWriter(file)
		asak := fmt.Sprintf("aws_secret_access_key = %v\naws_session_token = %v", secretTempkey, tempToken)
		_, err = w.WriteString(asak)
		check(err)
		err = w.Flush()
		check(err)
	} else {
		log.Printf(usg.Get.ExclamationMark, "The file %v already exists!\n", credentialsFile)
		log.Println(usg.Get.ExclamationMark, "modifying the values")
		ModifyCredentials(profile, accessTempKey, secretTempkey, tempToken)
		log.Println(usg.Get.Tick, "done")
		return
	}

	log.Println(usg.Get.Tick, "File created successfully", credentialsPath)
}

// ModifyCredentials to Modify the credentials in the file if they exist.
func ModifyCredentials(profile string, accessTempKey string, secretTempkey string, tempToken string) {

	var tempCredentials = []string{
		fmt.Sprintf("[%v]", profile),
		fmt.Sprintf("aws_access_key_id = \"%v\"", accessTempKey),
		fmt.Sprintf("aws_secret_access_key = \"%v\"", secretTempkey),
		fmt.Sprintf("aws_session_token = \"%v\"", tempToken),
	}

	output := strings.Join(tempCredentials, "\n")
	//err := ioutil.WriteFile(credentialsPath, []byte(output), 0644
	err := os.WriteFile(credentialsPath, []byte(output), 0644)
	checkFatal(err)
	return
}
