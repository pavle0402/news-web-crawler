package handlers

import (
	"crawler/models"
	"crawler/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const UserAgent = "CrawlerBotPavle"

func CrawlerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed.", http.StatusMethodNotAllowed)
		return
	}

	var payload models.CrawlerRequest
	err := json.NewDecoder(r.Body).Decode((&payload))
	if err != nil {
		log.Println("Error decoding payload:", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	results := make(map[string]interface{})
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	for _, req := range payload.Requests {
		if req.URL == "" || req.Depth < 0 {
			http.Error(w, "Invalid payload format", http.StatusBadRequest)
			return
		}

		// Process the initial URL
		pageData, err := processURL(req.URL, req.Depth, client, req.Keywords)
		if err != nil {
			log.Printf("Error processing URL %s: %v", req.URL, err)
			continue // Skip this URL but continue with others
		}

		results[req.URL] = pageData
		log.Printf("Successfully processed URL: %s", req.URL)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	prettyJSON, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Write(prettyJSON)
}

// NewsData represents a single news article with its headline and content
type NewsData struct {
	Headline string `json:"headline"`
	Text     string `json:"text"`
	TextLink string `json:"text_link"`
	Category string `json:"category"`
}

// Helper function to process a single URL and its depth
func processURL(url string, depth int32, client *http.Client, keywords []string) (map[string]interface{}, error) {
	// Check robots.txt first
	allowed, err := utils.ScrapeAllowed(url, UserAgent)
	if err != nil {
		return nil, fmt.Errorf("error checking robots.txt: %v", err)
	}
	if !allowed {
		return nil, fmt.Errorf("crawling not allowed by robots.txt")
	}

	// Fetch the page
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	// Parse the document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	// Extract data
	title := doc.Find("title").Text()

	// First pass - collect all headlines and their links
	headlinesData := collectHeadlines(doc)

	// Second pass - process articles if depth > 0
	if depth > 0 {
		processArticleContent(headlinesData, client)
	}

	// Apply keyword filtering as a final step
	filteredNews := filterByKeywords(headlinesData, keywords)

	return map[string]interface{}{
		"title":     title,
		"headlines": filteredNews,
	}, nil
}

// headlineInfo represents collected data for a headline
type headlineInfo struct {
	headline string
	text     string
	textLink string
}

// collets all headlines and their links from the document
func collectHeadlines(doc *goquery.Document) map[string]*headlineInfo {
	headlinesMap := make(map[string]*headlineInfo)

	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		tag := goquery.NodeName(s)
		if strings.HasPrefix(tag, "h") && len(tag) == 2 {
			headline := utils.CleanHeadline(strings.TrimSpace(s.Text()))
			if headline != "" {
				_, link := utils.ExtractTextAndLink(s)

				if _, exists := headlinesMap[headline]; !exists {
					headlinesMap[headline] = &headlineInfo{
						headline: headline,
						textLink: link,
					}
				}
			}
		}
	})

	return headlinesMap
}

// fetches and extracts content for each headline
func processArticleContent(headlinesMap map[string]*headlineInfo, client *http.Client) {
	for _, info := range headlinesMap {
		if info.textLink != "" {
			textContent, err := utils.ExtractTextFromURL(info.textLink, UserAgent, client)
			if err == nil {
				info.text = textContent
			}
		}
	}
}

// keyword filtering for the collected news data
func filterByKeywords(headlinesMap map[string]*headlineInfo, keywords []string) []NewsData {
	kwMap := utils.BuildKeywordsMap(keywords)
	filteredNews := make([]NewsData, 0)

	if len(kwMap) == 0 {
		for _, info := range headlinesMap {
			filteredNews = append(filteredNews, NewsData{
				Headline: info.headline,
				Text:     info.text,
				TextLink: info.textLink,
				Category: "",
			})
		}
		return filteredNews
	}

	for _, info := range headlinesMap {
		//applied to both, headlineinfo(headline and link) and actual text extracted - may not always be precise
		// needed additional filtering - figure it out
		headlineMatch, headlineKw := utils.ContainsKeywords(info.headline, kwMap)
		contentMatch, contentKw := utils.ContainsKeywords(info.text, kwMap)

		if headlineMatch || contentMatch {
			category := headlineKw
			if category == "" {
				category = contentKw
			}

			filteredNews = append(filteredNews, NewsData{
				Headline: info.headline,
				Text:     info.text,
				TextLink: info.textLink,
				Category: category,
			})
		}
	}

	return filteredNews
}
