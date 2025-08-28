This is the MCP server of [atest](https://github.com/linuxsuren/api-testing).

## Get started

Please start the MCP server with [atest](https://github.com/linuxsuren/api-testing) gRPC port.

```shell
atest-store-mcp server --runner-address 127.0.0.1:64385
```

## MCP Server

```json
{
  "mcpServers": {
    "atest": {
      "name": "atest",
      "type": "streamableHttp",
      "description": "The MCP server of atest",
      "isActive": true,
      "baseUrl": "http://localhost:7845",
      "disabledAutoApproveTools": []
    }
  }
}
```
