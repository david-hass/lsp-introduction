import * as path from 'path';
import { workspace, ExtensionContext, window } from 'vscode';
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  const config = workspace.getConfiguration('flow');
  const serverPath = config.get<string>('serverPath');

  if (!serverPath) {
    window.showErrorMessage("Flow LSP path is not configured!");
    return;
  }

  console.log(`Starting Flow LSP from: ${serverPath}`);

  const serverOptions: ServerOptions = {
    command: serverPath,
    transport: TransportKind.stdio,
    args: []
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: 'file', language: 'flow' }],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher('**/.clientrc')
    }
  };

  client = new LanguageClient(
    'flowLsp',
    'Flow Language Server',
    serverOptions,
    clientOptions
  );

  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
