package utils

type PageStructureArray struct {
	ImgType      string `json:"type"`
	Width        int    `json:"width"`
	Src          string `json:"src"`
	ContentStart string `json:"contentStart"`
	Height       int    `json:"height"`
}

type PageStructure struct {
	Pages []PageStructureArray `json:"pages"`
}

type ReadBleProduct struct {
	ShowSquareThumbnailInRecommendation bool          `json:"showSquareThumbnailInRecommendation"`
	NextReadableProductUri              string        `json:"nextReadableProductUri"`
	ImageUrisDigest                     string        `json:"imageUrisDigest"`
	PublishedAt                         string        `json:"publishedAt"`
	Title                               string        `json:"title"`
	Toc                                 *string       `json:"toc"`
	PageStructure                       PageStructure `json:"pageStructure"`
	HasPurchased                        bool          `json:"hasPurchased"`
	Permalink                           string        `json:"permalink"`
	IsPublic                            bool          `json:"isPublic"`
	FinishReadingNotificationUri        *string       `json:"finishReadingNotificationUri"`
	PrevReadableProductUri              *string       `json:"prevReadableProductUri"`
	ID                                  string        `json:"id"`
}

type GigaResponse struct {
	ReadableProduct ReadBleProduct `json:"readableProduct"`
}

type GigaPages struct {
	PageUrl string
	Width   int
	Height  int
}


type ComiciResult struct {
	ImageURL string `json:"imageURL"` 
	Scramble string `json:"scramble"` 
	Sort int `json:"sort"` 
	Width int `json:"width"` 
	Height int `json:"height"` 
}

type RequestComici struct {
	TotalPages int `json:"totalPages"`
	Result []ComiciResult `json:"result"`
}