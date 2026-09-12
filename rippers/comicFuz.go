package rippers

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"NewHakuneko/utils"

	pbBook "NewHakuneko/protoBookRequests"
	pb "NewHakuneko/protoRequests"

	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"google.golang.org/protobuf/proto"
)

func getPagesFuz(url string, cookie string) ([]*pb.ViewerPage, []*pbBook.ViewerPage, *string, error) {

	var cID uint64
	var response utils.NextData
	var err error

	if strings.Split(url, "/")[len(strings.Split(url, "/"))-3] == "book" {
		// Make the request, using an external utility.		
		cID, err = strconv.ParseUint(strings.Split(url, "/")[len(strings.Split(url, "/"))-1], 10, 32)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("%v", err)
		}
	} else {
		// Make the request, using an external utility.
		res, err := utils.Request(url, cookie, true)
		// Err check
		if err != nil {
			return nil, nil, nil, fmt.Errorf("%v", err)
		}

		// Parse the body
		doc, err := goquery.NewDocumentFromReader(res.Body)
		// Error check.
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Failed to parse HTML: %v", err)
		}
		// Close the request
		defer res.Body.Close()

		// Extract the needed data
		comicFuzPageData := doc.Find("#__NEXT_DATA__").Text()

		// Create the response variable with the correct type (imported from another file)
		var response utils.NextData
		// Handle the JSON
		err = json.Unmarshal([]byte(comicFuzPageData), &response)
		// Error check
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Error parsing JSON: %v", err)
		}

		cID, err = strconv.ParseUint(response.Props.PageProps.ChapterId, 10, 32)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("%v", err)
		}
	}

	var newUrl string
	isBook := false

	if strings.Split(url, "/")[len(strings.Split(url, "/"))-3] == "book" {
		newUrl = "https://api.comic-fuz.com/v1/book_viewer_2"
		isBook = true
	} else {
		newUrl = "https://api.comic-fuz.com/v1/web_manga_viewer"
	}

	if isBook {
		requestPayload := &pbBook.BookWebMangaViewerRequest{
			DeviceInfo: &pbBook.BookDeviceInfo{
				DeviceType: 2,
			},
			ChapterId: uint32(cID),
			UseTicket: false,
		}
		requestBuffer, err := proto.Marshal(requestPayload)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to marshal request: %v", err)
		}
		res, err := utils.RequestFuz(newUrl, cookie, true, requestBuffer)
		// Err check
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Error in request: %v", err)
		}
		// Close the request
		defer res.Body.Close()

		var reader io.Reader = res.Body

		if res.Header.Get("Content-Encoding") == "gzip" {
			gzReader, err := gzip.NewReader(res.Body)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("Failed to create gzip reader: %v", err)
			}
			defer gzReader.Close()
			reader = gzReader
		}

		responseBuffer, err := io.ReadAll(reader)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("%v", err)
		}

		responsePayload := &pbBook.WebMangaViewerResponse{}

		err = proto.Unmarshal(responseBuffer, responsePayload)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		data := responsePayload.Pages
		if data == nil {
			return nil, nil, nil, fmt.Errorf("no data found in response")
		}
	
		finalTitle := responsePayload.Info.Title

		return nil, responsePayload.Pages, &finalTitle, nil
	} else {
		requestPayload := &pb.WebMangaViewerRequest{
			DeviceInfo: &pb.DeviceInfo{
				DeviceType: 2,
			},
			ChapterId: uint32(cID),
			UseTicket: false,
		}
		requestBuffer, err := proto.Marshal(requestPayload)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to marshal request: %v", err)
		}
		res, err := utils.RequestFuz(newUrl, cookie, true, requestBuffer)
		// Err check
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Error in request: %v", err)
		}
		// Close the request
		defer res.Body.Close()

		var reader io.Reader = res.Body

		if res.Header.Get("Content-Encoding") == "gzip" {
			gzReader, err := gzip.NewReader(res.Body)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("Failed to create gzip reader: %v", err)
			}
			defer gzReader.Close()
			reader = gzReader
		}

		responseBuffer, err := io.ReadAll(reader)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("%v", err)
		}

		responsePayload := &pb.WebMangaViewerResponse{}

		err = proto.Unmarshal(responseBuffer, responsePayload)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to unmarshal response: %v", err)
		}

		data := responsePayload.GetData()
		if data == nil {
			return nil, nil, nil, fmt.Errorf("no data found in response")
		}

		pages := data.GetPages()
		finalTitle := response.Props.PageProps.Data.ChapterMainName + response.Props.PageProps.Data.ChapterSubName

		return pages, nil, &finalTitle, nil
	}
}

