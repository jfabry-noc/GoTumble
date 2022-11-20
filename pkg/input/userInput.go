package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jfabry-noc/GoTumble/pkg/auth"
)

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
