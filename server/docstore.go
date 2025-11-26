package main

import (
	"log"
	"net/url"
	"sync"
)

var (
	documentStore = make(map[string]string)
	storeMutex    = &sync.Mutex{}
)

func storeDocument(uri string, text string) bool {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	normalizedURI, err := url.QueryUnescape(uri)
	if err != nil {
		log.Printf("un-escape URI failed: %s", uri)
		return false
	}
	documentStore[normalizedURI] = text
		return true
}

func getDocument(uri string) (string, bool) {
	storeMutex.Lock()
	defer storeMutex.Unlock()
	normalizedURI, err := url.QueryUnescape(uri)
	if err != nil {
		log.Printf("un-escape URI failed: %s", uri)
		return "", false
	}
	doc, ok := documentStore[normalizedURI]
	return doc, ok
}
