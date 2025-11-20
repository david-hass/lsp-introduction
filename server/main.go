package main

import (
	"bufio"
	//	"context"
	"encoding/json"
	"fmt"
	"github.com/david-hass/lsp-introduction/tree_sitter_flow"
	sitter "github.com/smacker/go-tree-sitter"
	"io"
	"log"
	//	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	documentStore = make(map[string]string)
	storeMutex    = &sync.Mutex{}
	tsParser      *sitter.Parser
	flowLang      *sitter.Language
)

func main() {
	// lsp communicates via stdin/stdout
	f, err := os.OpenFile("/tmp/flow_lsp.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("--- flowlang server started ---")

	tsParser = sitter.NewParser()
	flowLang = sitter.NewLanguage(tree_sitter_flow.Language())
	tsParser.SetLanguage(flowLang)
	log.Println("flow tree-sitter-parser loaded")

	reader := bufio.NewReader(os.Stdin)

	// main read loop
	for {
		// 1. read header (Content-Length)
		contentLength, err := readContentLength(reader)
		if err != nil {
			if err == io.EOF {
				log.Println("client connection closed")
				return
			}
			log.Printf("failed to read header: %v", err)
			continue
		}

		// 2. read body
		body := make([]byte, contentLength)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			log.Printf("failed to read body: %v", err)
			continue
		}

		log.Printf("received: %s", string(body))
		handleMessage(body)
	}
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

func handleMessage(body []byte) {
	var request struct {
		ID     *json.RawMessage `json:"id"`
		Method string           `json:"method"`
		Params *json.RawMessage `json:"params"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		log.Printf("failed to parse: %v", err)
		return
	}

	// 'id' == nil means Notification -> no answer needed
	// 'id' != nil means Request -> answer expected

	switch request.Method {
	case "initialize":
		log.Println("Client 'initialize' received")

		// tell client about hover capabilities
		response := InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: FullSync,
				HoverProvider:    true,
			},
		}
		sendResponse(request.ID, response)

	case "initialized":
		log.Println("Client 'initialized' received")

	case "textDocument/hover":
		log.Println("Hover-Request received")

		hoverContent := MarkupContent{
			Kind:  "plaintext",
			Value: "hello from lsp server",
		}

		hoverResponse := Hover{
			Contents: hoverContent,
		}

		sendResponse(request.ID, hoverResponse)

	case "shutdown":
		log.Println("Shutdown-Request received")
		sendResponse(request.ID, nil)

	case "exit":
		log.Println("Exit-Notification received, shutting down...")
		os.Exit(0)

	default:
		log.Printf("unknown method: %s", request.Method)
		if request.ID != nil {
			// ...
		}
	}
}

func sendResponse(id *json.RawMessage, result interface{}) {
	type Response struct {
		ID     *json.RawMessage `json:"id"`
		Result interface{}      `json:"result"`
	}

	resp := Response{ID: id, Result: result}
	body, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Marshalling failed: %v", err)
		return
	}

	// LSP-Protokoll format: Header + \r\n\r\n + Body
	responseStr := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(body), body)

	// to stdout
	fmt.Print(responseStr)

	log.Printf("Answer sent: %s", string(body))
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

type TextDocumentSync int

const (
	FullSync TextDocumentSync = 1
)

type ServerCapabilities struct {
	TextDocumentSync TextDocumentSync
	HoverProvider    bool `json:"hoverProvider"`
	// e.g. CompletionProvider, DefinitionProvider ...
}

type Hover struct {
	Contents MarkupContent `json:"contents"`
}

type MarkupContent struct {
	Kind  string `json:"kind"` // "plaintext" or "markdown"
	Value string `json:"value"`
}
