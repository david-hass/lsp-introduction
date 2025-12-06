package main

import (
	"context"
	"fmt"
	"log"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

const definitionsQueryStr = `
		(source_definition name: (string_literal) @def.name)
		(task_definition   name: (string_literal) @def.name)
		(sink_definition   name: (string_literal) @def.name)
	`

const referencesQueryStr = `
		(prop_input input: (identifier) @ref.name)
`

func diagnosticsFeature(_ context.Context, uri string, content string) *DiagnosticsResponse {
	sourceBytes := []byte(content)

	// parse document
	tree := tsParser.Parse([]byte(content), nil)
	root := tree.RootNode()
	defer tree.Close()

	log.Println("querying for definitions")
	definedNames := queryForDefinitions(root, sourceBytes)

	log.Println("checking for invalid references")
	diagnostics := checkForInvalidReferences(root, sourceBytes, definedNames)

	response := DiagnosticsResponse{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	return &response
}

func queryForDefinitions(root *sitter.Node, sourceBytes []byte) map[string]bool {
	definedNames := make(map[string]bool)

	defQuery, _ := sitter.NewQuery(flowLang, definitionsQueryStr)
	defer defQuery.Close()

	defQC := sitter.NewQueryCursor()
	defer defQC.Close()

	matches := defQC.Matches(defQuery, root, sourceBytes)

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		// there should be just one capture: @def.name
		for _, capture := range match.Captures {
			start := capture.Node.StartByte()
			end := capture.Node.EndByte()

			if end > uint(len(sourceBytes)) || end-1 == start {
				continue
			}

			name := string(sourceBytes[start+1 : end-1])
			log.Printf("definition found: %s", name)
			definedNames[name] = true
		}
	}

	return definedNames
}

func checkForInvalidReferences(root *sitter.Node, sourceBytes []byte, definedNames map[string]bool) []Diagnostic {
	diagnostics := make([]Diagnostic, 0)

	refQuery, _ := sitter.NewQuery(flowLang, referencesQueryStr)
	defer refQuery.Close()

	refQC := sitter.NewQueryCursor()
	defer refQC.Close()

	matches := refQC.Matches(refQuery, root, sourceBytes)

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, capture := range match.Captures {
			if capture.Node.Kind() == "identifier" {
				start := capture.Node.StartByte()
				end := capture.Node.EndByte()

				if end > uint(len(sourceBytes)) {
					continue
				}

				refName := string(sourceBytes[start:end])
				log.Printf("reference found: %s", refName)

				if !definedNames[refName] {
					log.Printf("invalid reference found: %s", refName)

					startPoint := capture.Node.StartPosition()
					endPoint := capture.Node.EndPosition()

					diagnostic := Diagnostic{
						Range: Range{
							Start: Position{Line: int(startPoint.Row), Character: int(startPoint.Column)},
							End:   Position{Line: int(endPoint.Row), Character: int(endPoint.Column)},
						},
						Severity: DiagnosticsError,
						Source:   "flow-lsp",
						Message:  fmt.Sprintf("invalid reference: task '%s' is not defined yet", refName),
					}

					diagnostics = append(diagnostics, diagnostic)
				}
			}
		}

	}

	return diagnostics
}
