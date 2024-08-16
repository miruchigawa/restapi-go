package downloader

import (
	"github.com/gocolly/colly/v2"

	models "miruchigawa.moe/restapi/internal/models/downloader"
)

func TiktokDownloader(url string) (*models.TiktokResult, error) {
	result := &models.TiktokResult{}
	c := colly.NewCollector()

	c.OnHTML("div.flex h2", func(e *colly.HTMLElement) {
		result.Nickname = e.Text
	})

	c.OnHTML("div.flex a.font-extrabold", func(e *colly.HTMLElement) {
		result.Username = e.Text
	})

	c.OnHTML("div.flex a > img", func(e *colly.HTMLElement) {
		result.Avatar = e.Attr("src")
	})

	c.OnHTML("div.flex p", func(e *colly.HTMLElement) {
		result.Description = e.Text
	})

	c.OnHTML("div.flex div.flex span", func(e *colly.HTMLElement) {
		index := e.Index
		switch index {
		case 0:
			result.Played = e.Text
		case 1:
			result.Commented = e.Text
		case 2:
			result.Saved = e.Text
		case 3:
			result.Shared = e.Text
		case 4:
			result.Song = e.Text
		}
	})

	c.OnHTML("#button-download-ready a", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		switch e.Index {
		case 0:
			result.Video = link
		case 2:
			result.Audio = link
		case 4:
			result.Thumbnail = link
		}
	})

	err := c.Post("https://ttsave.app/download", map[string]string{
		"language_id": "1",
		"query":       url,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
