package postfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

// NewPost creates the path for the temporary file to hold the post content.
func PostFilePath(format string) string {
	basePath := "/tmp/gotumble"
	rightNow := time.Now()
	filePath := fmt.Sprintf(
		"%v_%v%v%v_%v%v%v",
		basePath,
		rightNow.Year(),
		rightNow.Month(),
		rightNow.Day(),
		rightNow.Hour(),
		rightNow.Minute(),
		rightNow.Second(),
	)
	if format == "markdown" {
		filePath = fmt.Sprintf("%v%v", filePath, ".md")
	} else {
		filePath = fmt.Sprintf("%v%v", filePath, ".html")
	}

	return filePath
}

// CreateFile creates the temporary file and opens it in the $EDITOR.
func CreateFile(filePath string, editorPath string) error {
	cmd := exec.Command(editorPath, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// ReadFile will open the file in question and return the string representation.
func ReadFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}
