package main

import "encoding/json"

type Response struct {
	ID     *json.RawMessage `json:"id"`
	Result any              `json:"result"`
}

type Notification struct {
	Method string `json:"method"`
	Params any    `json:"params"`
}

type DiagnosticsResponse struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type DiagnosticSeverity int

const (
	DiagnosticsError   DiagnosticSeverity = 1
	DiagnosticsWarning DiagnosticSeverity = 2
	DiagnosticsInfo    DiagnosticSeverity = 3
	DiagnosticsHint    DiagnosticSeverity = 4
)

type Diagnostic struct {
	Range    Range              `json:"range"`
	Severity DiagnosticSeverity `json:"severity"` // 1: Error, 2: Warning, 3: Info, 4: Hint
	Source   string             `json:"source"`
	Message  string             `json:"message"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type MarkupContent struct {
	Kind  string `json:"kind"` // "plaintext" or "markdown"
	Value string `json:"value"`
}

type HoverResponse struct {
	Contents MarkupContent `json:"contents"`
}
