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

func validateBlog(client tumblr.TumblrClient, newBlog string) bool {
	return client.VerifyBlog(newBlog)
}

func forceValidBlog(client tumblr.TumblrClient, inputController input.InputController, configuration auth.AuthConfig) {
	for {
		newBlogId := inputController.UpdateBlogSelection()
		fmt.Printf("New blog: %v\n", newBlogId)
		if validateBlog(client, newBlogId) {
			client.Blog = newBlogId
			configuration.Instance = newBlogId
			input.ConfigUpdate(configuration, false)
			return
		} else {
			fmt.Printf("%v doesn't appear to be a valid blog ID for this account.\n", newBlogId)
		}
	}
}

// Manage loading of the configuration file from the main application loop.
func loadConfig(inputController *input.InputController) auth.AuthConfig {
	configuration, err := auth.LoadConfig()
	for {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				fmt.Println("Config file not found.")
				inputController.PromptConfig("")
				configuration, err = auth.LoadConfig()
			} else {
				fmt.Printf("Error loading configuration: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Printf("Config loaded for: %v\n", configuration.Instance)
		return configuration
	}
}

func main() {
	// Check for a config file. If not found, prompt for details to create.
	var configuration auth.AuthConfig

	// Instantiate an input controller for driving the menu system.
	inputController := input.CreateController()

	// Load the configuration.
	configuration = loadConfig(&inputController)

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
			postContent, tempFile := inputController.CreatePost(configuration.Format, loadEditor())
			if postContent != "" {
				tags := inputController.GetTags()
				postError := client.AddPost(postContent, configuration.Format, tags)
				inputController.PostAftermath(postError, tempFile)
			}
		} else if menuChoice == 2 {
			fmt.Println("Updating blog selection.")
			newBlogId := inputController.UpdateBlogSelection()
			fmt.Printf("New blog: %v\n", newBlogId)
			if validateBlog(client, newBlogId) {
				client.Blog = newBlogId
				configuration.Instance = newBlogId
				input.ConfigUpdate(configuration, false)
			} else {
				fmt.Printf("%v doesn't appear to be a valid blog ID for this account.\n", newBlogId)
			}
		} else if menuChoice == 3 {
			configuration.Format = input.ToggleFormat(configuration.Format)
			input.ConfigUpdate(configuration, false)
		} else if menuChoice == 4 {
			fmt.Println("Overwriting the entire config file.")
			inputController.PromptConfig(configuration.Format)

			// Reload the configuration.
			configuration = loadConfig(&inputController)

			// Confirm the specified blog is legitimate for the current user.
			if validateBlog(client, configuration.Instance) {
				client.Blog = configuration.Instance
				continue
			} else {
				fmt.Printf("Blog of %v does not appear to be valid for this user.\n", configuration.Instance)
				forceValidBlog(client, inputController, configuration)
			}
		} else if menuChoice == 5 {
			input.UpdateEditorInstr()
		} else {
			fmt.Println("Goodbye!")
			os.Exit(0)
		}
	}
}
