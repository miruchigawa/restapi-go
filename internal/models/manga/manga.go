package manga

type MangaInfo struct {
    ID            string
    Title         string
    AltTitles     any
    Description   string
    Status        string
    ReleaseDate   int
    ContentRating string
    LastVolume    string
    LastChapter   string
    Image         string
}

type SearchResults struct {
	CurrentPage int           
	Results     []MangaInfo
}
