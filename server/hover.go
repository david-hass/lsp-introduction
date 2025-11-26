package main

import (
	"fmt"
	"log"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

func hoverFeature(textDocument TextDocumentIdentifier,position     Position) (*HoverResponse,bool){
	doc, ok := getDocument(textDocument.URI)
	if !ok {
		log.Printf("document for hover not found: %s", textDocument.URI)
		return nil,false
	}

	tree := tsParser.Parse([]byte(doc), nil)
	rootNode := tree.RootNode()

	point := sitter.Point{Row: uint(position.Line), Column: uint(position.Character)}

	node := rootNode.NamedDescendantForPointRange(point, point)

	if node == nil {
		log.Printf("no node found at %v:%v", position.Line, position.Character)
		return nil,false
	}

	log.Printf("node kind: %s", node.Kind())

	hoverText := ""
	switch node.Kind() {
	case "source_definition":
		hoverText = "defines a **source**."
	case "source_body_path":
		hoverText = "specifies the file path from which the data should be read."
	case "task_definition":
		hoverText = "defines a **task** in the pipeline."
	case "task_body_input":
		hoverText = "specifies where the input data comes from."
	case "task_body_transformer":
		hoverText = "specifies how the data is to be processed."
	case "sink_definition":
		hoverText = "defines the **sink**, where the data is stored."
	case "sink_body_input":
		hoverText = "specifies where the input data comes from."
	case "sink_body_path":
		hoverText = "specifies the file path where the data is saved to."
	case "identifier":
		start := node.StartByte()
		end := node.EndByte()
		if end > uint(len(doc)) {
			break
		}
		text := doc[start:end]
		hoverText = fmt.Sprintf("identifier: `%s`", text)
	case "string_literal":
		start := node.StartByte()
		end := node.EndByte()
		if end > uint(len(doc)) {
			break
		}
		text := doc[start:end]
		hoverText = fmt.Sprintf("value: `%s`", text)
	default:
		break
	}

	hoverContent := MarkupContent{
		Kind:  "markdown",
		Value: hoverText,
	}

	hoverResponse := HoverResponse{
		Contents: hoverContent,
	}

	return &hoverResponse, true
}
