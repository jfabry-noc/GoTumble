package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/jfabry-noc/GoTumble/pkg/auth"
	"github.com/jfabry-noc/GoTumble/pkg/input"
)

func main() {
	// Need to have this check for a config file. If not found, prompt for details
	// to create.
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
	os.Exit(0)

}
