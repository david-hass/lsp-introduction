package tree_sitter_flowlangparser_test

import (
	"testing"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_flowlangparser "github.com/tree-sitter/tree-sitter-flowlangparser/bindings/go"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(tree_sitter_flowlangparser.Language())
	if language == nil {
		t.Errorf("Error loading flowlang parser grammar")
	}
}
