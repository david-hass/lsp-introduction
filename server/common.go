package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

func readContentLength(reader *bufio.Reader) (int, error) {
	var length int
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}

		if strings.HasPrefix(line, "Content-Length:") {
			parts := strings.Split(line, ":")
			lengthStr := strings.TrimSpace(parts[1])
			length, err = strconv.Atoi(lengthStr)
			if err != nil {
				return 0, fmt.Errorf("invalid Content-Length: %s", parts[1])
			}
		}

		if line == "\r\n" {
			if length == 0 {
				return 0, fmt.Errorf("no Content-Length found")
			}
			return length, nil
		}
	}
}

func sendResponse(id *json.RawMessage, result any) {

	resp := Response{ID: id, Result: result}
	body, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Marshalling failed: %v", err)
		return
	}

	fmt.Printf("Content-Length: %d\r\n\r\n%s", len(body), body)

	log.Printf("Answer sent: %s", string(body))
}

func sendNotification(method string, params any) {
	notif := Notification{Method: method, Params: params}

	body, _ := json.Marshal(notif)

	fmt.Printf("Content-Length: %d\r\n\r\n%s", len(body), body)

	log.Printf("Answer sent: %s", string(body))
}
