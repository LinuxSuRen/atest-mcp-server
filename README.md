This is the MCP server of [atest](https://github.com/linuxsuren/api-testing).

## Get started

Please start the MCP server with [atest](https://github.com/linuxsuren/api-testing) gRPC port.

```shell
atest-store-mcp server --runner-address 127.0.0.1:64385
```

or start in Docker

```shell
docker run -p 7845:7845 ghcr.io/linuxsuren/atest-mcp-server --runner-address 127.0.0.1:64385
```

or start in npx

```shell
npx atest-mcp-server-launcher@latest server --mode=stdio --runner-address=localhost:64385
```

You can also set the MCP server mode with the `--mode` flag.

```shell
atest-store-mcp server --runner-address 127.0.0.1:64385 --mode=[sse|stdio]
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

or as stdio mode:
```json
{
  "mcpServers": {
    "atest": {
      "name": "atest-mcp-stdio",
      "type": "stdio",
      "description": "",
      "isActive": true,
      "command": "atest-store-mcp",
      "args": [
        "server",
        "-m=stdio",
        "--runner-address=localhost:64385"
      ]
    }
  }
}
```
