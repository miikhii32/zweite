package rippers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"NewHakuneko/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func getPagesGanGan(url string, cookie string) ([]string, *string, error) {
	// Make the request, using an external utility.
	res, err := utils.Request(url, cookie, false)
	// Err check
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}

	// Parse the body
	doc, err := goquery.NewDocumentFromReader(res.Body)
	// Error check.
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to parse HTML: %v", err)
	}
	// Close the request
	defer res.Body.Close()

	// Extract the needed data
	gigaPagesData := doc.Find("#__NEXT_DATA__").Text()

	// Create the response variable with the correct type (imported from another file)
	var response utils.GanResponse
	// Handle the JSON
	err = json.Unmarshal([]byte(gigaPagesData), &response)
	// Error check
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing JSON: %v", err)
	}

	// Make a new variable with the pages, with the imported type.
	var Pages []string

	// Remake the array with the variable.
	for _, value := range response.Props.PageProps.Data.Pages {
		Pages = append(Pages, value.Image.ImageURL)
	}

	// Final return.
	return Pages, &response.Props.PageProps.Data.ChapterName, nil
}

// Function to download the pages with Giga.
func downloadPagesGanGan(pages []string, title string, ctx context.Context, folder string) error {
	// Make the directory of the downloaded file
	os.Mkdir(path.Join(folder, title), os.ModePerm)


	runtime.EventsEmit(ctx, "title-get", "Downloading "+title)

	// Loop over all the pages
	for index, value := range pages {
		// Get the full path.
		fullPath := filepath.Join(folder, title, strconv.Itoa(index+1)+".png")

		if (value == "") {
			continue
		}

		// Make the request.
		res, err := utils.Request("https://www.ganganonline.com" + value, "", false)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		out, err := os.Create(fullPath)
		if err != nil {
			return err
		}
		// Close the request.
		defer res.Body.Close()

		_, err = io.Copy(out, res.Body)
		if err != nil {
			return err
		}

		// Final log message.
		runtime.EventsEmit(ctx, "title-get", "Downloaded and De-Scrambled "+fullPath)
	}

	runtime.EventsEmit(ctx, "title-get", "Done! Enjoy!")

	// Return no error if loop completes normally.
	return nil
}

// Main ripping function for GigaReader.
func RipGanGan(url string, cookie string, ctx context.Context, folder string) {
	// Get the pages, Title, and error in case.
	pages, title, error := getPagesGanGan(url, cookie)
	// Err check
	if error != nil {
		runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
	} else {
		// Downloading and de scrambling the pages.
		error = downloadPagesGanGan(pages, *title, ctx, folder)
		if error != nil {
			runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
		}
	}
}