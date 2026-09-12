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
	"os"
	"path"
	"path/filepath"
	"strconv"

	"NewHakuneko/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func getPagesGiga(url string, cookie string) ([]utils.GigaPages, *string, error) {
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
	gigaPagesData, exists := doc.Find("#episode-json").Attr("data-value")
	// Check if it exists.
	if !exists {
		return nil, nil, fmt.Errorf("Element does not exist, %v", err)
	}

	// Create the response variable with the correct type (imported from another file)
	var response utils.GigaResponse
	// Handle the JSON
	err = json.Unmarshal([]byte(gigaPagesData), &response)
	// Error check
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing JSON: %v", err)
	}

	// Check if chapter is available.
	if response.ReadableProduct.HasPurchased == false && response.ReadableProduct.IsPublic == false {
		return nil, nil, fmt.Errorf("It seems like you don't have access to this content. Maybe try adding a cookie string, or if you are adding one, check to see if you've purchased the desired content.")
	}

	// Make a new variable with the pages, with the imported type.
	var Pages []utils.GigaPages

	// Remake the array with the variable.
	for _, value := range response.ReadableProduct.PageStructure.Pages {
		if value.Height != 0 {
			Pages = append(Pages, utils.GigaPages{PageUrl: value.Src, Height: value.Height, Width: value.Width})
		}
	}

	// Final return.
	return Pages, &response.ReadableProduct.Title, nil
}

// Function to download the pages with Giga.
func downloadPagesGiga(pages []utils.GigaPages, title string, ctx context.Context, folder string) error {
	// Make the directory of the downloaded file
	os.Mkdir(path.Join(folder, title), os.ModePerm)

	runtime.EventsEmit(ctx, "title-get", "Downloading "+title)

	runtime.EventsEmit(ctx, "total-page", strconv.Itoa(len(pages)))

	// Loop over all the pages
	for index, value := range pages {
		// Get the full path.
		fullPath := filepath.Join(folder, title, strconv.Itoa(index+1)+".png")

		// Make the request.
		res, err := utils.Request(value.PageUrl, "", false)
		// Err check.
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		// Close the request.
		defer res.Body.Close()

		// Decode the body of the image.
		srcImg, _, err := image.Decode(res.Body)
		// Err check.
		if err != nil {
			return fmt.Errorf("failed to unscramble image: %w", err)
		}

		// Set the bounds of the image
		bounds := srcImg.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Set the dividing and multiplying number.
		div := 4
		mul := 8

		// Convert image to RGBA format.
		rgbaSrc := image.NewRGBA(bounds)

		// Draw initial image.
		draw.Draw(rgbaSrc, bounds, srcImg, bounds.Min, draw.Src)

		// Create a clean destination canvas
		rgbaDest := image.NewRGBA(bounds)
		// Fill the background first
		draw.Draw(rgbaDest, bounds, rgbaSrc, bounds.Min, draw.Src)
		// Calculate total tiles.
		totalTiles := div * mul
		// Calculate the heights and widths of the tiles.
		fixedWidth := (width / totalTiles) * mul
		fixedHeight := (height / totalTiles) * mul

		// Map out where each tile goes based on the inverted matrix translation
		for y := 0; y < div; y++ {
			for x := 0; x < div; x++ {
				sourceX := y * fixedWidth
				sourceY := x * fixedHeight

				destX := x * fixedWidth
				destY := y * fixedHeight

				srcRect := image.Rect(sourceX, sourceY, sourceX+fixedWidth, sourceY+fixedHeight)
				destRect := image.Rect(destX, destY, destX+fixedWidth, destY+fixedHeight)

				draw.Draw(rgbaDest, destRect, rgbaSrc, srcRect.Min, draw.Src)
			}
		}

		// Create the final file.
		outFile, err := os.Create(fullPath)
		// Err check.
		if err != nil {
			return fmt.Errorf("Failed to create output file: %w", err)
		}
		// Close the final file.
		defer outFile.Close()

		// Encode the final file.
		if err := png.Encode(outFile, rgbaDest); err != nil {
			return fmt.Errorf("Failed to encode final file: %w", err)
		}

		// Final log message.
		runtime.EventsEmit(ctx, "new-page", "")
	}

	runtime.EventsEmit(ctx, "title-get", "Done! Enjoy!")

	// Return no error if loop completes normally.
	return nil
}

// Main ripping function for GigaReader.
func RipGiga(url string, cookie string, ctx context.Context, folder string) {
	// Get the pages, Title, and error in case.
	pages, title, error := getPagesGiga(url, cookie)
	// Err check
	if error != nil {
		runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
	} else {
		// Downloading and de scrambling the pages.
		error = downloadPagesGiga(pages, *title, ctx, folder)
		if error != nil {
			runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
		}
	}
}