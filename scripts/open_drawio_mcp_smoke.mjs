import { Client } from "file:///C:/Users/11818/AppData/Roaming/npm/node_modules/@drawio/mcp/node_modules/@modelcontextprotocol/sdk/dist/esm/client/index.js";
import { StdioClientTransport } from "file:///C:/Users/11818/AppData/Roaming/npm/node_modules/@drawio/mcp/node_modules/@modelcontextprotocol/sdk/dist/esm/client/stdio.js";

const transport = new StdioClientTransport({
  command: "C:\\nvm4w\\nodejs\\node.exe",
  args: [
    "C:\\Users\\11818\\AppData\\Roaming\\npm\\node_modules\\@drawio\\mcp\\src\\index.js",
  ],
});

const client = new Client({
  name: "codex-drawio-open-smoke",
  version: "1.0.0",
});

await client.connect(transport);
const result = await client.callTool({
  name: "open_drawio_xml",
  arguments: {
    dark: "false",
    lightbox: false,
    content:
      '<mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/><mxCell id="u" value="用户" style="whiteSpace=wrap;html=1;fillColor=#ffffff;strokeColor=#000000;" vertex="1" parent="1"><mxGeometry x="80" y="80" width="120" height="60" as="geometry"/></mxCell></root></mxGraphModel>',
  },
});
console.log(JSON.stringify(result, null, 2));
await client.close();