// DecryptAESCBC decrypts AES-CBC encrypted bytes using hex-encoded key and IV strings
func DecryptAESCBC(ciphertext []byte, hexKey, hexIV string) ([]byte, error) {
	// 1. Decode Key and IV from hex strings to raw byte slices
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("invalid hex key: %w", err)
	}

	iv, err := hex.DecodeString(hexIV)
	if err != nil {
		return nil, fmt.Errorf("invalid hex IV: %w", err)
	}

	// 2. Initialize AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher block: %w", err)
	}

	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("ciphertext length (%d) is not a multiple of block size (%d)", len(ciphertext), block.BlockSize())
	}

	// 3. Decrypt blocks using CBC mode
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// 4. Strip PKCS#7 padding if present
	plaintext = pkcs7Unpad(plaintext)

	return plaintext, nil
}

// Helper to remove standard PKCS#7 padding
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	padding := int(data[length-1])
	if padding > length || padding > aes.BlockSize {
		return data // Return unpadded if invalid padding
	}
	return data[:length-padding]
}

// Function to download the pages with Giga.
func downloadPagesFuz(pages []*pb.ViewerPage, pagesBook []*pbBook.ViewerPage, title string, ctx context.Context, folder string) error {
	// Make the directory of the downloaded file
	os.Mkdir(path.Join(folder, title), os.ModePerm)

	runtime.EventsEmit(ctx, "title-get", "Downloading "+ title)
	runtime.EventsEmit(ctx, "total-page", strconv.Itoa(len(pages)))

	if pages != nil {

		// Loop over all the pages
		for index, value := range pages {
			if value.GoogleAds != nil {
				break
			}
			if !value.Image.IsExtraPage {
				// Get the full path.
				fullPath := filepath.Join(folder, title, strconv.Itoa(index+1)+".png")

				// Make the request.
				res, err := utils.Request("https://img.comic-fuz.com"+value.Image.ImageUrl, "", false)
				// Err check.
				if err != nil {
					return fmt.Errorf("%v", err)
				}
				// Close the request.
				defer res.Body.Close()

				// 2. Read the entire encrypted response body into memory
				encryptedBytes, err := io.ReadAll(res.Body)
				if err != nil {
					return fmt.Errorf("failed to read response body: %w", err)
				}

				// 3. Decrypt using the AES-CBC function
				decryptedBytes, err := DecryptAESCBC(encryptedBytes, value.Image.EncryptionKey, value.Image.Iv)
				if err != nil {
					return fmt.Errorf("failed to decrypt image: %w", err)
				}

				// 4. Save the decrypted output file
				err = os.WriteFile(fullPath, decryptedBytes, 0644)
				if err != nil {
					return fmt.Errorf("failed to write decrypted file: %w", err)
				}

				// Final log message.
				runtime.EventsEmit(ctx, "new-page", "")
			}
		}

		runtime.EventsEmit(ctx, "title-get", "Done! Enjoy!")

		// Return no error if loop completes normally.
		return nil

	} else {
		// Loop over all the pages
		for index, value := range pagesBook {
			// Get the full path.
			fullPath := filepath.Join(folder, title, strconv.Itoa(index+1)+".png")

			// Make the request.
			res, err := utils.Request("https://img.comic-fuz.com"+value.Image.ImageUrl, "", false)
			// Err check.
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			// Close the request.
			defer res.Body.Close()

			// 2. Read the entire encrypted response body into memory
			encryptedBytes, err := io.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}

			// 3. Decrypt using the AES-CBC function
			decryptedBytes, err := DecryptAESCBC(encryptedBytes, value.Image.EncryptionKey, value.Image.Iv)
			if err != nil {
				return fmt.Errorf("failed to decrypt image: %w", err)
			}

			// 4. Save the decrypted output file
			err = os.WriteFile(fullPath, decryptedBytes, 0644)
			if err != nil {
				return fmt.Errorf("failed to write decrypted file: %w", err)
			}

			// Final log message.
			runtime.EventsEmit(ctx, "title-get", "Downloaded and De-Scrambled "+fullPath)
		}
	}

	runtime.EventsEmit(ctx, "title-get", "Done! Enjoy!")

	// Return no error if loop completes normally.
	return nil
}

// Main ripping function for GigaReader.
func RipFuz(url string, cookie string, ctx context.Context, folder string) {
	// Get the pages, Title, and error in case.
	pages, pagesBook, title, error := getPagesFuz(url, cookie)
	// Err check
	if error != nil {
		runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
	} else {
		// Downloading and de scrambling the pages.
		error = downloadPagesFuz(pages, pagesBook, *title, ctx, folder)
		if error != nil {
			runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
		}
	}
}
