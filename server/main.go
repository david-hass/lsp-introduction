package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/david-hass/lsp-introduction/parser"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

var (
	tsParser *sitter.Parser
	flowLang *sitter.Language
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
	flowLang = parser.GetLanguage()

	if err := tsParser.SetLanguage(flowLang); err != nil {
		log.Fatalf("failed to load tree sitter language: %v", err)
	}
	log.Println("tree sitter parser loaded.")

	log.Println("flow tree sitter parser loaded")

	reader := bufio.NewReader(os.Stdin)

	// main read loop
	for {
		contentLength, err := readContentLength(reader)
		if err != nil {
			if err == io.EOF {
				log.Println("client connection closed")
				return
			}
			log.Printf("failed to read header: %v", err)
			continue
		}

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

// 'id' == nil means Notification -> no answer needed
// 'id' != nil means Request -> answer expected
func handleMessage(body []byte) {
	var request Request
	if err := json.Unmarshal(body, &request); err != nil {
		log.Printf("failed to parse: %v", err)
		return
	}

	log.Printf("received: %v", request)

	switch request.Method {

	case "initialize":
		log.Println("client 'initialize' received")

		response := InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: FullSync,
				HoverProvider:    true,
			},
		}
		sendResponse(request.ID, response)

	case "initialized":
		log.Println("'initialized' received")

	case "textDocument/didOpen":
		log.Println("'textDocument/didOpen' received")

		var params DidOpenParams
		if err := json.Unmarshal(*request.Params, &params); err != nil {
			log.Printf("failed to parse didOpen: %v", err)
			return
		}

		txt := params.TextDocument.Text

		ok := storeDocument(params.TextDocument.URI, txt)
		if !ok {
			return
		}

		diagnosticsResponse := diagnosticsFeature(context.Background(), params.TextDocument.URI, txt)
		sendNotification("textDocument/publishDiagnostics", diagnosticsResponse)

	case "textDocument/didChange":
		log.Println("'textDocument/didChange' received")

		var params DidChangeParams
		if err := json.Unmarshal(*request.Params, &params); err != nil {
			log.Printf("failed to parse didChange: %v", err)
			return
		}

		if len(params.ContentChanges) > 0 {
			txt := params.ContentChanges[0].Text

			storeDocument(params.TextDocument.URI, txt)

			diagnosticsResponse := diagnosticsFeature(context.Background(), params.TextDocument.URI, txt)
			sendNotification("textDocument/publishDiagnostics", diagnosticsResponse)
		}

	case "textDocument/hover":
		log.Println("'textDocument/hover' received")

		var params HoverParams
		if err := json.Unmarshal(*request.Params, &params); err != nil {
			log.Printf("failed to parse HoverParams: %v", err)
			sendResponse(request.ID, nil)
			return
		}
		hoverResponse, ok := hoverFeature(params.TextDocument, params.Position)
		if !ok {
			sendResponse(request.ID, nil)
			return
		}

		sendResponse(request.ID, hoverResponse)

	case "shutdown":
		log.Println("'shutdown' received")

		sendResponse(request.ID, nil)

	case "exit":
		log.Println("'exit' received")

		os.Exit(0)

	default:
		log.Printf("received unknown method: '%s'", request.Method)

		if request.ID != nil {
			// ...
		}
	}
}
