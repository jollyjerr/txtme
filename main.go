package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

type tuser struct {
	Name      string
	PhoneTo   string
	PhoneFrom string
	Sid       string
	Token     string
}

func main() {
	setup()
	commands()

	err := app.Run(os.Args)
	check(err)
}

func setup() {
	app.Name = "txtme"
	app.Usage = "Send yourself a text message when a long script finishes."
	app.Version = "0.0.1"
	app.Action = send
	app.EnableBashCompletion = true
}

func commands() {
	app.Commands = []*cli.Command{} // future spot for commands
}

func send(c *cli.Context) error {
	user := generateUser()
	message := fmt.Sprintf("Yo %s, your script is finished! 🚀", user.Name)

	// Set account keys & information
	accountSid := user.Sid
	authToken := user.Token
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	// Pack up the data for our message
	msgData := url.Values{}
	msgData.Set("To", user.PhoneTo)
	msgData.Set("From", user.PhoneFrom)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			message := fmt.Sprintf("🚀 Yo %s, your script is finished! ", user.Name)
			fmt.Println(message)
		}
	} else {
		message = fmt.Sprintf("💩 Sorry %s, something went wrong!", user.Name)
		log.Print(resp.Status)
		fmt.Println(message)
	}
	return nil
}

func generateUser() tuser {
	user, perr := user.Current()
	check(perr)
	configPath := fmt.Sprintf("%s/.txtme.toml", user.HomeDir)

	var configuredUser tuser
	if _, err := toml.DecodeFile(configPath, &configuredUser); err != nil || configuredUser.Name == "" {
		name := askfor("Name")
		phoneTo := askfor("PhoneTo")
		phoneFrom := askfor("PhoneFrom")
		sid := askfor("Sid")
		token := askfor("Token")

		generatedUser := tuser{
			Name:      name,
			PhoneTo:   phoneTo,
			PhoneFrom: phoneFrom,
			Sid:       sid,
			Token:     token,
		}

		save(generatedUser, configPath)
		return generatedUser
	}
	return configuredUser
}

func askfor(item string) string {
	result := os.Getenv(item)
	if result == "" {
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("Hey 👋! What is your %s? I will save it for future use!", strings.ToLower(item)),
		}
		response, err := prompt.Run()
		check(err)
		result = response
	}

	return result
}

func save(user tuser, filePath string) {
	fmt.Println(filePath)
	f, _ := os.Create(filePath)
	if err := toml.NewEncoder(f).Encode(user); err != nil {
		log.Println(err)
	}
	err := f.Close()
	check(err)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
