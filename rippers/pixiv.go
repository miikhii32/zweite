package rippers

import (
	"io"
	"context"
	"encoding/json"
	"fmt"
	_ "image/png"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"NewHakuneko/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func getPagesPixiv(url string, cookie string) ([]utils.GigaPages, *string, error) {
	splitUrl := strings.Split(url, "/")

	lastElement := ""

	if len(splitUrl) > 0 {
		// 2. Access the last element
		lastElement = splitUrl[len(splitUrl)-1]
	} else {
		return nil, nil, fmt.Errorf("%v", "Please include a link!")
	}

	if strings.Contains(lastElement, "#") {
		lastElement = strings.Split(lastElement, "#")[0]
	}

	res, err := utils.Request("https://www.pixiv.net/ajax/illust/" + lastElement + "/pages?lang=en", "", true)

	if (err != nil){
		return nil, nil, fmt.Errorf("Failed to fetch pages. Here's the error: %v", err)	
	}
	defer res.Body.Close()


	// Create the response variable with the correct type (imported from another file)
	var response utils.PixivResponse
	// Handle the JSON
	err = json.NewDecoder(res.Body).Decode(&response)
	// Error check
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing JSON: %v", err)
	}

	// Make a new variable with the pages, with the imported type.
	var Pages []utils.GigaPages

	// Remake the array with the variable.
	for _, value := range response.Body {
		if value.Height != 0 {
			Pages = append(Pages, utils.GigaPages{PageUrl: value.Urls.Original, Height: value.Height, Width: value.Width})
		}
	}

	res, err = utils.Request("https://www.pixiv.net/ajax/illust/" + lastElement, "", true)

	if (err != nil){
		return nil, nil, fmt.Errorf("%v", err)	
	}
	defer res.Body.Close()

	// Create the response variable with the correct type (imported from another file)
	var titleRes utils.PixivResponseTitle
	// Handle the JSON
	err = json.NewDecoder(res.Body).Decode(&titleRes)
	// Error check
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing JSON: %v", err)
	}

	// Final return.
	return Pages, &titleRes.Body.IllustTitle, nil
}

// Function to download the pages with Giga.
func downloadPagesPixiv(pages []utils.GigaPages, title string, ctx context.Context, folder string) error {
	// Make the directory of the downloaded file
	os.Mkdir(path.Join(folder, title), os.ModePerm)

	runtime.EventsEmit(ctx, "title-get", "Downloading "+title)

	// Loop over all the pages
	for index, value := range pages {
		// Get the full path.
		fullPath := filepath.Join(folder, title, strconv.Itoa(index+1)+".png")

		// Make the request.
		res, err := utils.Request(value.PageUrl, "", true)
		// Err check.
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
func RipPixiv(url string, cookie string, ctx context.Context, folder string) {
	// Get the pages, Title, and error in case.
	pages, title, error := getPagesPixiv(url, cookie)
	// Err check
	if error != nil {
		runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
	} else {
		// Downloading and de scrambling the pages.
		error = downloadPagesPixiv(pages, *title, ctx, folder)
		if error != nil {
			runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
		}
	}
}