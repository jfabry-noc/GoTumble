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

func main() {
	// Check for a config file. If not found, prompt for details to create.
	var configuration auth.AuthConfig
	for {
		configuration, err := auth.LoadConfig()
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				fmt.Println("Config file not found.")
				input.PromptConfig()
			} else {
				fmt.Printf("Error type is: %T\n", err)
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
	selection := input.MainMenu(client.Blog)
	fmt.Printf("Selected option: %v\n", selection)
	os.Exit(0)
}
