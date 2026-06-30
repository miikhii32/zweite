/*
Copyright © 2026 mikhi32
*/
package cmd

// Imports
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
	"slices"
	"strconv"
	"strings"

	"NewHakuneko/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Don't ask why I made a function here... arrays in go are way too complicated lol.
func getGigaReaders() []string {
	gigaReaders := [...]string{
		"comic-days.com",
		"ichicomi.com",
		"shonenjumpplus.com",
		"tonarinoyj.jp",
		"andsofa.com",
		"morningtwo.com",
		"getsumagakichi.com",
		"bibliosirius.com",
		"kuragebunch.com",
		"comicbunch-kai.com",
		"viewer.heros-web.com",
		"comicborder.com",
		"comic-gardo.com",
		"comic-zenon.com",
		"magcomi.com",
		"comic-action.com",
		"comic-trail.com",
		"feelweb.jp",
		"www.sunday-webry.com",
		"comic-ogyaaa.com",
		"comic-earthstar.com",
		"ourfeel.jp",
		"comic-seasons.com",
		"comic-y-ours.com",
	}
	return gigaReaders[:]
}

func getComiciReaders() []string {
	comiciReaders := [...]string{
		"rimacomiplus.jp",
		"younganimal.com",
		"takecomic.jp",
		"kimicomi.com",
		"youngchampion.jp",
		"mangalt.jp",
		"mangaspa.nikkan-spa.jp",
		"piacomic.jp",
		"comic-room-base.com",
		"manga-zegra.com",
		"comic-growl.com",
		"asacomi.jp",
		"comicpash.jp",
		"comic.j-nbooks.jp",
		"hayacomic.jp",
		"championcross.jp",
		"kansai.mag-garden.co.jp",
		"comicride.jp",
		"carula.jp",
		"comic-medu.com",
		"comics.manga-bang.com",
		"bigcomics.jp",
		"studio.booklista.co.jp",
		"hanayume.com",
	}
	return comiciReaders[:]
}

// Function that checks the url it's getting and determines the reader.
func checkString(url string) string {
	gigaReaders := getGigaReaders()
	comiciReaders := getComiciReaders()
	isGiga := slices.Contains(gigaReaders, url)
	isComici := slices.Contains(comiciReaders, url)
	if isGiga {
		return "giga"
	} else if isComici {
		return "comici"
	} else {
		return ""
	}
}

// Get all the pages for Giga Reader.
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

	runtime.EventsEmit(ctx, "title-get", "Downloading "+ title)

	// Loop over all the pages
	for index, value := range pages {
		// Get the full path.
		fullPath := filepath.Join(folder, title, strconv.Itoa(index)+".png")

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
		runtime.EventsEmit(ctx, "title-get", "Downloaded and De-Scrambled " + fullPath)
	}
	
	runtime.EventsEmit(ctx, "title-get", "Done! Enjoy!")

	// Return no error if loop completes normally.
	return nil
}

// Main ripping function for GigaReader.
func ripGiga(url string, cookie string, ctx context.Context, folder string) {
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

// Function to get pages for Comici
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
	runtime.EventsEmit(ctx, "title-get", "Downloading " + title)

	// Loop over pages.
	for index, value := range pages {
		// Get the full path.
		fullPath := filepath.Join(folder, title, strconv.Itoa(index)+".png")

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
		runtime.EventsEmit(ctx, "title-get", "Downloaded and De-Scrambled " + fullPath)
	}
	
	// Final log
	runtime.EventsEmit(ctx, "title-get", "Done! Enjoy!")

	// Return nil when no error.
	return nil
}

// Main comici rip function
func ripComici(url string, cookie string, ctx context.Context, folder string) {
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


// Main function to organize everything
func RipMain(url string, cookie string, ctx context.Context, folder string) {
	if url == ""{
		runtime.EventsEmit(ctx, "error-emit", "You need to send something for us to rip!")
	} else {
		// Check what reader we are using.
		if !strings.HasPrefix(url, "https://"){
			runtime.EventsEmit(ctx, "error-emit", "Please submit a url!")
		} else {
			reader := checkString(strings.Split(url, "/")[2])
			// Rip accordingly.
			if reader == "giga" {
				ripGiga(url, cookie, ctx, folder)
			} else if reader == "comici" {
				ripComici(url, cookie, ctx, folder)
			} else {
				runtime.EventsEmit(ctx, "error-emit", "We currently do not support this url. Apologies!")
			}
		}
	}
}