package anime

import (
	"fmt"
	"strings"
	"strconv"
	"net/url"

	"github.com/gocolly/colly/v2"

	models "miruchigawa.moe/restapi/internal/models/anime"
)

var (
	baseURL string = "https://anitaku.pe"
	ajaxURL string = "https://ajax.gogocdn.net/ajax"
)

func Search(query string, page int) (*models.SearchResult, error) {
	searchResult := &models.SearchResult{
		CurrentPage: page,
		HasNextPage: false,
		Results:     []models.AnimeResult{},
	}

	c := colly.NewCollector()

	c.OnHTML("div.anime_name.new_series > div > div > ul > li.selected", func(e *colly.HTMLElement) {
		if e.DOM.Next().Length() > 0 {
			searchResult.HasNextPage = true
		}
	})	

	c.OnHTML("div.last_episodes > ul > li", func(e *colly.HTMLElement) {
		result := models.AnimeResult{
			ID:          strings.Split(e.ChildAttr("p.name > a", "href"), "/")[2],
			Title:       e.ChildText("p.name > a"),
			URL:         baseURL + e.ChildAttr("p.name > a", "href"),
			Image:       e.ChildAttr("div > a > img", "src"),
			ReleaseDate: strings.Split(e.ChildText("p.released"), "Released: ")[1],
			SubOrDub:    determineSubOrDub(e.ChildText("p.name > a")),
		}
		searchResult.Results = append(searchResult.Results, result)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\n", err)
	})

	url := fmt.Sprintf("%s/filter.html?keyword=%s&page=%d", baseURL, query, page)
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

func Info(id string) (*models.AnimeInfo, error) {
	result := &models.AnimeInfo{ Episodes: []models.Episode{} }

	if !strings.Contains(id, "gogoanime") {
		id = fmt.Sprintf("%s/category/%s", baseURL, id)
	}

	c := colly.NewCollector()

	c.OnHTML("body", func (e *colly.HTMLElement) {
		u, err := url.Parse(id)
		if err != nil {
			return
		}

		parts := strings.Split(u.Path, "/")
		if len(parts) > 2 {
			result.ID = parts[2]
		}

		result.Title = e.ChildText("section.content_left > div.main_body > div:nth-child(2) > div.anime_info_body_bg > h1")
		result.URL = id
		result.Image = e.ChildAttr("div.anime_info_body_bg > img", "src")
		result.ReleaseDate = strings.TrimSpace(strings.Split(e.ChildText("div.anime_info_body_bg > p:nth-child(8)"), "Released: ")[1])
		result.Description = strings.TrimPrefix(e.ChildText("div.anime_info_body_bg > div:nth-child(6)"), "Plot Summary: ")
		result.SubOrDub = determineSubOrDub(result.Title)
		result.Type = models.MediaFormat(strings.Split(strings.ToUpper(e.ChildText("div.anime_info_body_bg > p:nth-child(4) > a")), " ")[2])
		
		status := e.ChildText("div.anime_info_body_bg > p:nth-child(9) > a")
		switch status {
		case "Ongoing":
			result.Status = models.ONGOING
		case "Completed":
			result.Status = models.COMPLETED
		case "Upcoming":
			result.Status = models.NOT_YET_AIRED
		default:
			result.Status = models.UNKNOWN
		}

		result.OtherName = e.ChildText(".other-name a")
		e.ForEach("div.anime_info_body_bg > p:nth-child(7) > a", func(_ int, el *colly.HTMLElement) {
			result.Genres = append(result.Genres, el.Attr("title"))
		})
		
		epStart := e.DOM.Find("#episode_page > li").First().Find("a").AttrOr("ep_start", "")
		epEnd := e.DOM.Find("#episode_page > li").Last().Find("a").AttrOr("ep_end", "")
		movieID := e.ChildAttr("#movie_id", "value")
		alias := e.ChildAttr("#alias_anime", "value")
		
		episodes, err := FetchEpisode(epStart, epEnd, movieID, alias)
		if err != nil {
			return
		}

		result.Episodes = episodes
		result.TotalEpisodes = len(episodes)
	})

	if err := c.Visit(id); err != nil {
		return nil, err
	}

	return result, nil
}

func FetchEpisode(epStart, epEnd, movieID, alias string) ([]models.Episode, error) {
	var episodes []models.Episode

	c := colly.NewCollector()

	c.OnHTML("#episode_related > li", func(e *colly.HTMLElement) {
		episode := models.Episode{
			ID:     strings.Split(e.ChildAttr("a", "href"), "/")[1],
			Number: parseEpisodeNumber(e.ChildText("div.name")),
			URL:    fmt.Sprintf("%s/%s", baseURL, strings.TrimSpace(e.ChildAttr("a", "href"))),
		}
		episodes = append(episodes, episode)
	})

	if err := c.Visit(fmt.Sprintf("%s/load-list-episode?ep_start=%s&ep_end=%s&id=%s&default_ep=%d&alias=%s", ajaxURL, epStart, epEnd, movieID, 0, alias)); err != nil {
		return nil, err
	}

	for i, j := 0, len(episodes)-1; i < j; i, j = i+1, j-1 {
		episodes[i], episodes[j] = episodes[j], episodes[i]
	}

	return episodes, nil
}

func Downloads(id string) ([]models.EpisodeServer, error) {
	var servers []models.EpisodeServer

	c := colly.NewCollector()

	c.OnHTML("div.anime_video_body > div.anime_muti_link > ul > li", func(e *colly.HTMLElement) {
		url := e.ChildAttr("a", "data-video")
		if !strings.HasPrefix(url, "http") {
			url = "https:" + url
		}

		server := models.EpisodeServer{
			Name: strings.TrimSpace(strings.Replace(e.ChildText("a"), "Choose this server", "", 1)),
			URL:  url,
		}

		servers = append(servers, server)
	})

	if !strings.HasPrefix(id, baseURL) {
		id = fmt.Sprintf("%s/%s", baseURL, id)
	}

	err := c.Visit(id)
	if err != nil {
		return nil, err
	}

	return servers, nil
}

func parseEpisodeNumber(text string) float64 {
	number := strings.Replace(strings.TrimPrefix(text, "EP "), " ", "", -1)
	if num, err := strconv.ParseFloat(number, 64); err == nil {
		return num
	}
	return 0
}

func determineSubOrDub(title string) models.SubOrDub {
	if strings.Contains(strings.ToLower(title), "(dub)") {
		return models.DUB
	}
	return models.SUB
}
