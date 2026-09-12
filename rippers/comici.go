package rippers

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"NewHakuneko/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func getPagesComici(url string, cookie string) ([]utils.ComiciResult, *string, error) {
	// Make a request with the request function.
	res, err := utils.Request(url, cookie, false)
	// Err check.
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}

	// Parse HTML.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	// Err check.
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to parse HTML: %v", err)
	}
	// Close request
	defer res.Body.Close()

	// Extract viewer id and api domain
	viewerId, exists := doc.Find("#comici-viewer").Attr("data-comici-viewer-id")
	apiDomain, exists := doc.Find("#comici-viewer").Attr("data-api-domain")
	contentId, exists := doc.Find("#comici-viewer").Attr("data-content-id")
	// Extract title.
	title := doc.Find(".ep-main-h-h").Text()
	// Check they exit.
	if !exists {
		return nil, nil, fmt.Errorf("Does not exist")
	}
	// Check if the title is nothing, in that case, making the title the same thing as the viewerID.
	if title == "" {
		fmt.Println("Could not find a title. Using the viewerID as a placeholder title. Apologies.")
		title = viewerId
	}

	// Create a url by splicing it together.
	newUrl := "https://" + strings.Split(url, "/")[2] + apiDomain + "/book/contentsInfo?user-id=&comici-viewer-id=" + viewerId + "&page-from=0&page-to=1"
	if contentId != "" {
		newUrl = newUrl + "&contentId=" + contentId
	}

	// Make a second request.
	res, err = utils.Request(newUrl, cookie, true)
	// Err check.
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}

	// Read the bytes of the request.
	bodyBytes, err := io.ReadAll(res.Body)
	// Err check.
	if err != nil {
		return nil, nil, fmt.Errorf("Error while reading the request: %v", err)
	}
	// Declare the response variable.
	var response utils.RequestComici

	// Read the Json
	err = json.Unmarshal([]byte(bodyBytes), &response)
	// Err check.
	if err != nil {
		return nil, nil, fmt.Errorf("Error while decoding the josn.: %v", err)
	}
	// Close the request.
	defer res.Body.Close()

	// Get the total pages that we need to request.
	finalPage := response.TotalPages

	// Make second request.
	// (Before you make fun of this code, this is literally what officials do. There's no other way LMAO)
	// Recrate url with max pages pages.
	newUrl = "https://" + strings.Split(url, "/")[2] + apiDomain + "/book/contentsInfo?user-id=&comici-viewer-id=" + viewerId + "&page-from=0&page-to=" + strconv.Itoa(finalPage)
	if contentId != "" {
		newUrl = newUrl + "&contentId=" + contentId
	}

	// Make second request.
	res, err = utils.Request(newUrl, cookie, true)
	// Err check.
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}
	// Read the bytes of this second request.
	bodyBytes, err = io.ReadAll(res.Body)
	// Err check.
	if err != nil {
		return nil, nil, fmt.Errorf("Error while reading the request: %v", err)
	}
	// Read the json.
	err = json.Unmarshal([]byte(bodyBytes), &response)
	if err != nil {
		return nil, nil, fmt.Errorf("Error while decoding the josn.: %v", err)
	}

	// Close the request
	defer res.Body.Close()

	// If everything ok, then return the pages (no need to make a new Array for clarity like with Giga. Everthing's pretty clear already)
	return response.Result, &title, nil
}

// Get the coordinates to unscramble the image.
func getCoords(index, w, h int) (int, int, int, int) {
	// Set a grid size.
	gridSize := 4

	// Get width and height of the grids.
	tileW := w / gridSize
	tileH := h / gridSize

	// Get the x and y positions.
	x := (index / gridSize) * tileW
	y := (index % gridSize) * tileH

	// Return the x, y, and width and height.
	return x, y, tileW, tileH
}

// Download the pages comici.
func downloadPagesComici(pages []utils.ComiciResult, title string, ctx context.Context, folder string) error {
	// Create the directory with the title
	os.Mkdir(path.Join(folder, title), os.ModePerm)
	// Log to tell the user.
	runtime.EventsEmit(ctx, "title-get", "Downloading "+title)

	runtime.EventsEmit(ctx, "total-page", strconv.Itoa(len(pages)))

	// Loop over pages.
	for index, value := range pages {
		// Get the full path.
		fullPath := filepath.Join(folder, title, strconv.Itoa(index+1)+".png")

		// Make a request.
		res, err := utils.Request(value.ImageURL, "", true)
		// Err check
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		// Close the body.
		defer res.Body.Close()

		// Decode the image.
		srcImg, _, err := image.Decode(res.Body)
		// Err check.
		if err != nil {
			return fmt.Errorf("Error while decoding image: %w", err)
		}

		// Define bounds, height and width.
		bounds := srcImg.Bounds()
		imgW, imgH := bounds.Dx(), bounds.Dy()

		// Get the indicies of the scramble from the Json.
		var scrambleIndices []int
		if err := json.Unmarshal([]byte(value.Scramble), &scrambleIndices); err != nil {
			return fmt.Errorf("Error while decoding JSON: %w", err)
		}

		// Make RGBA layers.
		rgbaSrc := image.NewRGBA(bounds)

		// Draw basic canvas
		draw.Draw(rgbaSrc, bounds, srcImg, bounds.Min, draw.Src)
		canvas := image.NewRGBA(bounds)

		// Reconstruct grid structure array
		for destIdx, srcIdx := range scrambleIndices {
			sx, sy, tw, th := getCoords(srcIdx, imgW, imgH)
			dx, dy, _, _ := getCoords(destIdx, imgW, imgH)

			srcRect := image.Rect(sx, sy, sx+tw, sy+th)
			destRect := image.Rect(dx, dy, dx+tw, dy+th)

			draw.Draw(canvas, destRect, rgbaSrc, srcRect.Min, draw.Src)
		}

		// Create file
		outFile, err := os.Create(fullPath)
		// Erro check.
		if err != nil {
			return fmt.Errorf("Failed to create output file: %w", err)
		}
		// Close the file.
		defer outFile.Close()

		// Encode the file, and err check.
		if err := png.Encode(outFile, canvas); err != nil {
			return fmt.Errorf("Failed to encode output file: %w", err)
		}

		// Log that you're done with the page page.
		runtime.EventsEmit(ctx, "new-page", "")
	}

	// Final log
	runtime.EventsEmit(ctx, "title-get", "Done! Enjoy!")

	// Return nil when no error.
	return nil
}

// Main comici rip function
func RipComici(url string, cookie string, ctx context.Context, folder string) {
	// Get the pages, Title, and error in case.
	pages, title, error := getPagesComici(url, cookie)
	// Err check
	if error != nil {
		runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
	} else {
		// Downloading and de scrambling the pages.
		error = downloadPagesComici(pages, *title, ctx, folder)
		if error != nil {
			runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
		}
	}
}