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

type PixivResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"Message"`
	Body    []PixivBody `json:"Body"`
}

type PixivBody struct {
	Urls   PixivUrls `json:"Urls"`
	Width  int `json:"Width"`
	Height int `json:"Height"`
}

type PixivUrls struct {
	ThumbMini string `json:"ThumbMini"`
	Small     string `json:"Small"`
	Regular   string `json:"Regular"`
	Original  string `json:"Original"`
}

type PixivResponseTitle struct {
	Error   bool        `json:"error"`
	Message string      `json:"Message"`
	Body    PixivBodyTitle `json:"Body"`
}

type PixivBodyTitle struct {
	IllustTitle string `json:"IllustTitle"` 
}


type NextData struct {
	Props PropsData `json:"props"`
}

type PropsData struct {
	PageProps PagePropsData `json:"pageProps"`
}

type PagePropsData struct {
	ChapterId string `json:"chapterId"`
	Data PagePropsDataData `json:"data"`
}

type PagePropsDataData struct {
	ChapterMainName string `json:"chapterMainName"`
	ChapterSubName string `json:"chapterSubName"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	MangaId int `json:"mangaId"`
	MangaName string `json:"mangaName"`
}


type GanImage struct {
	ImageURL string `json:"ImageURL"`
}

type GanPageObject struct {
	Image GanImage `json:"Image"` 
}

type GanData struct {
	ChapterName string `json:"chapterName"`
	Pages []GanPageObject `json:"pages"`
}

type GanPageProps struct {
	Data GanData `json:"Data"`
}

type GanProps struct {
	PageProps GanPageProps `json:"PageProps"`
}

type GanResponse struct {
	Props GanProps `json:"props"`
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
	Sort     int    `json:"sort"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type RequestComici struct {
	TotalPages int            `json:"totalPages"`
	Result     []ComiciResult `json:"result"`
}
