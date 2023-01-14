package input

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jfabry-noc/GoTumble/pkg/auth"
	"github.com/jfabry-noc/GoTumble/pkg/postfile"
)

type InputController struct {
	InputBuffer *bufio.Reader
}

// CreateController instantiates a new InputController struct.
func CreateController() InputController {
	var controllerInstance InputController
	controllerInstance.InputBuffer = bufio.NewReader(os.Stdin)
	return controllerInstance
}

// getInput gathers user input from a prompt and returns the string sans newline.
func (i *InputController) getInput(message string) string {
	fmt.Println(message)
	fmt.Print("> ")
	userInput, err := i.InputBuffer.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to retrieve input with error: %v\n", err)
		os.Exit(1)
	}
	return strings.Trim(userInput, "\n")
}

func (i *InputController) Printer(message string) {
	fmt.Println(message)
}

func ToggleFormat(currentFormat string) string {
	if currentFormat == "html" {
		return "markdown"
	} else {
		return "html"
	}
}

// PromptConfig gets configuration details from the user.
func (i *InputController) PromptConfig(currentFormat string) {
	var config auth.AuthConfig
	fmt.Println("Please provide the following configuration information.")
	fmt.Println("This will be written to: ~/.config/gotumble/config.json")

	config.ConsumerKey = i.getInput("Consumer Key: ")
	config.ConsumerSecret = i.getInput("Consumer Secret: ")
	config.Token = i.getInput("Token: ")
	config.TokenSecret = i.getInput("Token Secret: ")
	config.Instance = i.getInput("Instance: ")

	if currentFormat == "" {
		config.Format = "html"
	} else {
		config.Format = ToggleFormat(currentFormat)
	}

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
	// Println generates a warning for the extra newline that I want for formatting.
	fmt.Printf("~/.zshrc or ~/.bashrc\n\n")
}

// UpdateBlogSelection modifies the current blog used for posts.
func (i *InputController) UpdateBlogSelection() string {
	return i.getInput("Enter blog ID.")
}

// MainMenu prints the main menu and gets user input on where to navigate.
func (i *InputController) MainMenu(blogName string, editor string, format string) int {
	fmt.Printf("--== Posting to %v in %v with editor: %v ==--\n", blogName, format, editor)
	fmt.Println("1. New text post")
	fmt.Println("2. New link post")
	fmt.Println("3. New video post.")
	fmt.Println("4. New quote post.")
	fmt.Println("5. Update blog selection")
	fmt.Println("6. Toggle format (HTML or Markdown)")
	fmt.Println("7. Overwrite config file")
	fmt.Println("8. View editor instructions")
	fmt.Println("9. Quit")

	var choice string
	var choiceInt int
	var err error
	for {
		choice = i.getInput("What would you like to do?")
		choiceInt, err = strconv.Atoi(choice)
		if err != nil {
			fmt.Printf("%v is not a valid choice. Please enter a number.\n", choice)
			continue
		}

		if choiceInt < 1 || choiceInt > 9 {
			fmt.Printf("%v is not a valid choice. Please select from the options above.\n", choiceInt)
			continue
		}

		// Break if both tests above passed.
		break
	}

	return choiceInt
}

