package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/jfabry-noc/GoTumble/pkg/auth"
	"github.com/jfabry-noc/GoTumble/pkg/input"
	"github.com/jfabry-noc/GoTumble/pkg/tumblr"
)

const editorVar = "EDITOR"

func loadEditor() string {
	return os.Getenv(editorVar)
}

func main() {
	// Check for a config file. If not found, prompt for details to create.
	var configuration auth.AuthConfig

	// Instantiate an input controller for driving the menu system.
	inputController := input.CreateController()

	// This is a problem because configuration is NOT updated outside of the
	// scope of this for loop. Meaning the call to configuration.Instance
	// immediately after it returns the default of empty string for that value...
	// Probably need another function that returns the solid values?
	configuration, err := auth.LoadConfig()
	for {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				fmt.Println("Config file not found.")
				inputController.PromptConfig()
				configuration, err = auth.LoadConfig()
			} else {
				fmt.Printf("Error loading configuration: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Printf("Config loaded for: %v\n", configuration.Instance)
		break
	}

	// Instantiate a Tumblr Client from the library.
	client := tumblr.CreateClient(
		configuration.ConsumerKey,
		configuration.ConsumerSecret,
		configuration.Token,
		configuration.TokenSecret,
		configuration.Instance,
	)

	// Start the main loop.
	for {
		menuChoice := inputController.MainMenu(client.Blog, loadEditor(), configuration.Format)
		if menuChoice == 1 {
			fmt.Println("Creating new post.")
		} else if menuChoice == 2 {
			fmt.Println("Updating blog selection.")
			newBlogId := inputController.UpdateBlogSelection()
			// Delete this later. Just leaving for the fmt import.
			fmt.Printf("New blog: %v\n", newBlogId)
			if client.VerifyBlog(newBlogId) {
				// NEXT: Need to figure out how to write this to the config file
				// along with having it update the current struct.
				client.Blog = newBlogId
				configuration.Instance = newBlogId
				input.ConfigUpdate(configuration, false)
			} else {
				fmt.Printf("%v doesn't appear to be a valid blog ID for this account.\n", newBlogId)
			}
		} else if menuChoice == 3 {
			if configuration.Format == "markdown" {
				configuration.Format = "html"
			} else if configuration.Format == "html" {
				configuration.Format = "markdown"
			}
			input.ConfigUpdate(configuration, false)
		} else if menuChoice == 4 {
			fmt.Println("Overwriting the entire config file.")
		} else if menuChoice == 5 {
			input.UpdateEditorInstr()
		} else {
			fmt.Println("Goodbye!")
			os.Exit(0)
		}
	}
}
