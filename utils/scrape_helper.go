package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractTextAndLink(s *goquery.Selection) (string, string) {
	link, exists := s.Find("a").Attr("href")
	if !exists {
		log.Print("Unable to find link.")
		link = ""
	}

	var text string
	nextNode := s.Next()
	for i := 0; i < 2 && nextNode.Length() > 0; i++ {
		paragraph := strings.TrimSpace(nextNode.Text())

		// Clean up the text
		paragraph = strings.ReplaceAll(paragraph, "\n", " ")
		paragraph = strings.ReplaceAll(paragraph, "\t", " ")
		paragraph = strings.ReplaceAll(paragraph, "\r", " ")
		paragraph = strings.Join(strings.Fields(paragraph), " ")

		if paragraph != "" {
			text += paragraph + " "
		}
		nextNode = nextNode.Next()
	}

	return strings.TrimSpace(text), link
}

func CleanHeadline(headline string) string {
	headline = strings.Trim(headline, "\"")
	headline = strings.ReplaceAll(headline, "\\", "")
	headline = strings.ReplaceAll(headline, "/", "")

	headline = strings.TrimSpace(headline)
	return headline
}

func ExtractTextFromURL(url string, UserAgent string, client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// More sophisticated content extraction
	var textBuilder strings.Builder
	doc.Find("article p, .article-content p, main p, .content p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			textBuilder.WriteString(text + " ")
		}
	})

	// Fallback if no content found in specific elements
	if textBuilder.Len() == 0 {
		doc.Find("p").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > 20 { // Filter out short paragraphs
				textBuilder.WriteString(text + " ")
			}
		})
	}

	// Get first 100 words
	fullText := textBuilder.String()
	words := strings.Fields(fullText)
	if len(words) > 100 {
		words = words[:100]
	}

	return strings.Join(words, " "), nil
}


//keywords based search
func BuildKeywordsMap(keywords []string) map[string]bool {
	kwMap := make(map[string]bool)
	for _, kw := range keywords {
		kwMap[strings.ToLower(kw)] = true
	}
	return kwMap
}

// ContainsKeywords checks if text contains any of the keywords
// Returns both a boolean indicating if any keyword was found and the matched keyword
func ContainsKeywords(text string, kwMap map[string]bool) (bool, string) {
	if len(kwMap) == 0 {
		return true, ""
	}

	lowerText := strings.ToLower(text)
	for kw := range kwMap {
		if strings.Contains(lowerText, kw) {
			return true, kw
		}
	}
	return false, ""
}
