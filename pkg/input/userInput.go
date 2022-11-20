package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jfabry-noc/GoTumble/pkg/auth"
)

// getInput gathers user input from a prompt and returns the string sans newline.
func getInput(message string, buffer *bufio.Reader) string {
	fmt.Println(message)
	fmt.Print("> ")
	userInput, err := buffer.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to retrieve input with error: %v\n", err)
		os.Exit(1)
	}
	return strings.Trim(userInput, "\n")
}

// PromptConfig gets configuration details from the user.
func PromptConfig() {
	var config auth.AuthConfig
	buffer := bufio.NewReader(os.Stdin)
	fmt.Println("Please provide the following configuration information.")
	fmt.Println("This will be written to: ~/.config/gotumble.json")

	config.ConsumerKey = getInput("Consumer Key: ", buffer)
	config.ConsumerSecret = getInput("Consumer Secret: ", buffer)
	config.Token = getInput("Token: ", buffer)
	config.TokenSecret = getInput("Token Secret: ", buffer)
	config.Instance = getInput("Instance: ", buffer)

	result, err := auth.WriteConfig(config)
	fmt.Print(result)
	if err != nil {
		os.Exit(1)
	}

}

// MainMenu prints the main menu and gets user input on where to navigate.
func MainMenu(blogName string) int {
	buffer := bufio.NewReader(os.Stdin)
	fmt.Printf("Current using blog: %v\n", blogName)
	fmt.Println("1. New post")
	fmt.Println("2. Update blog selection")
	fmt.Println("3. Overwrite config file")
	fmt.Println("4. Quit")

	var response string
	for {
		response = getInput("What would you like to do?", buffer)

		if response != "1" && response != "2" && response != "3" && response != "4" {
			fmt.Printf("%v is not a validate response! Please enter 1 - 4.\n", response)
		} else {
			break
		}
	}

	// Take the appropriate action.
	if response == "1" {
		fmt.Println("Creating new post.")
	} else if response == "2" {
		fmt.Println("Updating blog selection.")
	} else if response == "3" {
		fmt.Println("Overwriting the current config file.")
	} else if response == "4" {
		fmt.Println("Quitting the application")
		os.Exit(0)
	}

	intValue, _ := strconv.Atoi(response)
	return intValue
}
