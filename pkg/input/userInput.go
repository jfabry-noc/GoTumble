package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jfabry-noc/GoTumble/pkg/auth"
)

type InputController struct {
	InputBuffer *bufio.Reader
}

func CreateController() InputController {
	var controllerInstance InputController
	controllerInstance.InputBuffer = bufio.NewReader(os.Stdin)
	return controllerInstance
}

// getInput gathers user input from a prompt and returns the string sans newline.
func (b *InputController) getInput(message string) string {
	fmt.Println(message)
	fmt.Print("> ")
	userInput, err := b.InputBuffer.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to retrieve input with error: %v\n", err)
		os.Exit(1)
	}
	return strings.Trim(userInput, "\n")
}

// PromptConfig gets configuration details from the user.
func (b *InputController) PromptConfig() {
	var config auth.AuthConfig
	fmt.Println("Please provide the following configuration information.")
	fmt.Println("This will be written to: ~/.config/gotumble.json")

	config.ConsumerKey = b.getInput("Consumer Key: ")
	config.ConsumerSecret = b.getInput("Consumer Secret: ")
	config.Token = b.getInput("Token: ")
	config.TokenSecret = b.getInput("Token Secret: ")
	config.Instance = b.getInput("Instance: ")

	ConfigUpdate(config, true)
}

// ConfigUpdate writes an updated configuration file after the user has specified changes.
func ConfigUpdate(config auth.AuthConfig, hardStop bool) {
	fmt.Println("Updating config file...")
	result, err := auth.WriteConfig(config)
	fmt.Print(result)
	if err != nil {
		fmt.Printf("Error message was: %v\n", err)
		if hardStop {
			os.Exit(1)
		}
	}
}

// updateEditorInstr prints instructions for how to update the editor.
func UpdateEditorInstr() {
	fmt.Println("The editor seen on the screen is based on the $EDITOR environment variable.")
	fmt.Println("This can be temporarily updated for your current shell by running:")
	fmt.Println("\nexport EDITOR=\"/path/to/editor\"")
	fmt.Println("\nTo set this permanently, set it in your shell's config, e.g.")
	fmt.Println("~/.zshrc or ~/.bashrc\n")
}

// UpdateBlogSelection modifies the current blog used for posts.
func (b *InputController) UpdateBlogSelection() string {
	return b.getInput("Enter blog ID.")
}

// MainMenu prints the main menu and gets user input on where to navigate.
func (b *InputController) MainMenu(blogName string, editor string, format string) int {
	fmt.Printf("--== Posting to %v in %v with editor: %v ==--\n", blogName, format, editor)
	fmt.Println("1. New post")
	fmt.Println("2. Update blog selection")
	fmt.Println("3. Toggle format (HTML or Markdown)")
	fmt.Println("4. Overwrite config file")
	fmt.Println("5. View editor instructions")
	fmt.Println("6. Quit")

	var choice string
	var choiceInt int
	var err error
	for {
		choice = b.getInput("What would you like to do?")
		choiceInt, err = strconv.Atoi(choice)
		if err != nil {
			fmt.Printf("%v is not a valid choice. Please enter a number.\n", choice)
			continue
		}

		if choiceInt < 1 || choiceInt > 6 {
			fmt.Printf("%v is not a valid choice. Please select from the options above.\n", choiceInt)
			continue
		}

		// Break if both tests above passed.
		break
	}

	return choiceInt
}

// UpdateFormat will prompt the user for the new post format (HTML or Markdown)
func (b *InputController) UpdateFormat() string {
	choice := b.getInput("Select if you'd like HTML or Markdown for your posts.")

	choiceLowered := strings.ToLower(choice)
	for {
		if choiceLowered == "m" || choiceLowered == "markdown" {
			return "markdown"
		} else if choiceLowered == "h" || choiceLowered == "html" {
			return "html"
		} else if choiceLowered == "q" || choiceLowered == "quit" {
			return "quit"
		}
	}
}
