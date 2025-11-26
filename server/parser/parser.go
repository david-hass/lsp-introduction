package parser

// #cgo CFLAGS: -std=c11 -fPIC
// #cgo CFLAGS: -I. -I./tree_sitter
// #include "tree_sitter/api.h"
// TSLanguage *tree_sitter_flow();
import "C"
import (
	"unsafe"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// GetLanguage liefert die offizielle Sprach-Definition
func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_flow())
	return sitter.NewLanguage(ptr)
}
