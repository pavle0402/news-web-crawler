package utils

import (
	"net/http"
	"net/url"

	"github.com/temoto/robotstxt"
)

func ScrapeAllowed(urlStr string, userAgent string) (bool, error) {
	parsedUrl, err := url.Parse(urlStr)

	if err != nil {
		return false, err
	}

	robotsUrl := parsedUrl.Scheme + "://" + parsedUrl.Host + "/robots.txt"
	resp, err := http.Get(robotsUrl)
	if err != nil {
		return false, err
	}

	robotsData, err := robotstxt.FromResponse(resp)
	if err != nil {
		return false, err
	}

	group := robotsData.FindGroup(userAgent)
	return group.Test(parsedUrl.Path), nil
}
