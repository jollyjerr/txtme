package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

type tuser struct {
	name  string
	phone string
	sid   string
	token string
}

func main() {
	setup()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setup() {
	app.Name = "txtme"
	app.Usage = "Send yourself a text message when a long script finishes."
	app.Version = "0.0.1"
	app.Action = send
	app.EnableBashCompletion = true
}

func commands() {
	app.Commands = []*cli.Command{
		// {
		// 	Name: "Set SID"

		// }
	}
}

func send(c *cli.Context) error {
	user := generateUser()
	str := fmt.Sprintf("🚀 Yo %s, your script is finished! ", user.name)

	// Set account keys & information
	accountSid := user.sid
	authToken := user.token
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	// Pack up the data for our message
	msgData := url.Values{}
	msgData.Set("To", "NUMBER_TO")
	msgData.Set("From", "NUMBER_FROM")
	msgData.Set("Body", str)
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
			str := fmt.Sprintf("🚀 Yo %s, your script is finished! ", user.name)
			fmt.Println(str)
		}
	} else {
		str = fmt.Sprintf("💩 Sorry %s, something went wrong!", user.name)
		log.Print(resp.Status)
		fmt.Println(str)
	}
	return nil
}

func generateUser() tuser {
	name := locate("TXTME_USER_NAME")
	phone := locate("TXTME_USER_PHONE")
	sid := locate("TXTME_USER_SID")
	token := locate("TXTME_USER_TOKEN")

	generatedUser := tuser{
		name:  name,
		phone: phone,
		sid:   sid,
		token: token,
	}

	return generatedUser
}

func locate(item string) string {
	result := os.Getenv(item)
	if result == "" {
		str := fmt.Sprintf("😮 Uh oh! I could not find your %s anywhere!", item)
		fmt.Println(str)
		prompt := promptui.Prompt{
			Label: fmt.Sprintf("What is your %s? I will save it for future use!", item),
		}
		response, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		result = response
	}
	return result
}
