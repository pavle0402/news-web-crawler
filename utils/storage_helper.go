package utils

import (
	"fmt"
	"log"
	"time"
)

type NewsDocument struct {
	URL       string      `bson:"url"`
	Headlines interface{} `bson:"headlines"`
	Title     string      `bson:"title"`
	Created   time.Time   `bson:"created"`
}

func DataForInsert(data map[string]interface{}) []interface{} {
	var docs []interface{}

	location, err := time.LoadLocation("Europe/Belgrade")
	if err != nil {
		log.Printf("Error occured while setting up the  time: %v", err)
	}
	current_time := time.Now().In(location)

	for url, content := range data {
		contentMap, ok := content.(map[string]interface{})
		if !ok {
			continue
		}
		docs = append(docs, NewsDocument{
			URL:       url,
			Headlines: contentMap["headlines"],
			Title:     fmt.Sprintf("%v", contentMap["title"]),
			Created:   current_time,
		})
	}

	return docs
}