// UpdateFormat will prompt the user for the new post format (HTML or Markdown)
func (i *InputController) UpdateFormat() string {
	choice := i.getInput("Select if you'd like HTML or Markdown for your posts.")

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

// getInputChoice forces input to be one of two choices.
func (i *InputController) getInputChoice(message string, first string, second string) string {
	for {
		response := i.getInput(message)
		if strings.ToLower(response) == first || strings.ToLower(response) == second {
			return strings.ToLower(response)
		}
	}
}

// getPostContent Gets the content from the user for any text-based post.
func (i *InputController) getPostContent(postFormat string, editorPath string) (string, string) {
	tempFilePath := postfile.PostFilePath(postFormat)
	for {
		fmt.Printf("Opening temp file at '%v' with editor '%v'...\n", tempFilePath, editorPath)
		time.Sleep(2 * time.Second)
		fileCreated := postfile.CreateFile(tempFilePath, editorPath)
		if fileCreated != nil {
			fmt.Printf("Failed to create the file with error: %v\n", fileCreated)
			return "", tempFilePath
		}

		fileContent, err := postfile.ReadFile(tempFilePath)
		if err != nil {
			fmt.Printf("Failed to read the temporary file with error: %v\n", err)
			return "", tempFilePath
		}

		fmt.Printf("Preparing to post the following content:\n\n")
		fmt.Printf("%v\n", fileContent)

		message := "Post this? Y or N"
		response := i.getInputChoice(message, "y", "n")

		if response == "y" {
			fmt.Println("Posting!")
			return fileContent, tempFilePath
		}

		message = "Edit or discard? E or D"
		response = i.getInputChoice(message, "e", "d")

		if response == "e" {
			fmt.Println("Re-editing the post.")
		} else {
			fmt.Println("Discarding the post!")
			deletionError := postfile.DeleteFile(tempFilePath)
			if deletionError != nil {
				fmt.Printf("Failed to delete the file with error: %v\n", deletionError)
			}
			fmt.Printf("Successfully removed the file at: %v\n", tempFilePath)
			return "", tempFilePath
		}
	}
}

// CreateTextPost Wrapper function to just get text post content.
func (i *InputController) CreateTextPost(postFormat string, editorPath string) (string, string) {
	content, file := i.getPostContent(postFormat, editorPath)

	return content, file
}

func (i *InputController) CreateQuotePost(postFormat string, editorPath string) (string, string, string, string) {
	fmt.Println("Remember that quotes and their sources can only be HTML, not Markdown!")
	fmt.Println("Enter the quote first into the file.")
	fmt.Print("Enter to continue...")
	_, _ = i.InputBuffer.ReadString('\n')
	quote, quoteFile := i.getPostContent(postFormat, editorPath)

	fmt.Println("Enter the source, if any, into this file.")
	fmt.Print("Enter to continue...")
	source, sourceFile := i.getPostContent(postFormat, editorPath)

	return quote, source, quoteFile, sourceFile
}

// validateUrl Checks if a URL is valid. If not, an error is returned.
func validateUrl(link string) error {
	_, err := url.ParseRequestURI(link)

	return err
}

// CreateVideoPost Wrapper method to get a URL for a video and optional caption.
func (i *InputController) CreateVideoPost(postFormat string, editorPath string) (string, string, string) {
	// Not sure if this is true... test it.
	fmt.Println("Remember that video captions can only be HTML, not Markdown!")
	link := i.getUrl("Enter video URI to share.")
	if link == "" {
		return "", "", ""
	}

	content, file := i.getPostContent(postFormat, editorPath)

	return link, content, file

}

// getUrl Gathers a valid URL from the user or Q to quit.
func (i *InputController) getUrl(message string) string {
	var link string
	for {
		link = i.getInput(message)

		if strings.ToLower(link) == "q" {
			return ""
		}

		err := validateUrl(link)

		if err == nil {
			break
		} else {
			fmt.Printf("%v was not a valid URL. Try again or enter 'Q' to quit.\n", link)
		}
	}

	return link
}

// CreateLinkPost Wrapper function to get a URL and optional description for a link.
func (i *InputController) CreateLinkPost(postFormat string, editorPath string) (string, string, string) {
	fmt.Println("Remember that link descriptions can only be HTML, not Markdown!")
	link := i.getUrl("Enter the URL to share.")
	if link == "" {
		return "", "", ""
	}
	/*
		var link string
		for {
			link = i.getInput("Enter the URL to share.")

			if strings.ToLower(link) == "q" {
				return "", "", ""
			}

			err := validateUrl(link)

			if err == nil {
				break
			} else {
				fmt.Printf("%v was not a valid URL. Try again or enter 'Q' to quit.\n", link)
			}
		}
	*/
	content, file := i.getPostContent(postFormat, editorPath)

	return link, content, file
}

// PostAftermath prints if a post was successful and cleans up the temporary file.
func (i *InputController) PostAftermath(err error, filePaths []string) {
	if err != nil {
		fmt.Printf("Failed to create post with error: %v\n", err)
	} else {
		fmt.Println("Post added successfully!")
	}

	for _, filePath := range filePaths {
		fmt.Printf("Cleaning up temporary file at: %v\n", filePath)
		err = postfile.DeleteFile(filePath)
		if err != nil {
			fmt.Printf("Failed to delete temporary file with error: %v\n", err)
		}
	}
}

// promptSecrets will prompt the user for each for each of the config file's secrets.
func (i *InputController) promptSecrets(secretName string) string {
	message := fmt.Sprintf("Enter the value for the new: %v\n", secretName)
	return i.getInput(message)
}

func (i *InputController) GetTags() string {
	message := "Enter tags, if any, separated by commas."
	rawTags := i.getInput(message)
	if len(rawTags) > 0 {
		return rawTags
	}
	return ""
}
