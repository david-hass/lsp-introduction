"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode_1 = require("vscode");
const node_1 = require("vscode-languageclient/node");
let client;
function activate(context) {
    const config = vscode_1.workspace.getConfiguration('flow');
    const serverPath = config.get('serverPath');
    if (!serverPath) {
        vscode_1.window.showErrorMessage("Flow LSP path is not configured!");
        return;
    }
    console.log(`Starting Flow LSP from: ${serverPath}`);
    const serverOptions = {
        command: serverPath,
        transport: node_1.TransportKind.stdio,
        args: []
    };
    const clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'flow' }],
        synchronize: {
            fileEvents: vscode_1.workspace.createFileSystemWatcher('**/.clientrc')
        }
    };
    client = new node_1.LanguageClient('flowLsp', 'Flow Language Server', serverOptions, clientOptions);
    client.start();
}
function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}
//# sourceMappingURL=extension.js.map