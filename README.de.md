[![en](https://img.shields.io/badge/lang-en-red.svg)](https://github.com/david-hass/lsp-introduction/blob/main/README.md)
[![de](https://img.shields.io/badge/lang-de-yellow.svg)](https://github.com/david-hass/lsp-introduction/blob/main/README.de.md)


Eine exemplarische Implementierung des Language Server Protocols in Go.

Dieses Projekt demonstriert die Funktionsweise moderner IDE-Features, indem ein eigener Language Server für eine eigene Domain Specific Language, genannt flow, entwickelt wurde.

Die Implementierung nutzt einen per tree-sitter generierten Parser. Glücklicherweise erzeugt tree-sitter neben dem Parser auch zugehörige Bindings für diverse Sprachen, darunter auch go.
Zudem wurde das go-Package [github.com/tree-sitter/go-tree-sitter](https://github.com/tree-sitter/go-tree-sitter) genutzt, welche eine Reihe von Funktionalitäten zur verwendung eines solchen Parsers bietet.

### Projektstruktur

<pre>
.
├── grammar/    # Die Tree-sitter Grammatik
├── server/     # Der Language Server (Go)
├── vscode/     # Der VS Code Client (TypeScript Extension)
└── examples/   # Beispiel-Dateien (.flow)
</pre>

### FlowLang

FlowLang ist eine minimale (proof-of-concept) DSL zur Definition von Daten-Pipelines.

Die Sprache erlaubt, zu definieren, woher die Daten stammen (source):
```hcl
source "raw_user_data" {
    path: "/data/users.csv"
    encoding: "utf-8"
}
```

welche Arten von Verarbeitung auf diese Daten angewandt werden (task):
```hcl
task "filter_active_users" {
    # Nimmt Output von "raw_user_data" als Input
    input: raw_user_data
    transformer: "filter_by_column 'status' == 'active'"
}
```


und wohin die Daten transportiert werden (sink):
```hcl
sink "active_user_report" {
    # Nimmt den Output des Anonymisierungs-Tasks
    input: filter_active_users
    path: "/reports/active_users.json"
}
```
Die Grammatik ist sehr einfach gehalten und die Notwendigkeit, des explizieten Festlegens von Präzedenzen oder allgemeiner, der Umgang mit Eindeutigkeitsproblemen, wurde vermieden. 

Der Server analysiert den von dem generierten tree-sitter Parser erzeugten CST (concrete syntax tree) und stellt Editoren über das LSP verschiedene Features bereit.

### Implementierte Features

1. [Contextual Hover](https://github.com/david-hass/lsp-introduction/blob/main/server/hover.go) (textDocument/hover)
    Bewege den Mauszeiger über Schlüsselwörter wie **sink**, **task** usw. um Dokumentation zur Syntax zu erhalten.

2. [Semantische Diagnose](https://github.com/david-hass/lsp-introduction/blob/main/server/diagnostics.go) (textDocument/publishDiagnostics)
    Der Server prüft logische Referenzen. Wenn z.B. ein **task** in **input** verwendet wird, der nicht existiert, wird dies als Fehler markiert.


### Build & Installation

Das Projekt wurde unter folgender Umgebung entwickelt:
Linux 6.16.8-1-MANJARO x86_64 GNU/Linux
go 1.25.1
npm 11.6.1
gcc 15.2.1

1. Grammatik generieren **(Optional)**
    Der C-Code des Parsers liegt bereits in server/parser/. Falls Änderungen an der Grammatik (grammar/grammar.js) vorgenommen werden, muss sie neu generiert werden:

    <pre>
    cd grammar
    npm install
    npx tree-sitter generate && npx tree-sitter build
    \# Danach müssen die Dateien manuell nach ../server/parser/ kopiert werden
    \# alternativ das Skript: ../server/copyparser.sh nutzen
    </pre>

2. Server bauen

    Der Server muss kompiliert werden, damit die Clients ihn starten können.

    <pre>
    cd server
    go mod tidy
    CGO_ENABLED=1 go build -o flow-lsp
    </pre>


    Test: Führe ./flow-lsp aus. Es sollte starten und auf Input warten. Außerdem sollte eine /tmp/flow_lsp.log
    (unter Windows wahrscheinlich C:\Users\Name\AppData\Local\Temp\flow_lsp.log)
    mit folgendem Inhalt erstellt worden sein:
    <pre>
    2025/11/27 10:35:07 --- flowlang server started ---
    2025/11/27 10:35:07 tree sitter parser loaded.
    2025/11/27 10:35:07 flow tree sitter parser loaded
    </pre>

### Integration in VS Code

Der VS Code Client befindet sich im Ordner vscode/.

Abhängigkeiten installieren:
<pre>
cd vscode
npm install
</pre>

1. Führe in VSCode den Befehl: "Developer: Install Extension from Location..." aus
2. Wähle dann den vscode Ordner
3. Öffne nun die VSCode Einstellungen und suche nach "flow"
4. Trage bei Flow: Server Path den absoluten Pfad zu deiner flow-lsp Binary ein


### Integration in Neovim

Folgendermaßen könnte die Integration unter Verwendung von lazy.nvim stattfinden.
In der plugin config():
<pre>
local flow_lsp_path = vim.fn.expand("~/projects/lsp-introduction/server/flow-lsp")
local util = require 'lspconfig.util'

vim.lsp.config.flow = {
  cmd = { flow_lsp_path },
  filetypes = { "flow" },
  root_dir = util.find_git_ancestor() or util.path.dirname,
  capabilities = capabilities,
}
vim.lsp.enable('flow')
</pre>


### Referenzen

[LSP](https://microsoft.github.io/language-server-protocol/)

[Tree-sitter](https://tree-sitter.github.io/tree-sitter/)

[go-tree-sitter Bindings](https://github.com/tree-sitter/go-tree-sitter)
