package rippers

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"NewHakuneko/utils"

	pb "NewHakuneko/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"google.golang.org/protobuf/proto"
)

const protoSchema = `
  syntax = "proto3";

  message DeviceInfo {
    string secret = 1;
    string appVer = 2;
    int32 deviceType = 3;
    string osVer = 4;
    bool isTablet = 5;
    int32 imageQuality = 6;
  }

  message UserPoint {
    uint32 free = 1;
    uint32 paid = 2;
  }

  message WebMangaViewerRequest {
    DeviceInfo deviceInfo = 1;
    bool useTicket = 2;
    UserPoint consumePoint = 3;
    uint32 chapterId = 4;
    ChapterArgument chapterArgument = 5;

    message ChapterArgument {
      uint32 mangaId = 1;
      int32 position = 2;
    }
  }

  message Image {
    string imageUrl = 1;
    string urlScheme = 2;
    string iv = 3;
    string encryptionKey = 4;
    uint32 imageWidth = 5;
    uint32 imageHeight = 6;
    bool isExtraPage = 7;
    uint32 extraId = 8;
    uint32 extraIndex = 9;
    uint32 extraSlotId = 10;
  }

  message WebView {
    string url = 1;
  }

  message LastPage {}
  message GoogleAds {}

  message ViewerPage {
    Image image = 1;
    WebView webview = 2;
    LastPage lastPage = 3;
    GoogleAds googleAds = 4;
  }

  message BookViewerResponse {
    string viewerTitle = 1;
    repeated ViewerPage pages = 2;
  }

  message WebMangaViewerResponse {
    BookViewerResponse data = 2;
  }
`;


func getPagesFuz(url string, cookie string) ([]*utils.ViewerPage, *string, error) {
	// Make the request, using an external utility.
	res, err := utils.Request(url, cookie, true)
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
	comicFuzPageData := doc.Find("#__NEXT_DATA__").Text()
	fmt.Println(comicFuzPageData)

	// Create the response variable with the correct type (imported from another file)
	var response utils.NextData
	// Handle the JSON
	err = json.Unmarshal([]byte(comicFuzPageData), &response)
	// Error check
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing JSON: %v", err)
	}

	fmt.Println(response)

	/*
	// Check if chapter is available.
	if response.ReadableProduct.HasPurchased == false && response.ReadableProduct.IsPublic == false {
		return nil, nil, fmt.Errorf("It seems like you don't have access to this content. Maybe try adding a cookie string, or if you are adding one, check to see if you've purchased the desired content.")
	}
	*/

	cID, err := strconv.ParseUint(response.Props.PageProps.ChapterId, 10, 32)
	if (err != nil){
		return nil, nil, err
	}

	requestPayload := &pb.WebMangaViewerRequest{
		DeviceInfo: &pb.DeviceInfo{
			DeviceType: 2,
		},
		ChapterId: uint32(cID),
		UseTicket: false,
	}

	requestBuffer, err := proto.Marshal(requestPayload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	res, err = utils.RequestFuz("https://api.comic-fuz.com/v1/web_manga_viewer", cookie, true, requestBuffer)
	// Err check
	if err != nil {
		return nil, nil, fmt.Errorf("%v", err)
	}
	// Close the request
	defer res.Body.Close()

	var reader io.Reader = res.Body

	if res.Header.Get("Content-Encoding") == "gzip" {
    	gzReader, err := gzip.NewReader(res.Body)
    	if err != nil {
     	   fmt.Println("Failed to create gzip reader: %v", err)
    	}
    	defer gzReader.Close()
    	reader = gzReader
	}

	responseBuffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}

	responsePayload := &pb.WebMangaViewerResponse{}

	err = proto.Unmarshal(responseBuffer, responsePayload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	data := responsePayload.GetData()
	if data == nil {
		return nil, nil, fmt.Errorf("no data found in response")
	}

	pages := data.GetPages()
	finalTitle := response.Props.PageProps.Data.ChapterMainName + response.Props.PageProps.Data.ChapterSubName

	return pages, &finalTitle, nil
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
func downloadPagesFuz(pages []*utils.ViewerPage, title string, ctx context.Context, folder string) error {
	fmt.Println(pages)
	// Make the directory of the downloaded file
	os.Mkdir(path.Join(folder, title), os.ModePerm)

	runtime.EventsEmit(ctx, "title-get", "Downloading "+title)

	// Loop over all the pages
	for index, value := range pages {
		if value.GoogleAds != nil { 
			break
		}
		if (!value.Image.IsExtraPage){
				// Get the full path.
				fullPath := filepath.Join(folder, title, strconv.Itoa(index+1)+".png")
		
				// Make the request.
				res, err := utils.Request("https://img.comic-fuz.com" + value.Image.ImageUrl, "", false)
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
	pages, title, error := getPagesFuz(url, cookie)
	// Err check
	if error != nil {
		runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
	} else {
		// Downloading and de scrambling the pages.
		error = downloadPagesFuz(pages, *title, ctx, folder)
		if error != nil {
			runtime.EventsEmit(ctx, "error-emit", fmt.Sprintf("%v", error))
		}
	}
}