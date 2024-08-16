package manga

import (
    "fmt"
    "net/http"
    "net/url"
    "errors"
    "encoding/json"

	models "miruchigawa.moe/restapi/internal/models/manga"
)

var (
    baseURL string = "https://mangadex.org"
    apiURL string = "https://api.mangadex.org"
)

type MangadexSearchResponse struct {
	Result string `json:"result"`
	Data   []struct {
		ID         string `json:"id"`
		Attributes struct {
			Title         map[string]string `json:"title"`
			AltTitles     interface{}       `json:"altTitles"`
			Description   map[string]string `json:"description"`
			Status        string            `json:"status"`
			Year          int               `json:"year"`
			ContentRating string            `json:"contentRating"`
			LastVolume    string            `json:"lastVolume"`
			LastChapter   string            `json:"lastChapter"`
		} `json:"attributes"`
		Relationships []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"relationships"`
	} `json:"data"`
}

func Search(query string, page, limit int) (*models.SearchResults, error) {
    if page <= 0 {
        return nil, errors.New("page number must be greater than 0")
    }

    if limit > 100 {
        return nil, errors.New("limit must be less than or equal to 100")
    }

    if limit*(page-1) >= 10000 {
        return nil, errors.New("not enough results")
    }

    params := url.Values{}
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("title", query)
	params.Set("offset", fmt.Sprintf("%d", limit*(page-1)))
	params.Set("order[relevance]", "desc")

	resp, err := http.Get(fmt.Sprintf("%s/manga?%s", apiURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response MangadexSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if response.Result != "ok" {
		return nil, errors.New("failed to fetch manga results")
	}

	results := &models.SearchResults{
	    CurrentPage: page,
	    Results: []models.MangaInfo{},
	}

	for _, manga := range response.Data {
        var coverArt string
		for _, rel := range manga.Relationships {
			if rel.Type == "cover_art" {
				coverArt, err = fetchCoverImage(rel.ID)
				if err != nil {
					return nil, err
				}
				break
			}
		}

	    results.Results = append(results.Results, models.MangaInfo{
	        ID: manga.ID,
	        Title: manga.Attributes.Title["en"],
	        AltTitles: manga.Attributes.AltTitles,
	        Description: manga.Attributes.Description["en"],
	        Status: manga.Attributes.Status,
	        ReleaseDate: manga.Attributes.Year,
	        ContentRating: manga.Attributes.ContentRating,
			LastVolume:    manga.Attributes.LastVolume,
			LastChapter:   manga.Attributes.LastChapter,
			Image: fmt.Sprintf("%s/covers/%s/%s", baseURL, manga.ID, coverArt),
	    })
	}

	return results, nil
}


type CoverResponse struct {
	Data struct {
		Attributes struct {
			FileName string `json:"fileName"`
		} `json:"attributes"`
	} `json:"data"`
}

func fetchCoverImage(id string) (string, error) {
    url := fmt.Sprintf("%s/cover/%s", apiURL, id)

    resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch cover image, status code: %d", resp.StatusCode)
	}

	var response CoverResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Data.Attributes.FileName, nil
}


