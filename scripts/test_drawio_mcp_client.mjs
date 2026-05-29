import { Client } from "file:///C:/Users/11818/AppData/Roaming/npm/node_modules/@drawio/mcp/node_modules/@modelcontextprotocol/sdk/dist/esm/client/index.js";
import { StdioClientTransport } from "file:///C:/Users/11818/AppData/Roaming/npm/node_modules/@drawio/mcp/node_modules/@modelcontextprotocol/sdk/dist/esm/client/stdio.js";

const transport = new StdioClientTransport({
  command: "C:\\nvm4w\\nodejs\\node.exe",
  args: [
    "C:\\Users\\11818\\AppData\\Roaming\\npm\\node_modules\\@drawio\\mcp\\src\\index.js",
  ],
});

const client = new Client({
  name: "codex-drawio-sdk-test",
  version: "1.0.0",
});

await client.connect(transport);
const tools = await client.listTools();
console.log(JSON.stringify(tools, null, 2).slice(0, 2000));
await client.close();
