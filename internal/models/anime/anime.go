package anime

type SearchResult struct {
	CurrentPage int           `json:"CurrentPage"`
	HasNextPage bool          `json:"HasNextPage"`
	Results     []AnimeResult `json:"Results"`
}

type AnimeResult struct {
	ID          string `json:"Id"`
	Title       string `json:"Title"`
	URL         string `json:"Url"`
	Image       string `json:"Image"`
	ReleaseDate string `json:"ReleaseDate"`
	SubOrDub    SubOrDub `json:"SubOrDub"`
}

type MediaFormat string 

const (
	ANIME MediaFormat = "ANIME"
	MOVIE MediaFormat = "MOVIE"
)

type MediaStatus string 

const (
	ONGOING       MediaStatus = "ONGOING"
	COMPLETED     MediaStatus = "COMPLETED"
	NOT_YET_AIRED MediaStatus = "NOT_YET_AIRED"
	UNKNOWN       MediaStatus = "UNKNOWN"
)

type SubOrDub string 

const (
	SUB SubOrDub = "SUB"
	DUB SubOrDub = "DUB"
)

type Episode struct {
	ID     string
	Number float64
	URL    string
}

type AnimeInfo struct {
	ID 					string `json:"Id`
	Title 			string `json:"Title"`
	URL 				string `json:"Url"`
	Image 			string `json:"Image"`
	ReleaseDate string `json:"ReleaseDate"`
	Description string `json:"Description"`
	SubOrDub    SubOrDub 		`json:"SubOrDub`
	Type        MediaFormat `json:"Type"`
	Status      MediaStatus `json:"Status"`
	OtherName   string 			`json:"OtherName"`
	Genres      []string 		`json:"Genres"`
	TotalEpisodes int 			`json:"TotalEpisodes"`
	Episodes    []Episode 	`json:"Episodes"`
}

type EpisodeServer struct {
	Name string `json:"Name"`
	URL  string `json:"URL"`
}
