package anime

type SearchResult struct {
	CurrentPage int           
	HasNextPage bool          
	Results     []AnimeResult 
}

type AnimeResult struct {
	ID          string 
	Title       string 
	URL         string 
	Image       string 
	ReleaseDate string 
	SubOrDub    SubOrDub 
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
	ID 					string 
	Title 			string 
	URL 				string 
	Image 			string 
	ReleaseDate string 
	Description string 
	SubOrDub    SubOrDub 		
	Type        MediaFormat 
	Status      MediaStatus 
	OtherName   string 			
	Genres      []string 		
	TotalEpisodes int 			
	Episodes    []Episode 	
}

type EpisodeServer struct {
	Name string 
	URL  string 
}
