package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	runner "runtime"
	
	"NewHakuneko/cmd"
	
	"golang.org/x/net/html"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}


func Body(doc *html.Node) (*html.Node, error) {
    var body *html.Node
    var crawler func(*html.Node)
    crawler = func(node *html.Node) {
        if node.Type == html.ElementNode && node.Data == "body" {
            body = node
            return
        }
        for child := node.FirstChild; child != nil; child = child.NextSibling {
            crawler(child)
        }
    }
    crawler(doc)
    if body != nil {
        return body, nil
    }
    return nil, errors.New("Missing <body> in the node tree")
}

func (a *App) SelectTargetDirectory() (string, error) {
	directoryPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Download Folder",
	})

	if err != nil {
		return "", err
	}

	// Returns empty string if the user closed/cancelled the dialog
	return directoryPath, nil
}

// GetDefaultDocumentsFolder returns the absolute path to the OS Documents directory
func (a *App) GetDefaultDocumentsFolder() (string, error) {
	// 1. Get the current user's home directory path natively
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var documentsDir string

	// 2. Adjust the folder name based on the operating system
	if runner.GOOS == "windows" {
		// On Windows, the folder is usually named "Documents" inside the profile folder
		documentsDir = filepath.Join(homeDir, "Documents", "HebiStuff")
	} else if runner.GOOS == "darwin" {
		// On macOS, it's also named "Documents" inside the user home folder
		documentsDir = filepath.Join(homeDir, "Documents", "HebiStuff")
	} else {
		// On Linux/Unix, it usually follows XDG user dirs, defaulting to "Documents"
		documentsDir = filepath.Join(homeDir, "Documents", "HebiStuff")
	}

	return documentsDir, nil
}

func (a *App) Rip(url string, cookies string, folder string)  {
	cmd.RipMain(url, cookies, a.ctx, folder)
}
