package downloader

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"

	models "miruchigawa.moe/restapi/internal/models/downloader"
)

func GetMediafireInfo(url string) (*models.MediafireInfo, error) {
	info := &models.MediafireInfo{}
	c := colly.NewCollector()

	c.OnHTML("body", func(e *colly.HTMLElement) {
		doc := e.DOM

		info.URL = strings.TrimSpace(doc.Find("#downloadButton").AttrOr("href", ""))

		intro := doc.Find("div.dl-info > div.intro")
		info.Filename = strings.TrimSpace(intro.Find("div.filename").Text())
		info.Filetype = strings.TrimSpace(intro.Find("div.filetype > span").Eq(0).Text())

		re := regexp.MustCompile(`\(\.(.*?)\)`)
		match := re.FindStringSubmatch(strings.TrimSpace(intro.Find("div.filetype > span").Eq(1).Text()))
		if len(match) > 1 {
			info.Ext = strings.TrimSpace(match[1])
		} else {
			info.Ext = "bin"
		}

		li := doc.Find("div.dl-info > ul.details > li")
		info.Uploaded = strings.TrimSpace(li.Eq(1).Find("span").Text())
		info.Filesize = strings.TrimSpace(li.Eq(0).Find("span").Text())
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return info, nil
}
